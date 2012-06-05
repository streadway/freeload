package freeload

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
)

// Contains the structured fields of a complete Get
type Result struct {
	RequestURI string         `json:"-"`
	DataURI    string         `json:"uri,omitempty"`
	Response   *http.Response `json:"-"`
	Error      *ErrResult     `json:"err,omitempty"`
}

// Marshal Errors into JSON, if you find a way of marshalling error directly
// to JSON without this embedded type, please replace.
type ErrResult struct {
	error
}

func (me *ErrResult) MarshalJSON() ([]byte, error) {
	return json.Marshal(me.Error())
}

// Fetch-me-a data URI!  Uses an injected client for tests or
// production that will setup the proper proxy
func Get(c http.Client, url string) Result {
	res, err := c.Get(url)
	if err != nil {
		return Result{url, "", nil, &ErrResult{err}}
	}

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return Result{url, "", res, &ErrResult{err}}
	}

	if contentTypes, ok := res.Header["Content-Type"]; ok {
		return Result{url, Format(contentTypes[0], body), res, nil}
	} else {
		return Result{url, Format("", body), res, nil}
	}

	panic("unreachable")
}

// Given a content type potentially from a trusted server and the fully read
// payload body, build the data-uri format defined in:
// http://tools.ietf.org/html/rfc2397
func Format(contentType string, body []byte) string {
	// len(data:;base64,) + rest as capacity
	uri := make([]byte, 0, 16+len(contentType)+base64.StdEncoding.EncodedLen(len(body)))

	uri = append(uri, []byte("data:")...)

	// Split out each param, to extract the mime and maybe the charset
	parts := strings.Split(contentType, ";")

	// Append the mimetype
	uri = append(uri, []byte(parts[0])...)

	// Pass through parameters, including 'charset=', defined as a parameter by
	// containing an equals sign in the part
	for _, part := range parts[1:] {
		if strings.Index(part, "=") > 0 {
			uri = append(uri, []byte(";")...)
			uri = append(uri, []byte(strings.TrimSpace(part))...)
		}
	}

	// Append the optional base64 encoding and the start of the payload (,)
	uri = append(uri, []byte(";base64,")...)

	// Append the base64 encoded content
	encoded := make([]byte, base64.StdEncoding.EncodedLen(len(body)))
	base64.StdEncoding.Encode(encoded, body)

	return string(append(uri, encoded...))
}

// Extract the list of absolute URLs from the query parameters in an http.Request
// this is the command decoder that will take the compressed prefix/suffix scheme
// and return expanded URLs
//
// Query parameters:
//
// p = 1..1 prefix, such as http://s3.amazonaws.com/base/path/00
// i = 0..n inner, such as 101-26273-x23sn, if none exists, only use prefix
// s = 0..1 suffix, such as _m.png will be applied to all URLs including prefix-only requests
func DecodeUrls(query url.Values) (urls []string, err error) {
	prefix := query.Get("p")
	if prefix == "" {
		err = fmt.Errorf("Must contain the (p)refix query parameter")
		return
	}

	suffix := query.Get("s")
	inners := query["i"]

	if len(inners) == 0 {
		urls = append(urls, prefix+suffix)
	} else {
		for _, inner := range inners {
			urls = append(urls, prefix+inner+suffix)
		}
	}

	return
}

func GetAll(c http.Client, urls []string, after time.Duration) (results map[string]Result) {
	results = make(map[string]Result, len(urls))
	requests := make(chan Result, len(urls))
	timeout := &ErrResult{fmt.Errorf("Request time %v exceeded", after)}

	for _, u := range urls {
		results[u] = Result{u, "", nil, timeout}
	}

	for _, u := range urls {
		go func(url string) { requests <- Get(c, url) }(u)
	}

	deadline := time.After(after)

	for {
		select {
		case res := <-requests:
			results[res.RequestURI] = res
		case <-deadline:
			return
		}
	}

	return
}

// Find the lowest max-age of the sucessful results to use for this response
// return -1 if any of the upstreams are uncacheable
func MaxAge(results map[string]Result) (maxAge int) {
	param := "max-age="
	maxAge = -1

	for _, res := range results {
		if res.Error == nil && res.Response != nil {
			cc := res.Response.Header.Get("Cache-Control")
			if cc == "" {
				// This response doesn't have cache control, so none are cacheable
				return -1
			}

			if i := strings.Index(cc, param); i >= 0 {
				if candidate, err := strconv.Atoi(cc[len(param):]); err != nil {
					if candidate > maxAge {
						maxAge = candidate
					}
				}
			}
		}
	}

	return
}

func WriteResponseJSON(w http.ResponseWriter, results map[string]Result) {
	maxAge := MaxAge(results)
	if maxAge > 0 {
		w.Header().Set("Cache-Control", fmt.Sprintf("public,max-age=%d", maxAge))
	} else {
		w.Header().Set("Cache-Control", "private,no-store,max-age=0")
	}

	w.Header().Set("Content-Type", "text/json;charset=utf-8")

	if err := json.NewEncoder(w).Encode(results); err != nil {
		http.Error(w, fmt.Sprintf("JSON encoding error: %v", err), 500)
	}
}

func LoadJSONHandler(c http.Client, timeout time.Duration) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case "GET":
			urls, err := DecodeUrls(r.URL.Query())
			if err != nil {
				http.Error(w, fmt.Sprintf("bad query parameters: %v", err), 400)
			}

			WriteResponseJSON(w, GetAll(c, urls, timeout))
		default:
			http.Error(w, "unsupported method", 400)
		}
	})
}

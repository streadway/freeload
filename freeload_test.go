// Copyright (C) 2012 Sean Treadway <treadway@gmail.com>, SoundCloud Ltd.
// All rights reserved.  See README.md for license details.

package freeload

import (
	"compress/gzip"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
	"time"
)

func matchFormat(t *testing.T, actual, expected string, msg string) {
	if actual != expected {
		t.Errorf("%s: expected %v, got: %v", msg, expected, actual)
	}
}

func backendClient(srv *httptest.Server) http.Client {
	url, _ := url.Parse(srv.URL)
	return http.Client{
		Transport: &http.Transport{
			Proxy: http.ProxyURL(url),
		},
	}
}

func backendPathEcho(t *testing.T) *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(r.URL.Path))
	}))
}

func backendPathEchoTimeoutAfterOne(t *testing.T) *httptest.Server {
	var i int
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if i > 0 {
			time.Sleep(1e7)
		}
		w.Header().Set("Cache-Control", "max-age=10")
		w.Write([]byte(r.URL.Path))
		i++
	}))
}

func TestFetch(t *testing.T) {
	origin := backendPathEcho(t)
	defer origin.Close()

	srv := httptest.NewServer(
		GzipHandler(
			CORS("*",
				LoadJSONHandler(backendClient(origin), 1e6))))
	defer srv.Close()

	req, _ := http.NewRequest("GET", "http://test/?p=http://test/&i=ohai&i=lol", nil)
	req.Header.Set("Accept-Encoding", "gzip")

	cli := backendClient(srv)
	res, err := (&cli).Do(req)

	if err != nil {
		t.Error("request failed")
	}

	if res.Header.Get("Access-Control-Allow-Origin") != "*" {
		t.Error("not CORS ready")
	}

	if res.Header.Get("Content-Encoding") != "gzip" {
		t.Error("not encoded with gzip")
	}

	if !strings.Contains(res.Header.Get("Vary"), "Accept-Encoding") {
		t.Error("vary header incorrect")
	}

	gz, _ := gzip.NewReader(res.Body)
	body, err := ioutil.ReadAll(gz)

	// TODO decode format and test fields individually

	if !strings.Contains(string(body), "http://test/ohai") {
		t.Error("does not contain ohai")
	}

	if !strings.Contains(string(body), "http://test/lol") {
		t.Error("does not contain lol")
	}
}

func TestGetAll(t *testing.T) {
	srv := backendPathEchoTimeoutAfterOne(t)
	defer srv.Close()

	u1, u2 := "http://test/ohai?1", "http://test/ohai?2"

	res := GetAll(backendClient(srv), []string{u1, u2}, 1e7)

	if len(res) != 2 {
		t.Error("Not enough results returned")
	}

	var good, bad Result

	if res[u1].Error != nil {
		good, bad = res[u2], res[u1]
	} else {
		good, bad = res[u1], res[u2]
	}

	matchFormat(t, good.DataURI, "data:text/plain;charset=utf-8;base64,L29oYWk=",
		"Test server doesn't match")

	matchFormat(t, bad.DataURI, "",
		"Timeout should not contain a DataURI")

	cc := good.Response.Header.Get("Cache-Control")
	if cc != "max-age=10" {
		t.Errorf("the mock backend should respond with cacheable content: %v", cc)
	}

	if bad.Response != nil {
		t.Error("the bad response should not have headers")
	}
}

func TestDecodePrefixSuffixUrl(t *testing.T) {
	q, _ := url.ParseQuery("p=http://foo/bar/0&s=.win")
	urls, _ := DecodeUrls(q)
	if urls[0] != "http://foo/bar/0.win" {
		t.Error("Could use prefix only URL")
	}
}

func TestDecodePrefixUrl(t *testing.T) {
	q, _ := url.ParseQuery("p=http://foo/bar/0.win")
	urls, _ := DecodeUrls(q)
	if urls[0] != "http://foo/bar/0.win" {
		t.Error("Could use prefix only URL")
	}
}

func TestDecodeMultipleUrls(t *testing.T) {
	q, _ := url.ParseQuery("p=http://foo/bar/0&i=123&i=456&s=.win")
	urls, _ := DecodeUrls(q)

	if urls[0] != "http://foo/bar/0123.win" {
		t.Error("Could not rebuild URLs")
	}

	if urls[1] != "http://foo/bar/0456.win" {
		t.Error("Could not rebuild URLs")
	}
}

func TestFormat(t *testing.T) {
	res := Format("", []byte("ohai"))
	matchFormat(t, res, "data:;base64,b2hhaQ==", "No content type")

	res = Format("text/plain", []byte("ohai"))
	matchFormat(t, res, "data:text/plain;base64,b2hhaQ==", "text/plain")

	res = Format("text/plain;charset=win", []byte("ohai"))
	matchFormat(t, res, "data:text/plain;charset=win;base64,b2hhaQ==", "Could not format charset")

	res = Format("text/plain; bad=space ", []byte("ohai"))
	matchFormat(t, res, "data:text/plain;bad=space;base64,b2hhaQ==", "Could not strip spaces around the parameters")

	res = Format("text/plain;bad", []byte("ohai"))
	matchFormat(t, res, "data:text/plain;base64,b2hhaQ==", "Could not strip non-parameters")
}

func TestGet(t *testing.T) {
	srv := backendPathEcho(t)
	defer srv.Close()

	res := Get(backendClient(srv), "http://test/ohai")

	matchFormat(t, res.DataURI, "data:text/plain;charset=utf-8;base64,L29oYWk=", "Test server doesn't match")
}

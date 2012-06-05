package main 

import (
	"freeload"
	"net/http"
	"net/url"
	"log"
	"flag"
	"time"
)

var (
	proxy = flag.String("proxy", "", "URL of the proxy server for outbound requests")
	jsonRoot = flag.String("json", "/json", "Root path to deliver JSON")
	origins = flag.String("origins", "*", "Origins to allow with CORS")
	timeout = flag.Float64("timeout", 0.5, "Maximum time per origin request in seconds")
	bind = flag.String("http", ":7433", "Host:Port to listen on")
)

func main() {
	var client http.Client

	flag.Parse()

	if *proxy != "" {
		proxyURL, err := url.Parse(*proxy)
		if err != nil {
			log.Fatalf("Invalid proxy URL: %v", err)
		}
		client.Transport = &http.Transport { Proxy: http.ProxyURL(proxyURL) }
	}

	http.Handle(*jsonRoot,
		freeload.GzipHandler(
			freeload.CORS(*origins,
				freeload.LoadJSONHandler(client, time.Duration(*timeout * 1e9)))))

	log.Fatal(http.ListenAndServe(*bind, nil))
}

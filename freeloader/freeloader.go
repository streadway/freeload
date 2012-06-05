// Copyright (C) 2012 Sean Treadway <treadway@gmail.com>, SoundCloud Ltd.
// All rights reserved.  See README.md for license details.

/*
Server for the freeload package
*/
package main

import (
	"flag"
	"github.com/streadway/freeload"
	"log"
	"net/http"
	"net/url"
	"time"
)

var (
	bind     = flag.String("http", ":7433", "Host:Port to listen on")
	proxy    = flag.String("proxy", "", "URL of the proxy server for outbound requests")
	jsonRoot = flag.String("json", "/json", "Root path to deliver JSON")
	origins  = flag.String("origins", "*", "Origins to allow with CORS")
	timeout  = flag.Float64("timeout", 0.5, "Maximum time per origin request in seconds")
)

func main() {
	flag.Parse()

	var client http.Client

	if *proxy != "" {
		proxyURL, err := url.Parse(*proxy)
		if err != nil {
			log.Fatalf("Invalid proxy URL: %v", err)
		}
		client.Transport = &http.Transport{Proxy: http.ProxyURL(proxyURL)}
	}

	http.Handle(*jsonRoot,
		freeload.GzipHandler(
			freeload.CORS(*origins,
				freeload.LoadJSONHandler(client, time.Duration(*timeout*1e9)))))

	log.Fatal(http.ListenAndServe(*bind, nil))
}

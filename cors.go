// Copyright (C) 2012 Sean Treadway <treadway@gmail.com>, SoundCloud Ltd.
// All rights reserved.  See README.md for license details.

package freeload

import (
	"net/http"
)

// Wraps a handler, if you want to accept all hosts, use "*"
// Will short circuit preflight requests indicated with the OPTIONS method
func CORS(origins string, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Methods", "GET")
		w.Header().Set("Access-Control-Allow-Headers", "Accept, Authorization, Content-Type, Origin")
		w.Header().Set("Access-Control-Allow-Origin", origins)

		if r.Method == "OPTIONS" && len(r.Header["Origin"]) > 0 {
			return
		}

		next.ServeHTTP(w, r)
	})
}

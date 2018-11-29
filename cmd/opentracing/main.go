// webapp: a standalone example Negroni / Gorilla based webapp.
//
// This example demonstrates basic usage of Appdash in a Negroni / Gorilla
// based web application. The entire application is ran locally (i.e. on the
// same server) -- even the Appdash web UI.
// -------------
// code copy form https://github.com/sourcegraph/appdash/blob/master/examples/cmd/webapp-opentracing/main.go
package main

import (
	"log"
	"net"
	"net/http"
	"net/url"
	"time"

	"sourcegraph.com/sourcegraph/appdash"
	"sourcegraph.com/sourcegraph/appdash/traceapp"
)

func main() {
	// Create a recent in-memory store, evicting data after 20s.
	//
	// The store defines where information about traces (i.e. spans and
	// annotations) will be stored during the lifetime of the application. This
	// application uses a MemoryStore store wrapped by a RecentStore with an
	// eviction time of 20s (i.e. all data after 20s is deleted from memory).
	memStore := appdash.NewMemoryStore()
	store := &appdash.RecentStore{
		MinEvictAge: 60 * time.Second,
		DeleteStore: memStore,
	}

	l, err := net.Listen("tcp", ":3001")
	if err != nil {
		log.Fatal(err)
	}
	proto := "plaintext TCP (no security)"
	log.Printf("appdash collector listening on %s (%s)", ":3001", proto)
	cs := appdash.NewServer(l, appdash.NewLocalCollector(store))
	go cs.Start()

	// Start the Appdash web UI on port 8700.
	//
	// This is the actual Appdash web UI -- usable as a Go package itself, We
	// embed it directly into our application such that visiting the web server
	// on HTTP port 8700 will bring us to the web UI, displaying information
	// about this specific web-server (another alternative would be to connect
	// to a centralized Appdash collection server).
	url, err := url.Parse("http://localhost:8700")
	if err != nil {
		log.Fatal(err)
	}
	tapp, err := traceapp.New(nil, url)
	if err != nil {
		log.Fatal(err)
	}
	tapp.Store = store
	tapp.Queryer = memStore
	log.Println("Appdash web UI running on HTTP :8700")
	log.Fatal(http.ListenAndServe(":8700", tapp))
}

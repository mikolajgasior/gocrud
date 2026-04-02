package handler

import (
	"net/http"
)

func AddHandler(serveMux *http.ServeMux, uri string, handler http.HandlerFunc) {
	// if uri is empty, do nothing - we cannot allow overriding '/' by one handler
	if uri == "" {
		return
	}

	// uri should not have a trailing slash
	subMux := http.NewServeMux()
	subMux.HandleFunc("/", handler)
	serveMux.Handle(uri+"/", http.StripPrefix(uri, subMux))
}

func AddHandlers(serveMux *http.ServeMux, handlers map[string]http.HandlerFunc) {
	for uri, handler := range handlers {
		serveMux.HandleFunc(uri, handler)
	}
}

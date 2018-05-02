package tool

import (
	"net/http"
)

type JsonHandlerFunc http.HandlerFunc
type HtmlHandlerFunc http.HandlerFunc

func (f JsonHandlerFunc) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if shouldFailPurposeHeader(w, r) {
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	setCommonHeaders(w)
	f(w, r)
}

func (f HtmlHandlerFunc) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if shouldFailPurposeHeader(w, r) {
		return
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	setCommonHeaders(w)
	f(w, r)
}

func setCommonHeaders(w http.ResponseWriter) {
	// Prevent caching of any of the requests so that we can use GET for things like /reboot
	w.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")
	w.Header().Set("Pragma", "no-cache")
	w.Header().Set("Expires", "0")

	// Allow any origin to hit the miner. Potential danger here, but if someone has access
	// to your miner network, then you've got bigger problems anyway since they have a default
	// root password (and some machines even have telnet enabled).
	w.Header().Set("Access-Control-Allow-Origin", "*")
}

func shouldFailPurposeHeader(w http.ResponseWriter, r *http.Request) bool {
	// https://www.google.com/search?q=X-Purpose%3A+preview
	if len(r.Header.Get("X-Purpose")) > 0 {
		w.Header().Set("Content-Type", "text/plain")
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("No preview allowed"))
		return true
	}
	return false
}
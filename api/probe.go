package api

import "net/http"

func (a Api) probe(w http.ResponseWriter, r *http.Request) {
	_, _ = w.Write([]byte("ok"))
}

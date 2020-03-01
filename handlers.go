package router

import "net/http"

type MethodNotAllowed struct {
}

func (m *MethodNotAllowed) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	w.WriteHeader(http.StatusMethodNotAllowed)
}


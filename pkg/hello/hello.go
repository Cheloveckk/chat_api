package hello

import (
	"net/http"
)

func hello(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("test", "hello")
	w.Write([]byte("hello"))

}
func HelloHandle(handler *http.ServeMux) {
	handler.HandleFunc("/hello", hello)
}

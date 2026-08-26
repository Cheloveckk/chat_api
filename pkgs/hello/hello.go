package hello

import (
	"net/http"
	"strconv"
)

var Counter = 0

func hello(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("test", "hello")
	w.Write([]byte("hello" + strconv.Itoa(Counter)))
	Counter += 1

}
func HelloHandle(handler *http.ServeMux) {
	handler.HandleFunc("/hello", hello)
}

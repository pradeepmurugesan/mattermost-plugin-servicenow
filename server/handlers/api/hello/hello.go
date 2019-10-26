package hello

import (
	"fmt"
	"net/http"
)

// SayHello hello world
func SayHello(writer http.ResponseWriter, request *http.Request) {
	fmt.Fprint(writer, "Hello, world!!!!")
}

package wxhelper

import (
	"io"
	"log"
	"net/http"
	"os"
)

func HandlerFunc(w http.ResponseWriter, response *http.Request) {
	log.Println("received message")
	io.Copy(os.Stdout, response.Body)
}

package main

import (
	"net/http"

	"github.com/OkutaniDaichi0106/mcp-go/mcp"
)

func main() {
	server := mcp.NewServer("http-server")

	http.HandleFunc("/mcp", func(w http.ResponseWriter, r *http.Request) {
		sess, err := server.AcceptHTTP(w, r)
		if err != nil {
			return
		}
		defer sess.Close()

		//
	})

	http.ListenAndServe(":8080", nil)
}

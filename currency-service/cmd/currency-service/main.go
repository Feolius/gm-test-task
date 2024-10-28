package main

import (
	_ "github.com/go-sql-driver/mysql"
	"log"
	"net/http"
)

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("hello"))
	})
	err := http.ListenAndServe(":3001", mux)
	log.Fatal(err)
}

package main

import (
	"flag"
	"fmt"
	"net/http"
)

func main() {
	var port int
	flag.IntVar(&port, "p", 3000, "Client port")
	flag.Parse()
	
	fmt.Printf("Client listening on port %v\n", port)
	http.ListenAndServe(fmt.Sprintf(":%v", port), http.FileServer(http.Dir("public")))
}
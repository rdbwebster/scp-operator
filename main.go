package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
)

var kubeconfig string

func init() {
	flag.StringVar(&kubeconfig, "kubeconfig", "./.kube/config", "Path to kubeconfig")
	flag.Parse()
}

func main() {

	fmt.Fprintf(os.Stderr, "Listening on  :8080 \n")
	fmt.Fprintf(os.Stderr, "kubeconfig %s \n", kubeconfig)
	router := NewRouter()
	log.Fatal(http.ListenAndServe(":8080", router))
}

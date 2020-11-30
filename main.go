package main

import (
        "log"
        "net/http"
        "fmt"
        "os"
)

func main() {

        fmt.Fprintf(os.Stderr, "Listening on  :8080");
	router := NewRouter()
        log.Fatal(http.ListenAndServe(":8080", router))
}

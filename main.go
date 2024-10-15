package main

import (
        "net/http"
        "fmt"

        "github.com/prometheus/client_golang/prometheus/promhttp"
)

func main() {

        http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
        fmt.Fprintf(w, "Hello!")
        })

        http.Handle("/metrics", promhttp.Handler())
        http.ListenAndServe(":2112", nil)
}


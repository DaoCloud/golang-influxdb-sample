package main

import (
    "log"

    "github.com/influxdb/influxdb/client"
)

var ic *client.Client

func connectStorage() {
    host, err := url.Parse(fmt.Sprintf("http://%s:%d", "localhost", 8086))
    if err != nil {
        log.Fatal(err)
    }

    con, err := client.NewClient(client.Config{URL: *host})
    if err != nil {
        log.Fatal(err)
    }

    dur, ver, err := con.Ping()
    if err != nil {
        log.Fatal(err)
    }

    log.Printf("Happy as a hippo! %v, %s", dur, ver)
    ic = con
}

func hello(res http.ResponseWriter, req *http.Request) {
    res.Write(fmt.Sprintf("Hello World, %d!", 1234))
}

func main() {
    connectStorage()
    http.HandleFunc("/", hello)

    log.Println("Start listening...")
    if err := http.ListenAndServe(":80", nil); err != nil {
        log.Fatal(err)
    }
}



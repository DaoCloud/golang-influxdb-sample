package main

import (
    "os"
    "testing"
)

func Test_FetchFromInflux(t *testing.T) {
    host = os.Getenv("INFLUXDB_PORT_8086_TCP_ADDR")
    port = os.Getenv("INFLUXDB_PORT_8086_TCP_PORT")
    db = os.Getenv("INFLUXDB_INSTANCE")
    user = os.Getenv("INFLUXDB_USERNAME")
    password = os.Getenv("INFLUXDB_PASSWORD")

    connect()
    create()
    insert()

    m := query()
    if m["color"] != "red" || m["shape"] != "circle" {
        t.Error("expect red circle")
    }
}

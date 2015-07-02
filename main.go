package main

import (
    "log"
    "fmt"
    "time"
    "os"
    "net/http"
    "net/url"

    "github.com/influxdb/influxdb/client"
)

var ic *client.Client

var host string
var port string
var db string
var user string
var password string

func connect() {
    u, err := url.Parse(fmt.Sprintf("http://%s:%s", host, port))
    if err != nil {
        log.Fatal(err)
    }

    ic, err = client.NewClient(client.Config{URL: *u})
    if err != nil {
        log.Fatal(err)
    }

    if _, _, err := ic.Ping(); err != nil {
        log.Fatal(err)
    }

    ic.SetAuth(user, password)
}

func create() {
    // Workaround, since daocloud influxdb haven't privision an instance
    // create the db instance here
    q := client.Query{
        Command:  fmt.Sprintf("create database %s", db),
        Database: db,
    }

    // ignore the error of existing database
    ic.Query(q)
}


func insert() {
    var (
        shapes     = []string{"circle", "rectangle", "square", "triangle"}
        colors     = []string{"red",    "blue",      "green",  "yellow"  }
        sampleSize = 4
        pts        = make([]client.Point, sampleSize)
    )

    for i := 0; i < sampleSize; i++ {
        pts[i] = client.Point{
            Measurement: "shapes",
            Tags: map[string]string{
                "color": colors[i],
                "shape": shapes[i],
            },
            Fields: map[string]interface{}{
                "value": i,
            },
            Time: time.Now(),
        }
    }

    bps := client.BatchPoints{
        Points:          pts,
        Database:        db,
        RetentionPolicy: "default",
    }

    _, err := ic.Write(bps)
    if err != nil {
        log.Println("Insert data error:")
        log.Fatal(err)
    }
}

func query() map[string]string {
    q := client.Query{
        Command:  "select * from shapes where value = 0",
        Database: db,
    }

    response, err := ic.Query(q)
    if err != nil {
        log.Println("Error, ", err)
        return nil
    }

    if response.Error() != nil {
        log.Println("Response error, ", response.Error())
        return nil
    }

    result := response.Results[0]
    if result.Err != nil {
        log.Println("Result error, ", result.Err)
        return nil
    }

    serie := result.Series[0]
    if serie.Err != nil {
        log.Println("Serie error, ", serie.Err)
        return nil
    }

    return serie.Tags
}

func hello(res http.ResponseWriter, req *http.Request) {
    m := query()
    res.Write([]byte(fmt.Sprintf("The first shape is %s %s!", m["color"], m["shape"])))
}

func main() {
    host = os.Getenv("INFLUXDB_PORT_8086_TCP_ADDR")
    port = os.Getenv("INFLUXDB_PORT_8086_TCP_PORT")
    db = os.Getenv("INFLUXDB_INSTANCE")
    user = os.Getenv("INFLUXDB_USERNAME")
    password = os.Getenv("INFLUXDB_PASSWORD")

    // workaround, daocloud influxdb have not privision db instance
    if len(db) == 0 {
        db = "mydb"
    }

    connect()
    log.Println("Successfully connect to influxdb ...")

    // prepare data
    create()
    insert()

    http.HandleFunc("/", hello)

    log.Println("Start listening...")
    if err := http.ListenAndServe(":80", nil); err != nil {
        log.Fatal(err)
    }
}

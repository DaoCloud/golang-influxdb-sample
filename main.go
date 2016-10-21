package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/influxdata/influxdb/client/v2"
)

var ic client.Client

var host string
var port string
var db string
var user string
var password string

func connect() {
	u := fmt.Sprintf("http://%s:%s", host, port)
	var err error
	ic, err = client.NewHTTPClient(client.HTTPConfig{Addr: u, Username: user, Password: password})
	if err != nil {
		log.Fatal(err)
	}
	if _, _, err := ic.Ping(time.Second * 5); err != nil {
		log.Fatal(err)
	}

}

func insert() {
	// Create a new point batch
	bp, err := client.NewBatchPoints(client.BatchPointsConfig{
		Database:  db,
		Precision: "s",
	})

	if err != nil {
		log.Fatalln("Error: ", err)
	}

	// Create a point and add to batch
	tags := map[string]string{"cpu": "cpu-total"}
	fields := map[string]interface{}{
		"idle":   10.1,
		"system": 53.3,
		"user":   46.6,
	}
	pt, err := client.NewPoint("cpu_usage", tags, fields, time.Now())

	if err != nil {
		log.Fatalln("Error: ", err)
	}

	bp.AddPoint(pt)

	// Write the batch
	ic.Write(bp)
}

func queryDB(cmd string) (res []client.Result, err error) {
	q := client.Query{
		Command:  cmd,
		Database: db,
	}
	if response, err := ic.Query(q); err == nil {
		if response.Error() != nil {
			return res, response.Error()
		}
		res = response.Results
	} else {
		return res, err
	}
	return res, nil
}

func hello(res http.ResponseWriter, req *http.Request) {
	q := fmt.Sprintf("SELECT * FROM %s LIMIT %d", "cpu_usage", 20)

	m, err := queryDB(q)
	if err != nil {
		fmt.Println("query failed!err:=%v", err)
	}
	res.Write([]byte(fmt.Sprintf("The first cpu_measure is %v!", m[0].Series[0])))
}

func main() {
	host = os.Getenv("INFLUXDB_PORT_8086_TCP_ADDR")
	port = os.Getenv("INFLUXDB_PORT_8086_TCP_PORT")
	db = os.Getenv("INFLUXDB_INSTANCE")
	user = os.Getenv("INFLUXDB_USERNAME")
	password = os.Getenv("INFLUXDB_PASSWORD")
	connect()
	log.Println("Successfully connect to influxdb ...")
	// workaround, daocloud influxdb have not privision db instance
	if len(db) == 0 {
		db = "mydb"
		// prepare data
		_, err := queryDB(fmt.Sprintf("CREATE DATABASE %s", db))
		if err != nil {
			log.Fatal(err)
		}
	}

	insert()

	http.HandleFunc("/", hello)

	log.Println("Start listening...")
	if err := http.ListenAndServe(":80", nil); err != nil {
		log.Fatal(err)
	}
}

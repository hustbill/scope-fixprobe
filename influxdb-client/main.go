package main



/*
* Read data from InfluxDB, and render it
*/


import (
	"encoding/json"
	"os"
	"net/http"
	"net/http/httptest"
	"net/url"
	//"strings"
	"time"
	"github.com/Sirupsen/logrus"
	"github.com/influxdata/influxdb/client"

)

const db = "canal-monitoring"

var log = logrus.New()
var command = "select * from ep_stats"

func main() {
	log.Out = os.Stdout
	log.Info("Test New Client to InfluxDB")
	TestNewClient()
	TestClient_Ping()
	TestClient_Query(db, command)
}


func TestNewClient() {
	config := client.Config{}
	_, err := client.NewClient(config)
	if err != nil {
		log.Fatalf("unexpected error.  expected %v, actual %v", nil, err)
	} else {
		log.Info("Connect to InfluxDB successfully")
	}
}

// helper functions
func emptyTestServer() *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(50 * time.Millisecond)
		w.Header().Set("X-Influxdb-Version", "1.2.2")
		return
	}))
}


func TestClient_Ping() {
	log.Info("TestClient_Ping")
	ts := emptyTestServer()
	defer ts.Close()

	u, _ := url.Parse(ts.URL)

	u.Host = "localhost:8086"  // set to InfluxDB default port : 8086

	config := client.Config{URL: *u}

	log.Info("u =", u)


	c, err := client.NewClient(config)
	if err != nil {
		log.Fatalf("unexpected error.  expected %v, actual %v", nil, err)
	}
	d, version, err := c.Ping()
	log.Info("version: ", version)
	log.Info("d: ", d)
	if err != nil {
		log.Fatalf("unexpected error.  expected %v, actual %v", nil, err)
	}
	if d.Nanoseconds() == 0 {
		log.Fatalf("expected a duration greater than zero.  actual %v", d.Nanoseconds())
	}
	if version != "1.2.2" {
		log.Fatalf("unexpected version.  expected %s,  actual %v", "1.2.2", version)
	}
}

func TestClient_Query(db string, command string) {
	log.Info("TestClient_Query")
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var data client.Response
		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode(data)
	}))
	defer ts.Close()

	u, _ := url.Parse(ts.URL)

	u.Host = "localhost:8086"  // set to InfluxDB default port : 8086

	config := client.Config{URL: *u}
	c, err := client.NewClient(config)
	if err != nil {
		log.Fatalf("unexpected error.  expected %v, actual %v", nil, err)
	}

	query := client.Query{}
	query.Database = db
	query.Command = command


	res, err := c.Query(query)
	log.Info("response: ", res)

	if err != nil {
		log.Fatalf("unexpected error.  expected %v, actual %v", nil, err)
	}
}

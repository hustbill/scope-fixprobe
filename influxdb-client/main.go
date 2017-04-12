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

	"fmt"
	"strconv"
)

const db = "ican"

var command = "select * from p2prxbyte"

var log = logrus.New()


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
	command = "select * from p2prxbyte, p2prxpkt"

	query.Command = command


	resp, err := c.Query(query)
	//log.Info("response: ", resp)

	if err != nil {
		log.Fatalf("unexpected error.  expected %v, actual %v", nil, err)
	}



	var myData [][]interface{} = make([][]interface{}, len(resp.Results[0].Series[0].Values))
	for i, d := range resp.Results[0].Series[0].Values {
		myData[i] = d
	}
	//
	//// row : p2prxpkt,
	//// host=10.145.240.148,
	//// spod_name=user-468431046-6wr62,
	//// spod_namespace=user-468431046-6wr62,
	//// dpod_name=weave-net-8sww0,
	//// dpod_namespace=kube-system,
	//// src=10.32.0.14,
	//// dst=10.145.240.148
	//// value=41
	//
	log.Printf("Result :")
	fmt.Println("", myData[0]) //first element in slice
	//
	//var bandwidth = myData[0][8]
	//log.Printf("bandwidth is %s", bandwidth)
	//log.Println("---------------------------")
	//fmt.Println("time: ", myData[0][0])
	//fmt.Println("dpod_name: ", myData[0][1])
	//fmt.Println("dpod_namespace: ", myData[0][2])
	//fmt.Println("dst: ", myData[0][3])
	//fmt.Println("host: ", myData[0][4])
	//fmt.Println("spod_name: ", myData[0][5])
	//fmt.Println("spod_namespace: ", myData[0][6])
	//fmt.Println("src: ", myData[0][7])
	//fmt.Printf("value: %s", myData[0][8])


	var testData [][]interface{} = make([][]interface{}, len(resp.Results[0].Series[1].Values))
	for i, d := range resp.Results[0].Series[1].Values {
		testData[i] = d
	}

	log.Printf("\nPacket :\n")
	fmt.Println("", testData[0]) //first element in slice
	fmt.Printf("col: %s\n", testData[0][0])
	fmt.Printf("col: %s\n", testData[0][1])
	fmt.Printf("col: %s\n", testData[0][2])
	fmt.Printf("col: %s\n", testData[0][3])
	fmt.Printf("col: %s\n", testData[0][4])

	fmt.Printf("value: %s\n", testData[0][8])


	log.Println("==========================")

	//bandwidth, err = resp.Results[0].Series[0].Values[0][8].(json.Number).Int64()
	//packet, err := resp.Results[0].Series[1].Values[0][8].(json.Number).Int64()
	//fmt.Printf("bandwidth = %d\t", bandwidth)
	//fmt.Printf("packet = %d", packet)

	result := [3]string{"-", "-", "-"}

	pod, ok := resp.Results[0].Series[0].Values[0][1].(string)
	if (ok) {
		log.Printf("spod_name : %s ", pod)
		result[0] = pod
	} else {
		log.Debugf("failed read data from InfluxDB, pod error", ok)
	}




	bandwidth, err := resp.Results[0].Series[0].Values[0][8].(json.Number).Int64()
	packet, err := resp.Results[0].Series[1].Values[0][8].(json.Number).Int64()

	result[1] = strconv.FormatInt(bandwidth, 10)
	result[2] = strconv.FormatInt(packet, 10)

	if err != nil {
		log.Errorf("failed read data from InfluxDB: %v", err)

	}

	log.Println("result:  ", result)


}

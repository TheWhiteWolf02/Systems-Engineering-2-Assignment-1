package main

import (
	"fmt"
	"net/http"
	"log"
	"io/ioutil"
	"github.com/go-zookeeper/zk"
	"encoding/json"
	"bytes"
	"time"
	"os"
)

const DEFAULT_PORT = "7000"

func handler(w http.ResponseWriter, req *http.Request) {
	fmt.Printf("Gserve request handle for method %s\n", req.Method)
	if req.Method == "POST" {
		// unmarshal, encode, marshal, convert to byte and send
		body, _ := ioutil.ReadAll(req.Body)
		var row RowsType

		json.Unmarshal(body, &row)

		if len(row.Row[0].Key) == 0 {
			// empty POST request
			return
		}

		fmt.Printf("log: re.Body: %s :: body: %s :: row: %s :: row.key: %s\n", req.Body, body, row, row.Row[0].Key)

		url := "http://hbase:8080/se2:library/" + row.Row[0].Key
		encRowsType := row.encode()
		byteArray, _ := json.Marshal(encRowsType)

		http.Post(url, "application/json", bytes.NewBuffer(byteArray))
	} else if req.Method == "GET" {
		// TODO:
		// go templating html
		// put -> get

		// os.getEnv(key) for gserve1 or gserve2
	}
}

func main() {
	time.Sleep(20*time.Second)
	fmt.Printf("Gserve version started: %s\n", os.Getenv("version"))

	conn, _, _ := zk.Connect([]string{"zookeeper:2181"}, time.Second)
	resp, _ := conn.Create("/" + os.Getenv("version"), nil, 1, []zk.ACL{{Perms: 31, Scheme: "world", ID: "anyone"}})

  	fmt.Printf("Response: %s\n", resp)

	http.HandleFunc("/library", handler)
	if err := http.ListenAndServe(":"+DEFAULT_PORT, nil); err != nil {
    	log.Fatal("ListenAndServe: ", err)
	}
}
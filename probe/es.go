package probe

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"strings"
	"time"
)

var ESUrl = strings.Replace(os.Getenv("ELASTIC_PORT"), "tcp", "http", 1)

func postStats(v interface{}) {
	timestamp := time.Now().Format("2006.01.02")
	var postData []byte
	w := bytes.NewBuffer(postData)
	json.NewEncoder(w).Encode(v)
	resp, err := http.Post(ESUrl+"/logstash-"+timestamp+"/monitor/", "application/json", w)
	if err != nil {
		log.Printf("postStats %v to %v:%v", v, ESUrl, err)
	} else {

		resp.Body.Close()
	}
}

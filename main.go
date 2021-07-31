package main

import (
	"bytes"
	"encoding/json"
	"log"
	"math/rand"
	"net/http"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

// we create a new custom metric of type counter
var userStatus = prometheus.NewCounterVec(
	prometheus.CounterOpts{
		Name: "http_request_get_user_status_count", // metric name
		Help: "Count of status returned by user.",
	},
	[]string{"user", "status"}, // labels
)

func init() {
	// we need to register the counter so prometheus can collect this metric
	prometheus.MustRegister(userStatus)
}

type MyRequest struct {
	User string
}

// the server will retrieve the user from the body, and randomly generate a status to return
func server(w http.ResponseWriter, r *http.Request) {
	var status string
	var user string
	defer func() {
		userStatus.WithLabelValues(user, status).Inc()
	}()
	var mr MyRequest
	json.NewDecoder(r.Body).Decode(&mr)

	if rand.Float32() > 0.8 {
		status = "4xx"
	} else {
		status = "2xx"
	}

	user = mr.User
	log.Println(user, status)
	w.Write([]byte(status))
}

// the producer will randomly select a user from a pool of users and send it to the server
func producer() {
	userPool := []string{"bob", "alice", "jack"}
	for {
		postBody, _ := json.Marshal(MyRequest{
			User: userPool[rand.Intn(len(userPool))],
		})
		requestBody := bytes.NewBuffer(postBody)
		http.Post("http://localhost:8080", "application/json", requestBody)
		time.Sleep(time.Second * 2)
	}
}

func main() {
	// the producer runs on its own goroutine
	go producer()

	http.Handle("/metrics", promhttp.Handler())
	http.HandleFunc("/", server)

	http.ListenAndServe(":8080", nil)
}

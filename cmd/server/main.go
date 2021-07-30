package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

// create a new counter vector
var userStatus = prometheus.NewCounterVec(
	prometheus.CounterOpts{
		Name: "http_request_get_user_status_count", // metric name
		Help: "Count of status returned by user.",
	},
	[]string{"user", "status"}, // labels
)

func init() {
	prometheus.MustRegister(userStatus)
}

type MyRequest struct {
	User   string
	Status string
}

func myhandler(w http.ResponseWriter, r *http.Request) {
	var status string
	var user string
	defer func() {
		// increment the counter on defer func
		userStatus.WithLabelValues(user, status).Inc()
	}()
	var mr MyRequest
	json.NewDecoder(r.Body).Decode(&mr)

	status = mr.Status
	user = mr.User
	log.Println(user, status)
	w.Write([]byte(status))

}

func main() {

	http.Handle("/metrics", promhttp.Handler())
	http.HandleFunc("/myhandler", myhandler)
	println("listening..")
	http.ListenAndServe(":8080", nil)
}

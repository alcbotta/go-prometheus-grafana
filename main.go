package main

import (
	"encoding/json"
	"net/http"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

// create a new counter vector
var userStatus = prometheus.NewCounterVec(
	prometheus.CounterOpts{
		Name: "http_request_get_user_status_count", // metric name
		Help: "Number status returned by user.",
	},
	[]string{"user", "status"}, // labels
)

func init() {
	prometheus.MustRegister(userStatus)
}



func myhandler(w http.ResponseWriter, r *http.Request) {
	var status string
	defer func() {
        // increment the counter on defer func
		getBookCounter.WithLabelValues(status).Inc()
	}()

	books, err := getBooks()
	if err != nil {
		status = "error"
		w.Write([]byte("something's wrong: " + err.Error()))
		return
	}

	resp, err := json.Marshal(books)
	if err != nil {
		status = "error"
		w.Write([]byte("something's wrong: " + err.Error()))
		return
	}

	status = "success"
	w.Write(resp)
}

func main() {

	http.Handle("/metrics", promhttp.Handler())
	http.HandleFunc("/myhandler", myhandler)
	println("listening..")
	http.ListenAndServe(":8080", nil)
}
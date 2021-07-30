package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"time"
)

func main() {
	userPool := []string{"bob", "alice", "jack"}
	statusPool := []string{"200", "400", "500"}
	for {
		//Encode the data
		postBody, _ := json.Marshal(map[string]string{
			"user":   userPool[rand.Intn(len(userPool))],
			"status": statusPool[rand.Intn(len(statusPool))],
		})

		requestBody := bytes.NewBuffer(postBody)
		_, err := http.Post("http://localhost:8080/myhandler", "application/json", requestBody)
		if err != nil {
			fmt.Println(err.Error())
		}

		time.Sleep(time.Second * 2)
	}

}

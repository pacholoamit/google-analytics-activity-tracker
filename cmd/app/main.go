package main

import (
	"io/ioutil"
	"log"
	"net/http"
)

func main() {

	makeGetRequest("https://jsonplaceholder.typicode.com/todos/1")
}

func makeGetRequest(url string) string {
	resp, err := http.Get(url)

	if err != nil {
		log.Fatalln(err)
	}

	body, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		log.Fatalln(err)
	}

	sb := string(body)

	return sb
}

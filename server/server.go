package server

// package main

import (
	// "errors"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
)

func RunServer() {

	router := mux.NewRouter() //.StrictSlash(true)

	watchers, err := ParseServices("config/services.json")
	if err != nil {
		log.Fatal(err)
	}
	StartWatchers(watchers)

	handlers := map[string]func() (string, error){
		"load": func() (string, error) {
			return toJSON(watchers), nil
		},
		"update": func() (string, error) {
			result := []map[string]interface{}{}
			currentTime := time.Now()
			for _, watcher := range watchers {
				result = append(result, watcher.Update(currentTime))
			}
			return toJSON(result), nil
		},
	}
	for route, handler := range handlers {
		router.HandleFunc("/"+route, wrap(handler))
	}

	router.HandleFunc("/client/config.json", wrap(curry(file, "config/client.json")))
	router.PathPrefix("/client/").Handler(http.FileServer(http.Dir(".")))

	http.Handle("/", router)

	log.Fatal(http.ListenAndServe(":54007", router))
}
func curry(f func(string) (string, error), param string) func() (string, error) {
	return func() (string, error) {
		return f(param)
	}
}

func file(filename string) (content string, err error) {
	bytes, err := ioutil.ReadFile(filename)
	return string(bytes), err
}

func wrap(f func() (string, error)) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, req *http.Request) {
		result, err := f()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		_, err = fmt.Fprint(w, result)
		if err != nil {
			// If you can't write content can you write an error? :D
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}
}

func ParseServices(location string) (result []*Watcher, err error) {
	// type Services struct {
	// 	Http []*Watcher
	// 	Rss  []*Watcher
	// }

	// var services Services
	var services map[string][]*Watcher
	bytes, err := ioutil.ReadFile(location)
	// bytes, err := ioutil.ReadFile("../config/services.json")
	if err != nil {
		return
	}
	// return string(bytes), err
	err = json.Unmarshal(bytes, &services)
	if err != nil {
		return
	}
	result = []*Watcher{}
	for serviceType, servicesOfType := range services {
		for _, service := range servicesOfType {
			for _, watchPage := range service.Pages {
				watchPage.Type = serviceType
			}
		}
		result = append(result, servicesOfType...)
	}
	return
	// a, _ := json.MarshalIndent(services, "", "  ")
	// fmt.Println(string(a))
}

func StartWatchers(watchers []*Watcher) {
	for i, watcher := range watchers {
		watcher.Init()
		fmt.Println("Watcher " + strconv.Itoa(i) + " started")
	}
}

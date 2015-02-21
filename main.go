package main

import (
	// "./cef2go/Release"
	// "encoding/json"

	"./server"
	"fmt"
	"log"
	"os/exec"
	"runtime"
)

func RunUI() {
	fmt.Println("Starting CEF browser - it will take about 10s")

	var path string
	if runtime.GOOS == "windows" {
		path = "build.bat"
	} else {
		path = "make"
	}

	cmd := exec.Cmd{Path: path, Dir: "./cef2go"}

	// cmd.Stdin = strings.NewReader("some input")
	// var out bytes.Buffer
	// cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		fmt.Println("Could not start CEF")
		log.Fatal(err)
	}

}

func main() {
	go server.RunServer()
	RunUI()

	// services, _ := server.ParseServices("config/services.json")
	// a, _ := json.MarshalIndent(services, "", "  ")
	// fmt.Println(string(a))
}

package main

import (
	"encoding/csv"
	_ "encoding/csv"
	"fmt"
	"io"
	_ "io"
	"log"
	_ "log"
	"net/http"
	"os"
	"strings"
	_ "strings"
	"time"
)

func main() {
	start := time.Now()
	arg := os.Args[1]
	csvfile, err := os.Open(arg)
	if err != nil {
		log.Fatalln("Couldn't open the csv file", err)
	}
	// Parse the file
	r := csv.NewReader(csvfile)
	// Iterate through the records
	data := [][]string{}
	for {
		// Read each record from csv
		domain, err := r.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatal(err)
		}
		url := domain[0]
		server := fetch(url)
		row := []string{url, server}
		data = append(data, row)
	}

	//write to output
	outputname := strings.Split(arg, ".")[0] + "_" + "output_a.csv"
	outputfile, err := os.Create(outputname)
	w := csv.NewWriter(outputfile)
	w.WriteAll(data)
	if err := w.Error(); err != nil {
		log.Fatalln("error writing csv:", err)
	}
	secs := time.Since(start).Seconds()
	fmt.Println("elapsed time:", secs)
}

func fetch(url string) (server string) {
	fmt.Printf("domain: " + url + "\n")
	server = "nil"
	timeout := time.Duration(10 * time.Second)
	client := http.Client{
		Timeout: timeout,
	}
	resp, err := client.Get("http://" + url)
	if err != nil {
		fmt.Println(err)
		return server
	}
	server = resp.Header.Get("Server")
	if (server == "") {
		server = "nil"
	}
	fmt.Printf("server: " + server + "\n")
	fmt.Println()
	return server
}

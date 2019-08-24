package main

import (
	"encoding/csv"
	_ "encoding/csv"
	"errors"
	"fmt"
	"io"
	_ "io"
	"log"
	_ "log"
	"net/http"
	"os"
	"strconv"
	"strings"
	_ "strings"
	"time"
)

//Have to say my implementation of parallelisation does not speedup that much :( please go easy on me, I'm new to Golang.
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
	ch := make(chan string, 3)
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
		go fetch3(url, ch)
		server := <-ch
		count := <-ch
		upgrade := <-ch
		row := []string{url, server, upgrade, count}
		data = append(data, row)
	}

	//write to output
	outputname := strings.Split(arg, ".")[0] + "_" + "output_c.csv"
	outputfile, err := os.Create(outputname)
	w := csv.NewWriter(outputfile)
	w.WriteAll(data)
	if err := w.Error(); err != nil {
		log.Fatalln("error writing csv:", err)
	}
	secs := time.Since(start).Seconds()
	fmt.Println("elapsed time:", secs)
}

func fetch3(url string, ch chan<- string) error{
	fmt.Printf("domain: " + url + "\n")
	timeout := time.Duration(15 * time.Second)
	upgrade := "false"
	count := "nil"
	count_int := 0
	server := "nil"
	client := &http.Client{
		//Transport: LogRedirects{},
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			fmt.Println("redirect", req.URL)
			if (strings.Contains(req.URL.String(), "https")) {
				//upgrade = "true"
				count_int++
			}
			if len(via) >= 10 {
				return errors.New("stopped after 10 redirects")
			}
			return nil
		},
		Timeout: timeout,
	}
	resp, err := client.Get("http://" + url)
	if err != nil {
		fmt.Println(err)
		fmt.Println()
		ch <- server
		ch <- count
		ch <- upgrade
		return err
	}
	fmt.Println("Count: ", count_int)
	server = resp.Header.Get("Server")
	fmt.Println("server: " + server)
	fmt.Println("Final URL is", resp.Request.URL)
	if (strings.Contains(resp.Request.URL.String(), "https")) {
		upgrade = "true"
		count = strconv.Itoa(count_int)
	}
	fmt.Println("if upgrade: ", upgrade)
	fmt.Println()
	ch <- server
	ch <- count
	ch <- upgrade
	return nil
}

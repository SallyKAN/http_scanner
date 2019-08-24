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

type LogRedirects struct {
	Transport http.RoundTripper
}

func (l LogRedirects) RoundTrip(req *http.Request) (resp *http.Response, err error) {
	t := l.Transport
	if t == nil {
		t = http.DefaultTransport
	}
	resp, err = t.RoundTrip(req)
	if err != nil {
		return
	}
	switch resp.StatusCode {
	case http.StatusMovedPermanently, http.StatusFound, http.StatusSeeOther, http.StatusTemporaryRedirect:
		fmt.Println("Request for", req.URL, "redirected with status", resp.StatusCode)
	}
	return
}

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
		server, count, upgrade := fetch2(url)
		row := []string{url, server, upgrade, count}
		data = append(data, row)
	}

	//outputname := arg + "_" + "output_b.csv"
	outputname := strings.Split(arg, ".")[0] + "_" + "output_b.csv"
	outputfile, err := os.Create(outputname)
	w := csv.NewWriter(outputfile)
	w.WriteAll(data)
	if err := w.Error(); err != nil {
		log.Fatalln("error writing csv:", err)
	}
	secs := time.Since(start).Seconds()
	fmt.Println("elapsed time: ", secs)
}

func fetch2(url string) (server string, count string, upgrade string) {
	fmt.Printf("domain: " + url + "\n")
	timeout := time.Duration(15 * time.Second)
	server = "nil"
	upgrade = "false"
	count = "nil"
	count_int := 0
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
		return server, count, upgrade
	}
	fmt.Println("Count: ", count_int)
	server = resp.Header.Get("Server")
	if(server == ""){
		server = "nil"
	}
	fmt.Println("server: " + server)
	fmt.Println("Final URL is", resp.Request.URL)
	if (strings.Contains(resp.Request.URL.String(), "https")) {
		upgrade = "true"
		count = strconv.Itoa(count_int)
	}
	fmt.Println("if upgrade: ", upgrade)
	fmt.Println()
	return server, count, upgrade
}

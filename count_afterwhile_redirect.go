package main

import (
	"encoding/csv"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
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
	ch := make(chan string, 1)
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
		//fmt.Println(url2)
		//url := "snapinspect.com"
		go fetch5(url, ch)
		afterwhile_redirect := <-ch
		row := []string{url, afterwhile_redirect}
		data = append(data, row)
	}

	//write to output
	//outputname := arg + "_" + "output_c.csv"
	outputname := strings.Split(arg, ".")[0] + "_" + "count_afterwhile_redirect.csv"
	outputfile, err := os.Create(outputname)
	w := csv.NewWriter(outputfile)
	w.WriteAll(data)
	if err := w.Error(); err != nil {
		log.Fatalln("error writing csv:", err)
	}
	secs := time.Since(start).Seconds()
	fmt.Println("elapsed time:", secs)
}

func fetch5(url string, ch chan<- string) error{
	fmt.Printf("domain: " + url + "\n")
	timeout := time.Duration(15 * time.Second)
	count_int := 0
	// "HTTPS redirect occurred after a while" means that there must be more than one redirect(2 redirects at least), 
	// also the first redirect must not upgrade to HTTPS and the final redirect must upgrade to HTTPS
	afterwhile_redirect := "false"
	first_redirect_upgrade := false
	upgrade := false
	client := &http.Client{
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			fmt.Println("redirect", req.URL)
			if (len(via) == 1 && strings.Contains(req.URL.String(), "https")) {
				first_redirect_upgrade = true
			}
			count_int = len(via)
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
		ch <- afterwhile_redirect
		return err
	}
	fmt.Println("Count: ", count_int)
	fmt.Println("Final URL is", resp.Request.URL)
	if (strings.Contains(resp.Request.URL.String(), "https")) {
		upgrade = true
	}
	fmt.Println("if upgrade: ", upgrade)
	if (count_int>1 && !first_redirect_upgrade && upgrade) {
		afterwhile_redirect = "true"
	}
	fmt.Println("if afterwhile redirect: ", afterwhile_redirect)
	fmt.Println()
	ch <- afterwhile_redirect
	return nil
}

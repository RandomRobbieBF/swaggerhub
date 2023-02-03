package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

type Api struct {
	Properties []struct {
		URL string `json:"url"`
	} `json:"properties"`
}

type Response struct {
	Apis []Api `json:"apis"`
}

func main() {
	search := flag.String("search", "", "The search query to use in the URL")
	flag.Parse()

	if *search == "" {
		flag.Usage()
		return
	}

	url := fmt.Sprintf("https://app.swaggerhub.com/apiproxy/specs?sort=BEST_MATCH&order=DESC&query=%s&limit=25", *search)
	resp, err := http.Get(url)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
		return
	}

	var response Response
	if err := json.Unmarshal(body, &response); err != nil {
		fmt.Println(err)
		return
	}

	urls := make(map[string]bool)
	for _, api := range response.Apis {
		for _, prop := range api.Properties {
			url := strings.Replace(prop.URL, " ", "", -1)
			url = strings.Replace(url, "\n", "", -1)
			if !urls[url] && url != "" {
				fmt.Println(url)
				urls[url] = true
			}
		}
	}
}

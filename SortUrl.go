package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"math"
	"net/http"
	"os"
	"strings"
	"sync"

	"github.com/xrash/smetrics"
)

var urlcontent []string = make([]string, 0)
var channel = make(chan string)
var wg sync.WaitGroup

func main() {

	for i := 1; i <= 30; i++ {
		wg.Add(1)
		go deal_url()
	}
	sc := bufio.NewScanner(os.Stdin)

	for sc.Scan() {
		url := sc.Text()
		if strings.Index(url, "http") == -1 {
			http_url := "http://" + url
			https_url := "https://" + url

			channel <- http_url
			channel <- https_url
		}
		channel <- url

	}
	wg.Wait()

}

func deal_url() {
	defer wg.Done()
	for url := range channel {

		count := 0
		percent := 0.1
		res, err := http.Get(url)
		if err != nil {
			continue
		}
		resp, err := ioutil.ReadAll(res.Body)
		if err != nil {
			continue
		}

		for _, value := range urlcontent {
			if math.Abs(float64(len(resp)-len(value))) > 50 {
				continue
			}

			if len(resp) <= 150 {
				percent = smetrics.JaroWinkler(string(resp), value, 0.7, 4)
			} else {
				percent = smetrics.JaroWinkler(string(resp)[100:150], value, 0.7, 4)
			}

			fmt.Println(percent)
			if percent > 0.95 {

				count = 1
				break
			}
		}
		if count != 1 {

			fmt.Println(url)
			urlcontent = append(urlcontent, string(resp)[0:80])
		}
		res.Body.Close()

	}

}

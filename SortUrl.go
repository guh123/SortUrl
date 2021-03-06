package main

import (
	"bufio"
	"fmt"
	"math"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/valyala/fasthttp"
	"github.com/xrash/smetrics"
)

var duration_Seconds time.Duration = (500 * 10) * time.Millisecond
var urlcontent []string = make([]string, 0)
var channel = make(chan string)
var wg sync.WaitGroup

func main() {

	for i := 1; i <= 30; i++ {
		wg.Add(1)
		go dealURL()
	}
	sc := bufio.NewScanner(os.Stdin)
	for sc.Scan() {
		url := sc.Text()
		if strings.Index(url, "http") == -1 {
			httpUrl := "http://" + url
			httpsUrl := "https://" + url

			channel <- httpUrl
			channel <- httpsUrl
		} else {
			channel <- url
		}

	}

	close(channel)
	wg.Wait()

}

func dealURL() {
	req := fasthttp.AcquireRequest()

	req.Header.SetMethod("GET")

	resp := fasthttp.AcquireResponse()
	client := &fasthttp.Client{}
	client.ReadTimeout = duration_Seconds

	defer func() {
		// 用完需要释放资源
		wg.Done()
		fasthttp.ReleaseResponse(resp)
		fasthttp.ReleaseRequest(req)
	}()

	for url := range channel {
		var text string
		count := 0
		percent := 0.1
		req.SetRequestURI(url)
		err := client.DoTimeout(req, resp, duration_Seconds)
		if err != nil {
			continue
		}

		response := resp.Body()
		for _, value := range urlcontent {
			if math.Abs(float64(len(response)-len(value))) > 50 {
				continue
			}

			if len(response) <= 150 {
				text = string(response)
			} else {
				text = string(response)[100:150]
			}
			percent = smetrics.JaroWinkler(text, value, 0.7, 4)
			if percent > 0.95 {

				count = 1
				break
			}
		}
		if count != 1 {

			fmt.Println(url)
			urlcontent = append(urlcontent, text)
		}

	}

}

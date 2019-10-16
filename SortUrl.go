package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"sync"
)

var urlcontent []string = make([]string, 0)
var channel = make(chan string)
var wg sync.WaitGroup

func SimilarText(first, second string, percent *float64) int {
	var similarText func(string, string, int, int) int
	similarText = func(str1, str2 string, len1, len2 int) int {
		var sum, max int
		pos1, pos2 := 0, 0

		// Find the longest segment of the same section in two strings
		for i := 0; i < len1; i++ {
			for j := 0; j < len2; j++ {
				for l := 0; (i+l < len1) && (j+l < len2) && (str1[i+l] == str2[j+l]); l++ {
					if l+1 > max {
						max = l + 1
						pos1 = i
						pos2 = j
					}
				}
			}
		}

		if sum = max; sum > 0 {
			if pos1 > 0 && pos2 > 0 {
				sum += similarText(str1, str2, pos1, pos2)
			}
			if (pos1+max < len1) && (pos2+max < len2) {
				s1 := []byte(str1)
				s2 := []byte(str2)
				sum += similarText(string(s1[pos1+max:]), string(s2[pos2+max:]), len1-pos1-max, len2-pos2-max)
			}
		}

		return sum
	}

	l1, l2 := len(first), len(second)
	if l1+l2 == 0 {
		return 0
	}
	sim := similarText(first, second, l1, l2)

	*percent = float64(sim*200) / float64(l1+l2)

	return sim
}
func main() {

	for i := 1; i <= 30; i++ {
		wg.Add(1)
		go deal_url()
	}
	sc := bufio.NewScanner(os.Stdin)

	for sc.Scan() {
		url := sc.Text()
		if strings.Index(url, "http") == -1 {
			url = "http://" + url
		}
		channel <- url

	}
	wg.Wait()

}

func deal_url() {
	defer wg.Done()
	for url := range channel {

		count := 0

		res, err := http.Get(url)
		if err != nil {
			continue
		}
		resp, err := ioutil.ReadAll(res.Body)
		if err != nil {
			continue
		}

		for _, value := range urlcontent {
			var percent *float64 = new(float64)
			SimilarText(string(resp), value, percent)

			fmt.Println(*percent)
			if *percent > 30 {

				count = 1
				break
			}
		}
		if count != 1 {

			fmt.Println(url)
			urlcontent = append(urlcontent, string(resp))
		}
		res.Body.Close()

	}

}

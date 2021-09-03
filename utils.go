package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"log"
	"net/url"
	"regexp"
	"sort"
	"strings"
)

func loadURLs() (urls map[string][]string) {
	resp, err := http.Get("https://github.com/DjMike238/pulitiotron/raw/master/urls.json")
	if err != nil {
		log.Println(err)
		return
	}
	defer resp.Body.Close()

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Println(err)
		return
	}

	json.Unmarshal(data, &urls)

	return
}

func getCleanURL(rawURL string, urls map[string][]string) string {
	u, err := url.Parse(rawURL)
	if err != nil {
		log.Println(err)
		return ""
	}

	hostArr := strings.Split(u.Hostname(), ".")
	host := hostArr[1]

	if hostArr[0] != "www" {
		host = hostArr[0]
	}

	if len(urls) > 0 && len(urls[host]) > 0 {
		rx := regexp.MustCompile(urls[host][0])
		return rx.ReplaceAllString(rawURL, urls[host][1])
	}

	return "unsupported"
}

func createSupportedList(urls map[string][]string) string {
	var sb strings.Builder

	keys := make([]string, 0, len(urls))

	for k := range urls {
		keys = append(keys, k)
	}

	sort.Strings(keys)

	for _, k := range keys {
		sb.WriteString(fmt.Sprintf("`%s`\n", k))
	}

	return sb.String()
}

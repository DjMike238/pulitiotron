package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"log"
	"net/url"
	"regexp"
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

	if len(urls) > 0 && len(urls[u.Host]) > 0 {
		rx := regexp.MustCompile(urls[u.Host][0])
		return rx.ReplaceAllString(rawURL, urls[u.Host][1])
	}

	return "unsupported"
}

func createSupportedList(urls map[string][]string) string {
	var sb strings.Builder

	for k, _ := range urls {
		sb.WriteString(fmt.Sprintf("`%s`\n", k))
	}

	return sb.String()
}

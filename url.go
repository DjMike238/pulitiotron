package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"regexp"
	"sort"
	"strings"

	tld "github.com/jpillora/go-tld"
)

// Discard any text before a query.
var beforeRx = regexp.MustCompile("(.* )?(.*)")

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
	url := beforeRx.FindStringSubmatch(rawURL)
	if url == nil {
		return ""
	}

	rawURL = url[2]

	u, err := tld.Parse(rawURL)
	if err != nil {
		log.Println(err)
		return ""
	}

	dom := u.Domain

	if len(urls) > 0 && len(urls[dom]) > 0 {
		rx := regexp.MustCompile(urls[dom][0])
		return rx.ReplaceAllString(rawURL, urls[dom][1])
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
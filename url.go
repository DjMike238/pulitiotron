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

// Custom type for different types of URLs.
type URLType uint

const (
	Empty URLType = iota
	Supported
	Unsupported

	// When trying to clean a Twitter URL, an option to send a
	// vxTwitter (vxtwitter.com) version of the URL will be added.
	Twitter
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

func getCleanURL(rawURL string, urls map[string][]string) (URLType, string) {
	url := beforeRx.FindStringSubmatch(rawURL)
	if url == nil {
		return Empty, ""
	}

	rawURL = url[2]

	u, err := tld.Parse(rawURL)
	if err != nil {
		log.Println(err)
		return Empty, ""
	}

	dom := u.Domain

	if len(urls) > 0 && len(urls[dom]) > 0 {
		rx := regexp.MustCompile(urls[dom][0])
		clean := rx.ReplaceAllString(rawURL, urls[dom][1])

		if dom == "twitter" {
			return Twitter, clean
		}

		return Supported, clean
	}

	return Unsupported, ""
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

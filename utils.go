package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/url"
	"os"
	"path"
	"regexp"
	"strings"
)

func getURLPath() string {
	cfgPath := path.Join(os.Getenv("HOME"), ".config", BOT_NAME)

	_, err := os.Stat(cfgPath)
	if os.IsNotExist(err) {
		os.Mkdir(cfgPath, 0755)
	}

	cfgFile := path.Join(cfgPath, "urls.json")

	file, err := os.OpenFile(cfgFile, os.O_RDWR|os.O_CREATE, 0666)
	if err != nil {
		log.Println(err)
	}
	defer file.Close()

	return cfgFile
}

func loadURLs() (urls map[string][]string) {
	path := getURLPath()

	data, err := ioutil.ReadFile(path)
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

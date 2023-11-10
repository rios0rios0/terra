package main

import (
	logger "github.com/sirupsen/logrus"
	"io/ioutil"
	"net/http"
	"regexp"
)

// Fetch the latest version of software from a URL
func fetchLatestVersion(url, regexPattern string) string {
	resp, err := http.Get(url)
	if err != nil {
		logger.Fatalf("Error fetching version info: %s", err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		logger.Fatalf("Error reading response body: %s", err)
	}

	re := regexp.MustCompile(regexPattern)
	matches := re.FindStringSubmatch(string(body))
	if len(matches) > 1 {
		return matches[1]
	}

	logger.Fatalf("No version match found")
	return ""
}

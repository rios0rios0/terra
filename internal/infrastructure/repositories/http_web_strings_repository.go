package repositories

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"regexp"
	"time"
)

const contextTimeout = 10 * time.Second

type HttpWebStringsRepository struct{}

func (it *HttpWebStringsRepository) FindStringMatchInURL(url, regexPattern string) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), contextTimeout)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return "", fmt.Errorf("error creating request: %w", err)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("error fetching version info: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("error reading response body: %w", err)
	}

	re := regexp.MustCompile(regexPattern)
	matches := re.FindStringSubmatch(string(body))
	if len(matches) > 1 {
		return matches[1], nil
	}

	return "", fmt.Errorf("no version match found, check the regex pattern: %s", regexPattern)
}

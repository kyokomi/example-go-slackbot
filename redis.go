package main

import (
	"fmt"
	"net/url"
)

func parseRedisURL(redisURL string) (string, string, error) {
	addr := fmt.Sprintf("%s:%d", "localhost", 6379)
	password := ""

	if redisURL == "" {
		return addr, password, nil
	}

	redisURLInfo, err := url.Parse(redisURL)
	if err != nil {
		return addr, password, err
	}

	addr = redisURLInfo.Host
	if redisURLInfo.User != nil {
		password, _ = redisURLInfo.User.Password()
	}
	return addr, password, nil
}

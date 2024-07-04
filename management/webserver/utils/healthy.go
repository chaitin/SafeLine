package utils

import (
	"fmt"
	"net/http"
)

func CheckHealthy(url string) bool {
	method := "GET"

	client := &http.Client{}
	req, err := http.NewRequest(method, url, nil)

	if err != nil {
		fmt.Println(err)
		return false
	}
	res, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return false
	}
	defer res.Body.Close()

	return res.StatusCode == 200
}

func CheckWafHealthy() bool {
	// todo: not check the healthy status
	return CheckHealthy("http://safeline-mario:3335") && CheckHealthy("http://safeline-detector:8001/stat")
}

package main

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/labstack/gommon/log"
)

func sendToFlume(message string) error {
	url := fmt.Sprintf("http://192.168.1.7:51401")
	b := []byte("")
	b = strconv.AppendQuote(b, message)
	params := fmt.Sprintf("[{\"body\":%s}]", b)
	// var jsonStr = []byte(params)

	fmt.Printf("social data param: %s\n", params)
	body := strings.NewReader(params)

	req, err := http.NewRequest("POST", url, body)
	if err != nil {
		log.Error(err)
	}
	req.Header.Set("Content-Type", "application/json; charset=UTF-8")
	resp, errC := http.DefaultClient.Do(req)
	if errC != nil {
		log.Error(errC)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK {
		return nil
	}

	return errC
}

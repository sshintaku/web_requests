package web_requests

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"github.com/sshintaku/cloud_types"
)

func PostWebRequest(url string, payload []byte, token *string) (*string, error) {
	request, httpError := http.NewRequest("POST", url, bytes.NewBuffer(payload))
	request.Header.Set("Content-Type", "application/json; charset=UTF-8")
	if token != nil {
		request.Header.Set("x-redlock-auth", *token)
	}

	if httpError != nil {
		return nil, httpError
	}
	var netClient = &http.Client{
		Timeout: time.Second * 50,
	}
	response, responseError := netClient.Do(request)
	if responseError != nil {
		return nil, responseError
	}
	body, httpReadError := ioutil.ReadAll(response.Body)
	if httpReadError != nil {
		return nil, httpReadError
	}
	responseBody := string(body)
	return &responseBody, nil

}

func GetJWTToken(host string, username string, password string) (*cloud_types.PrismaClient, error) {
	var client cloud_types.PrismaClient
	auth := cloud_types.Authentication{
		Username: username,
		Password: password,
	}
	jsonPayload, _ := json.Marshal(auth)

	var token cloud_types.JwtToken
	response, responseError := PostWebRequest(host, jsonPayload, nil)
	if responseError != nil {
		return nil, responseError
	}
	json.Unmarshal([]byte(*response), &token)
	client.Token = token.Token
	return &client, nil
}

func ProcessWebRequest(request *http.Request) ([]byte, error) {
	var netClient = &http.Client{
		Timeout: time.Second * 500,
	}
	response, responseError := netClient.Do(request)
	if responseError != nil {
		return nil, responseError
	}
	body, httpReadError := ioutil.ReadAll(response.Body)
	defer response.Body.Close()
	if httpReadError != nil {
		return nil, httpReadError
	}
	// responseBody := string(body)
	return body, nil
}

func PostMethod(url string, payload []byte, token string) (*string, error) {
	request, httpError := http.NewRequest("POST", url, bytes.NewBuffer(payload))
	request.Header.Add("x-redlock-auth", token)
	request.Header.Set("Content-Type", "application/json; charset=UTF-8")
	if httpError != nil {
		return nil, httpError
	}
	var netClient = &http.Client{
		Timeout: time.Second * 200,
	}
	response, responseError := netClient.Do(request)
	if responseError != nil {
		return nil, responseError
	}
	body, httpReadError := ioutil.ReadAll(response.Body)
	if httpReadError != nil {
		return nil, httpReadError
	}
	responseBody := string(body)
	return &responseBody, nil
}

func GetMethod(uri string, token string) ([]byte, error) {
	method := "GET"
	req, err := http.NewRequest(method, uri, nil)

	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	req.Header.Add("x-redlock-auth", token)
	req.Header.Add("Content-Type", "application/json")

	body, responseError := ProcessWebRequest(req)
	if responseError != nil {
		log.Fatalln(responseError)
		return nil, responseError

	}
	return body, nil

}

func GetComputeBaseUrl(token string) (string, error) {
	url := "https://api2.prismacloud.io/compute/config"
	result, _ := GetMethod(url, token)
	var base cloud_types.Twistlock
	json.Unmarshal(result, &base)
	return base.BaseURI, nil

}

func PutRequest(uri string, payload []byte, token string) string {
	request, httpError := http.NewRequest("PUT", uri, bytes.NewBuffer(payload))
	request.Header.Add("x-redlock-auth", token)
	request.Header.Add("Content-Type", "application/json")
	if httpError != nil {
		log.Fatalln(httpError)
	}
	result, resultError := ProcessWebRequest(request)

	if resultError != nil {
		log.Fatalln(resultError)
	}
	return string(result)
}

package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"
)

const (
	serverAddr string = ":8888"
)

var (
	clientId     string
	clientSecret string
	callbackUri  string = "http://localhost" + serverAddr + "/callback"
	authUrl      string
	tokenUrl     string
)

func init() {
	clientId = os.Getenv("CLIENT_ID")
	if clientId == "" {
		panic("set CLIENT_ID environment variable")
	}
	clientSecret = os.Getenv("CLIENT_SECRET")
	if clientSecret == "" {
		panic("set CLIENT_SECRET environment variable")
	}
	authUrl = os.Getenv("AUTH_URL")
	if authUrl == "" {
		panic("set AUTH_URL environment variable")
	}
	tokenUrl = os.Getenv("TOKEN_URL")
	if tokenUrl == "" {
		panic("set TOKEN_URL environment variable")
	}
}

func main() {
	log.Printf("server is running on %s\n", serverAddr)
	log.Printf("callbackUri is: %s\n", callbackUri)

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("try login via http://localhost" + serverAddr + "/login"))
	})

	http.HandleFunc("/login", func(w http.ResponseWriter, r *http.Request) {
		// Authorization Request
		redirectUrl, err := url.Parse(authUrl)
		if err != nil {
			panic(err)
		}

		q := redirectUrl.Query()
		q.Set("response_type", "code")
		q.Set("client_id", clientId)
		q.Set("redirect_uri", callbackUri)
		q.Set("state", "helloworld")
		redirectUrl.RawQuery = q.Encode()

		log.Println(redirectUrl.String())

		http.Redirect(w, r, redirectUrl.String(), http.StatusTemporaryRedirect)
	})

	http.HandleFunc("/callback", func(w http.ResponseWriter, r *http.Request) {
		// Extract Code from query component
		code := r.URL.Query().Get("code")

		fmt.Printf("code: %v\n", code)

		// Authorization Code Request
		body := exchangeCodeToToken(code)
		fmt.Printf("body: %v\n", body)

		w.Write([]byte(body))
	})

	err := http.ListenAndServe(serverAddr, nil)
	if err != nil {
		panic(err)
	}

}

func exchangeCodeToToken(code string) string {
	u, err := url.Parse(tokenUrl)
	if err != nil {
		panic(err)
	}

	data := url.Values{}
	data.Set("code", code)
	data.Set("redirect_uri", callbackUri)
	data.Set("grant_type", "authorization_code")

	req, err := http.NewRequest("POST", u.String(), strings.NewReader(data.Encode()))
	if err != nil {
		panic(err)
	}

	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	req.SetBasicAuth(clientId, clientSecret)

	client := http.DefaultClient
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}
	return string(body)
}

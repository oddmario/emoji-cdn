package utils

import "github.com/go-resty/resty/v2"

var HttpClient *resty.Client

func InitHttpClient() {
	HttpClient = resty.New()
}

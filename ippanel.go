//
// Copyright (c) 2020 IPPANEL SMS
// All rights reserved.
//
// Author: ippanel <dev@ippanel.com>

// Package ippanel is an official library for working with ippanel sms api.
// brief documentation for ippanel sms api provided at http://docs.ippanel.com
package ippanel

import (
	"net/http"
	"net/url"
	"time"
)

const (
	// ClientVersion is used in User-Agent request header to provide server with API level.
	ClientVersion = "2.0.0"

	// Endpoint points you to Ippanel REST API.
	Endpoint = "https://api2.ippanel.com/api/v1"

	// httpClientTimeout is used to limit http.Client waiting time.
	httpClientTimeout = 30 * time.Second
)

// Ippanel ...
type Ippanel struct {
	Apikey  string
	Client  *http.Client
	BaseURL *url.URL
}

// New create new ippanel sms instance
func New(apikey string) *Ippanel {
	u, _ := url.Parse(Endpoint)
	client := &http.Client{
		Transport: http.DefaultTransport,
		Timeout:   httpClientTimeout,
	}

	return &Ippanel{
		Apikey:  apikey,
		Client:  client,
		BaseURL: u,
	}
}

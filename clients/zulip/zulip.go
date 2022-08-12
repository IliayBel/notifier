package zulip

import (
	"github.com/dghubble/sling"
	"net/http"
	"time"
)

type Client struct {
	client   *http.Client
	sling    *sling.Sling
	Email    string
	APIKey   string
	common   service
	Messages *MessagesService
	Events   *EventsService
}

type service struct {
	client *Client
}

func New(httpClient *http.Client, baseURL string) *Client {
	if httpClient == nil {
		httpClient = &http.Client{
			Timeout: time.Second * 10,
		}
	}

	s := sling.New().Client(httpClient).Base(baseURL)
	c := &Client{sling: s}

	c.common.client = c
	c.Messages = (*MessagesService)(&c.common)
	c.Events = (*EventsService)(&c.common)

	return c
}

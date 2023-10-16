package odata

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"errors"
	"io"
	"net/http"

	"github.com/dainiauskas/go-log"
)

type Client struct {
	baseURL string
	auth    *BaseAuthorization
	header  http.Header
}

func NewClient() *Client {
	return &Client{}
}

func (c *Client) SetBaseURL(url string) {
	c.baseURL = url
}

func (c *Client) GetURL() string {
	return c.baseURL
}

func (c *Client) SetBaseCredentials(cred *BaseAuthorization) {
	c.auth = cred
}

func (c *Client) SetHeaders(h http.Header) {
	c.header = h
}

func (c *Client) Get(method string) ([]byte, http.Header, error) {
	return c.GetFromURL(c.baseURL + method)
}

func (c *Client) GetByURL() ([]byte, http.Header, error) {
	return c.GetFromURL(c.baseURL)
}

func (c *Client) PostAPI(data interface{}) ([]byte, error) {
	client := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		},
	}

	b, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest(http.MethodPost, c.baseURL, bytes.NewBuffer(b))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	req.SetBasicAuth(c.auth.Name, c.auth.Password)

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != 200 {
		return nil, errors.New(resp.Status)
	}
	defer resp.Body.Close()

	return io.ReadAll(resp.Body)
}

func (c *Client) PostFromURL(url string, b []byte) ([]byte, error) {
	client := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		},
	}

	req, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(b))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")

	req.SetBasicAuth(c.auth.Name, c.auth.Password)

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != 201 {
		return nil, errors.New(resp.Status)
	}
	defer resp.Body.Close()

	return io.ReadAll(resp.Body)
}

func (c *Client) DeleteFromURL(url string) error {
	client := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		},
	}

	req, err := http.NewRequest(http.MethodDelete, url, nil)
	if err != nil {
		return err
	}

	req.Header.Set("If-Match", "*")

	req.SetBasicAuth(c.auth.Name, c.auth.Password)

	resp, err := client.Do(req)
	if err != nil {
		return err
	}

	if resp.StatusCode != http.StatusNoContent {
		return errors.New(resp.Status)
	}

	return nil
}

func (c *Client) GetFromURL(url string) ([]byte, http.Header, error) {
	log.Trace("GetFromURL: %s", url)

	client := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		},
	}

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, nil, err
	}

	if c.header != nil {
		req.Header = c.header
	}

	req.SetBasicAuth(c.auth.Name, c.auth.Password)

	resp, err := client.Do(req)
	if err != nil {
		log.Error("client.Do error: %s", err)
		return nil, nil, err
	}

	b, err := io.ReadAll(resp.Body)
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return nil, nil, errors.New(resp.Status)
	}

	return b, resp.Header, err
}

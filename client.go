package odata

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"syscall"
	"time"

	"github.com/dainiauskas/go-log"
)

var (
	MaxIdleConnsPerHost = 1
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
			TLSClientConfig:     &tls.Config{InsecureSkipVerify: true},
			MaxIdleConnsPerHost: MaxIdleConnsPerHost,
			DisableKeepAlives:   true,
			IdleConnTimeout:     time.Millisecond * 100,
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
			TLSClientConfig:     &tls.Config{InsecureSkipVerify: true},
			MaxIdleConnsPerHost: MaxIdleConnsPerHost,
			DisableKeepAlives:   true,
			IdleConnTimeout:     time.Millisecond * 100,
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

	codes := map[int]bool{
		200: true,
		201: true,
	}

	if !codes[resp.StatusCode] {
		return nil, errors.New(resp.Status)
	}
	defer resp.Body.Close()

	return io.ReadAll(resp.Body)
}

func (c *Client) DeleteFromURL(url string) error {
	var resp *http.Response

	client := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig:     &tls.Config{InsecureSkipVerify: true},
			MaxIdleConnsPerHost: MaxIdleConnsPerHost,
			DisableKeepAlives:   true,
			IdleConnTimeout:     time.Millisecond * 100,
		},
	}

	req, err := http.NewRequest(http.MethodDelete, url, nil)
	if err != nil {
		return err
	}

	req.Header.Set("If-Match", "*")

	req.SetBasicAuth(c.auth.Name, c.auth.Password)

	for i := 0; i < 100; i++ {
		var err error

		resp, err = client.Do(req)
		if err != nil {
			if errors.Is(err, syscall.ECONNRESET) {
				continue
			}

			return err
		}

		if resp.StatusCode == http.StatusNotFound {
			return nil
		}

		if resp.StatusCode != http.StatusNoContent {
			return errors.New(resp.Status)
		}

		time.Sleep(time.Millisecond * 200)

		return nil
	}

	return nil
}

func (c *Client) GetFromURL(url string) ([]byte, http.Header, error) {
	log.Trace("GetFromURL: %s", url)

	client := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig:     &tls.Config{InsecureSkipVerify: true},
			MaxIdleConnsPerHost: MaxIdleConnsPerHost,
			DisableKeepAlives:   true,
			IdleConnTimeout:     time.Millisecond * 100,
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

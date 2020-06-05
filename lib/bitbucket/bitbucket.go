package bitbucket

import (
	"errors"
	"fmt"
	"io"
	"net/http"

	"golang.org/x/net/proxy"
)

type PagedResponse struct {
	Size          int  `json:"size"`
	Limit         int  `json:"limit"`
	Start         int  `json:"start"`
	IsLastPage    bool `json:"isLastPage"`
	NextPageStart int  `json:"nextPageStart"`
}

type Links struct {
	Clone []CloneLink `json:"clone,omitempty"`
	Self  *[]struct {
		Href string `json:"href"`
	} `json:"self,omitempty"`
}

type CloneLink struct {
	Href string `json:"href"`
	Name string `json:"name"`
}

type Client struct {
	baseEndpoint string
	Auth         interface{}
	HttpClient   *http.Client
	Projects     Projects
	Repositories Repositories
}

type BasicAuth struct {
	Username string
	Password string
}

type TokenAuth struct {
	Token string
}

func NewClientWithBasicAuth(endpoint, username, password string) *Client {
	c := &Client{
		baseEndpoint: endpoint,
		Auth:         BasicAuth{Username: username, Password: password},
		HttpClient:   new(http.Client),
	}

	c.Projects = Projects{c}
	c.Repositories = Repositories{c}
	return c
}

func NewClientWithTokenAuth(endpoint, token string) *Client {
	c := &Client{
		baseEndpoint: endpoint,
		Auth:         TokenAuth{Token: token},
		HttpClient:   new(http.Client),
	}
	c.Projects = Projects{c}
	c.Repositories = Repositories{c}
	return c
}

func (c *Client) SetSocksProxy(url string) error {
	dialer, err := proxy.SOCKS5("tcp", url, nil, proxy.Direct)
	if err != nil {
		return err
	}
	httpTransport := &http.Transport{}
	httpTransport.Dial = dialer.Dial
	c.HttpClient.Transport = httpTransport

	return nil
}

func (c *Client) RawRequest(method, url, text string) (io.ReadCloser, error) {
	// body := strings.NewReader(text)

	req, err := http.NewRequest(method, url, nil)
	if err != nil {
		return nil, err
	}
	if text != "" {
		req.Header.Set("Content-Type", "application/json")
	}

	err = c.authenticateRequest(req)
	if err != nil {
		return nil, err
	}

	return c.doRawRequest(req, false)
}

func (c *Client) authenticateRequest(req *http.Request) error {
	switch auth := c.Auth.(type) {
	case BasicAuth:
		req.SetBasicAuth(auth.Username, auth.Password)
	case TokenAuth:
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", auth.Token))
	default:
		return errors.New("Unsupported authentication method")
	}

	return nil
}

func (c *Client) doRawRequest(req *http.Request, emtpyResponse bool) (io.ReadCloser, error) {
	resp, err := c.HttpClient.Do(req)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		resp.Body.Close()
		return nil, errors.New(resp.Status)
	}

	if emtpyResponse {
		resp.Body.Close()
		return nil, nil
	}

	if resp.Body == nil {
		resp.Body.Close()
		return nil, errors.New("response body is nil")
	}

	return resp.Body, nil
}

package client

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"time"

	log "github.com/sirupsen/logrus"
)

type Client struct {
	HTTPClient *http.Client
	BaseURL    string
	Domain     string
	Token      string
	Verbose    bool
}

func NewClient(domain, token string, verbose bool) *Client {
	return &Client{
		HTTPClient: &http.Client{
			Timeout: time.Second * 30,
		},
		BaseURL: "https://www.duckdns.org/update?domains=%s&token=%s",
		Domain:  domain,
		Token:   token,
		Verbose: verbose,
	}
}

func (c *Client) SetHTTPClient(httpClient *http.Client) {
	c.HTTPClient = httpClient
}

func (c *Client) SetBaseURL(baseURL string) {
	c.BaseURL = baseURL
}

func (c *Client) SetDomain(domain string) {
	c.Domain = domain
}

func (c *Client) SetToken(token string) {
	c.Token = token
}

func (c *Client) SetVerbose(verbose bool) {
	c.Verbose = verbose
}

// Update updates the IP in DuckDNS, allowing the passing of an IP
func (u *Client) Update(ip, ipv6 string) error {
	log.Debugf("updating IP with %+v", u)

	url := fmt.Sprintf(u.BaseURL, u.Domain, u.Token)

	if ip != "" {
		url = url + "&ip=" + ip
	}

	if ipv6 != "" {
		url = url + "&ipv6=" + ipv6
	}

	if u.Verbose {
		url = url + "&verbose=true"
	}

	log.Debugf("using URL for update %s", url)

	resp, err := u.HTTPClient.Get(url)
	if err != nil {
		return err
	}

	log.Debugf("got http response %+v", resp)

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	log.Debugf("response body from request: %s", string(body))

	lines := strings.Split(string(body), "\n")
	if len(lines) > 0 && lines[0] == "OK" {
		return nil
	} else if len(lines) > 0 && lines[0] == "KO" {
		return fmt.Errorf("got failure response to update: %s", string(body))
	}

	return fmt.Errorf("unexpected response from update: %s", string(body))
}

package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
)

type Client struct {
	Host string
}

func New(h string) Client {
	return Client{Host: h}
}

func Default() Client {
	return Client{Host: "http://127.0.0.1:8000"}
}

func (c *Client) Setup(h string) {
	c.Host = h
}

func (c *Client) get(u string) (*http.Response, error) {
	uri := fmt.Sprintf("%s%s", c.Host, u)
	req, _ := http.NewRequest("GET", uri, nil)
	return c.makeRequest(req)
}

func (c *Client) makeRequest(req *http.Request) (*http.Response, error) {
	//log.Printf("%s request to: %s\n", req.Method, req.URL)
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	if resp.StatusCode >= 400 {
		err := errorHandler(resp)
		if err != nil {
			return resp, err
		}
	}
	return resp, nil
}

func (c *Client) Status(o string) (StatusResponse, error) {
	route := fmt.Sprintf("/api/workflows/v1/%s/status", o)
	var sr StatusResponse
	r, err := c.get(route)
	if err != nil {
		return sr, err
	}
	defer r.Body.Close()

	if err := json.NewDecoder(r.Body).Decode(&sr); err != nil {
		return sr, err
	}
	return sr, nil
}

func (c *Client) Metadata(o string, p url.Values) (MetadataResponse, error) {
	route := fmt.Sprintf("/api/workflows/v1/%s/metadata"+"?"+p.Encode(), o)
	var mr MetadataResponse
	r, err := c.get(route)
	if err != nil {
		return mr, err
	}

	defer r.Body.Close()
	if err := json.NewDecoder(r.Body).Decode(&mr); err != nil {
		return mr, err
	}
	return mr, err
}

func errorHandler(r *http.Response) error {
	var er = ErrorResponse{
		HTTPStatus: r.Status,
	}
	if err := json.NewDecoder(r.Body).Decode(&er); err != nil {
		log.Println("No json body in response")
	}
	return fmt.Errorf("submission failed. the server returned %#v", er)
}

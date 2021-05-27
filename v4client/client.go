package v4client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
)

type V4Client interface {
	ListBranch(owner, name string) ([]string, error)
}

const (
	defaultBaseURL = "https://api.github.com/graphql"
)

type v4Client struct {
	Client    *http.Client
	BaseURL     *url.URL
	AccessToken string

	Issue *IssueService
}

func NewClient(httpClient *http.Client, accessToken string) (*v4Client, error) {

	if httpClient == nil {
		httpClient = http.DefaultClient
	}

	baseURL, err := url.Parse(defaultBaseURL)

	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	c := &v4Client{
		BaseURL:     baseURL,
		AccessToken: accessToken,
		Client:      httpClient,
	}

	c.Issue = &IssueService{client: c}

	return c, nil
}

func (c *v4Client) Do(query string, variables map[string]interface{}, data interface{}) error {
	reqBody, err := json.Marshal(map[string]interface{}{"query": query, "variables": variables})
	if err != nil {
		return err
	}

	req, err := http.NewRequest(http.MethodPost, c.BaseURL.String(), bytes.NewBuffer(reqBody))
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/json; charset=utf-8")
	req.Header.Set("Authorization", fmt.Sprintf("token %s", c.AccessToken))
	req.Header.Set("Accept", "application/vnd.github.antiope-preview+json")
	req.Header.Set("Accept", "application/json")
	req.Header.Set("User-Agent", "hayashiki")
	log.Print(c.BaseURL.String())
	log.Print(query)
	log.Print(variables)

	resp, err := c.Client.Do(req)
	if err != nil {
		log.Print(err)
		return err
	}
	defer resp.Body.Close()

	return handleResponse(resp, data)
}

type graphQLResponse struct {
	Data   interface{}
	Errors []GraphQLError
}

type GraphQLError struct {
	Type    string
	Path    []string
	Message string
}

type GraphQLErrorResponse struct {
	Errors []GraphQLError
}

func handleResponse(resp *http.Response, data interface{}) error {
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	gr := &graphQLResponse{Data: data}
	err = json.Unmarshal(body, &gr)
	// For Debug
	// log.Print(string(body))
	if err != nil {
		return err
	}
	return nil
}


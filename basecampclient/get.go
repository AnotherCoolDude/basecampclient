package basecampclient

import (
	"fmt"
	"net/http"
)

// Get makes a get request to basecamp with the specified url
func (c *Client) Get(url string) (*http.Response, error) {
	req, err := http.NewRequest("GET", url, http.NoBody)
	if err != nil {
		fmt.Println("error making request " + err.Error())
		return nil, err
	}
	c.addHeader(req)
	resp, err := c.httpclient.Do(req)
	if err != nil {
		fmt.Println("error making request " + err.Error())
		return nil, err
	}
	return resp, nil
}

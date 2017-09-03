package heimdall

import (
	"io/ioutil"
	"net/http"
	"time"

	"github.com/pkg/errors"
)

type HTTPClient interface {
	Get(url string) (HeimdallResponse, error)
}

type Client struct {
	client *http.Client
}

func NewClient(config Config) *Client {
	timeout := config.timeoutInSeconds

	httpTimeout := time.Duration(timeout) * time.Second
	return &Client{
		client: &http.Client{
			Timeout: httpTimeout,
		},
	}
}

func (c *Client) Get(url string) (HeimdallResponse, error) {
	response := HeimdallResponse{}

	request, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return response, errors.Wrap(err, "GET - request creation failed")
	}

	return c.do(request)
}

func (c *Client) do(request *http.Request) (HeimdallResponse, error) {
	hr := HeimdallResponse{}
	var err error

	request.Close = true

	response, err := c.client.Do(request)
	if err != nil {
		return hr, err
	}

	defer response.Body.Close()

	var responseBytes []byte
	if response.Body != nil {
		responseBytes, err = ioutil.ReadAll(response.Body)
		if err != nil {
			return hr, err
		}
	}

	hr.body = responseBytes
	hr.statusCode = response.StatusCode

	return hr, err
}

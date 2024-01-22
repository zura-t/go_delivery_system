package httpclient

import "net/http"

//go:generate mockery --name=HttpClientI --with-expecter=true
type HttpClientI interface {
	SendRequest(request any, method string, url string) (*http.Response, error)
}

type HttpClient struct {
	Do func(req *http.Request) (*http.Response, error)
	c  *http.Client
}

var (
	Client HttpClient
)

func New() HttpClient {
	return HttpClient{
		c: &http.Client{},
	}
}

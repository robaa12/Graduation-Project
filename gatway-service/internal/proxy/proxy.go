package proxy

import (
	"fmt"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"

	"github.com/robaa12/gatway-service/config"
)

type ProxyService struct {
	config *config.Config
}

func NewProxyService(cfg *config.Config) *ProxyService {
	return &ProxyService{config: cfg}
}

func (p *ProxyService) UserServiceProxy() http.Handler {

	return p.createProxy(p.config.Services.UserService.URL)
}

func (p *ProxyService) ProductServiceProxy() http.Handler {
	return p.createProxy(p.config.Services.ProductService.URL)
}

func (p *ProxyService) OrderServiceProxy() http.Handler {
	return p.createProxy(p.config.Services.OrderService.URL)
}

func (p *ProxyService) createProxy(targetURL string) http.Handler {
	target, err := url.Parse(targetURL)
	fmt.Println(targetURL)
	if err != nil {
		panic(err)
	}

	return &httputil.ReverseProxy{
		Director: func(r *http.Request) {
			r.Header.Add("X-Forwarded-Host", r.Host)
			r.URL.Scheme = target.Scheme
			r.URL.Host = target.Host
			r.URL.RawPath = r.URL.Path
		},
		ModifyResponse: func(r *http.Response) error {
			return nil
		},
	}
}

// Helper function to remove the first path segment
func stripFirstPathSegment(path string) string {
	parts := strings.SplitN(path, "/", 3) // Split into at most 3 parts: ["", "products", "stores/1/products"]
	if len(parts) < 3 {
		return "/" // If there's nothing left, return root
	}
	return "/" + parts[2] // Return the remaining path
}

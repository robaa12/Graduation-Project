package proxy

import (
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"

	"github.com/robaa12/gatway-service/internal/config"
)

type ProxyService struct {
	serviceConfig *config.ServiceConfig
}

func NewProxyService(serviceConfig *config.ServiceConfig) http.Handler {
	return &ProxyService{serviceConfig: serviceConfig}
}

func (proxyService *ProxyService) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	proxy := proxyService.createProxy()

	proxy.ServeHTTP(w, r)
}

func (p *ProxyService) createProxy() http.Handler {
	target, err := url.Parse(p.serviceConfig.URL)
	log.Println(p.serviceConfig.URL)
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

package proxy

import (
	"github.com/robaa12/gatway-service/internal/config"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
)

type Service struct {
	serviceConfig *config.ServiceConfig
}

func NewProxyService(serviceConfig *config.ServiceConfig) http.Handler {
	return &Service{serviceConfig: serviceConfig}
}

func (proxyService *Service) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	proxy := proxyService.createProxy()

	proxy.ServeHTTP(w, r)
}

func (proxyService *Service) createProxy() http.Handler {
	target, err := url.Parse(proxyService.serviceConfig.URL)
	log.Println(proxyService.serviceConfig.URL)
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

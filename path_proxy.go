package path_proxy

import (
	"net/http"
	"net/http/httputil"
	"net/url"
	"regexp"
	"sort"
)

type Route struct {
	Pattern *regexp.Regexp
	Handler http.Handler
}

type ProxyDef struct {
	Pattern string
	Target  string
}

type ProxyDefSet []ProxyDef

type PathProxyHandler struct {
	Routes []*Route
	sorted bool
}

func (h *PathProxyHandler) sort() {
	if h.sorted {
		return
	}

	sort.Slice(h.Routes, func(i, j int) bool {
		return len(h.Routes[i].Pattern.String()) > len(h.Routes[j].Pattern.String())
	})

	h.sorted = true
}

func (h *PathProxyHandler) ProxyDefintions(proxies ProxyDefSet) error {
	for _, def := range proxies {
		if err := h.Proxy(def.Pattern, def.Target); err != nil {
			return nil
		}
	}

	return nil
}

func (h *PathProxyHandler) Proxy(pattern, targetUrl string) error {
	patternRe, err := regexp.Compile(pattern)
	if err != nil {
		return nil
	}

	urlParsed, err := url.Parse(targetUrl)
	if err != nil {
		return err
	}

	proxy := httputil.NewSingleHostReverseProxy(urlParsed)
	if err != nil {
		return nil
	}

	h.Handler(patternRe, proxy)
	return nil
}

func (h *PathProxyHandler) Handler(pattern *regexp.Regexp, handler http.Handler) {
	h.Routes = append(h.Routes, &Route{
		Pattern: pattern,
		Handler: handler,
	})
	h.sorted = false
}

func (h *PathProxyHandler) HandleFunc(pattern *regexp.Regexp, handler func(http.ResponseWriter, *http.Request)) {
	h.Handler(pattern, http.HandlerFunc(handler))
}

func (h *PathProxyHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h.sort()

	for _, Route := range h.Routes {
		if Route.Pattern.MatchString(r.URL.Path) {
			Route.Handler.ServeHTTP(w, r)
			return
		}
	}

	// no pattern matched; send 404 response
	http.NotFound(w, r)
}

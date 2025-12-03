package registry

import (
	"net/url"
	"sort"
	"strings"

	"Getway_market/internal/application"
	"Getway_market/internal/domain"
)

type MemoryRegistry struct {
	services map[string]*url.URL
	routes   []domain.Route
}

func NewMemoryRegistry(services map[string]string, routes []domain.Route) *MemoryRegistry {
	svcs := make(map[string]*url.URL, len(services))
	for name, raw := range services {
		u, err := url.Parse(raw)
		if err == nil {
			svcs[name] = u
		}
	}
	// сортируем по длине префикса (длинные раньше)
	rs := append([]domain.Route(nil), routes...)
	sort.Slice(rs, func(i, j int) bool {
		return len(rs[i].Prefix) > len(rs[j].Prefix)
	})
	return &MemoryRegistry{services: svcs, routes: rs}
}

func (m *MemoryRegistry) MatchByPrefix(path string) (domain.Route, bool) {
	for _, r := range m.routes {
		if strings.HasPrefix(path, r.Prefix) {
			return r, true
		}
	}
	return domain.Route{}, false
}

func (m *MemoryRegistry) FindRule(method, path string, rules []domain.Rule) (domain.Rule, map[string]string, bool) {
	method = strings.ToUpper(method)
	for _, r := range rules {
		if strings.ToUpper(r.Method) != method {
			continue
		}
		mr := application.MatchPattern(r.In, path)
		if mr.OK {
			return r, mr.Params, true
		}
	}
	return domain.Rule{}, nil, false
}

func (m *MemoryRegistry) ServiceBaseURL(name string) (*url.URL, bool) {
	u, ok := m.services[name]
	return u, ok
}

func (m *MemoryRegistry) TimeoutFor(route domain.Route) int {
	if route.TimeoutMs <= 0 {
		return 2000
	}
	return route.TimeoutMs
}

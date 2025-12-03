package application

import (
	"context"
	"errors"
	"net/url"
	"strings"

	"Getway_market/internal/domain"
)

var (
	ErrNoRoute      = errors.New("no route")
	ErrNoRule       = errors.New("no rule")
	ErrUnauthorized = errors.New("unauthorized")
)

type Registry interface {
	MatchByPrefix(path string) (domain.Route, bool)
	FindRule(method, path string, rules []domain.Rule) (domain.Rule, map[string]string, bool)
	ServiceBaseURL(name string) (*url.URL, bool)
	TimeoutFor(route domain.Route) int
}

type Auth interface {
	Validate(ctx context.Context, token string) error
}

type ResolveResult struct {
	TargetBase *url.URL
	Path       string
	Query      url.Values
	TimeoutMs  int
}

type RouterService struct {
	reg  Registry
	auth Auth
}

func NewRouterService(reg Registry, auth Auth) *RouterService {
	return &RouterService{reg: reg, auth: auth}
}

func (s *RouterService) Resolve(ctx context.Context, method, rawPath string, rawQuery url.Values, bearerToken string) (*ResolveResult, error) {
	route, ok := s.reg.MatchByPrefix(rawPath)
	if !ok {
		return nil, ErrNoRoute
	}
	if route.AuthRequired {
		if err := s.auth.Validate(ctx, bearerToken); err != nil {
			return nil, ErrUnauthorized
		}
	}
	rule, params, ok := s.reg.FindRule(method, rawPath, route.Rules)
	if !ok {
		return nil, ErrNoRule
	}
	base, ok := s.reg.ServiceBaseURL(route.ServiceName)
	if !ok {
		return nil, ErrNoRoute
	}
	// подставляем out + мержим query
	outPath, outQuery := ApplyOutPath(rule.Out, params, rawQuery)
	// аккуратно склеиваем base.Path и outPath
	finalPath := singleJoinPath(base.Path, outPath)
	return &ResolveResult{
		TargetBase: base,
		Path:       finalPath,
		Query:      outQuery,
		TimeoutMs:  s.reg.TimeoutFor(route),
	}, nil
}

func singleJoinPath(a, b string) string {
	switch {
	case a == "" || a == "/":
		if strings.HasPrefix(b, "/") {
			return b
		}
		return "/" + b
	case strings.HasSuffix(a, "/") && strings.HasPrefix(b, "/"):
		return a + b[1:]
	case !strings.HasSuffix(a, "/") && !strings.HasPrefix(b, "/"):
		return a + "/" + b
	default:
		return a + b
	}
}

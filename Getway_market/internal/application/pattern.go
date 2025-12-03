package application

import (
	"net/url"
	"strings"
)

type MatchResult struct {
	Params map[string]string
	OK     bool
}

// MatchPattern сопоставляет путь с шаблоном вида /a/:id/b
func MatchPattern(pattern, path string) MatchResult {
	pp := splitNoEmpty(pattern, '/')
	pv := splitNoEmpty(path, '/')
	if len(pp) != len(pv) {
		return MatchResult{OK: false}
	}
	params := make(map[string]string, 2)
	for i := range pp {
		if len(pp[i]) > 0 && pp[i][0] == ':' {
			params[pp[i][1:]] = pv[i]
			continue
		}
		if pp[i] != pv[i] {
			return MatchResult{OK: false}
		}
	}
	return MatchResult{Params: params, OK: true}
}

// ApplyOutPath подставляет :vars и собирает финальный path и query
func ApplyOutPath(out string, params map[string]string, incoming url.Values) (path string, query url.Values) {
	pathQuery := strings.SplitN(out, "?", 2)

	// path с подстановками
	path = substitute(pathQuery[0], params)

	// исходный query → основа
	q := url.Values{}
	for k, v := range incoming {
		for _, vv := range v {
			q.Add(k, vv)
		}
	}
	// дополняем из out (?id=:id), перетирая входные
	if len(pathQuery) == 2 {
		raw := pathQuery[1]
		parsed, _ := url.ParseQuery(substitute(raw, params))
		for k, v := range parsed {
			q.Del(k)
			for _, vv := range v {
				q.Add(k, vv)
			}
		}
	}
	return path, q
}

func substitute(s string, params map[string]string) string {
	for k, v := range params {
		s = strings.ReplaceAll(s, ":"+k, v)
	}
	return s
}

func splitNoEmpty(s string, sep byte) []string {
	if s == "" || s == "/" {
		return []string{}
	}
	parts := strings.Split(s, string(sep))
	res := make([]string, 0, len(parts))
	for _, p := range parts {
		if p != "" {
			res = append(res, p)
		}
	}
	return res
}

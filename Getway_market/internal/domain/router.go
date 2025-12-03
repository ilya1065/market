package domain

type Method string

type Rule struct {
	Method string `yaml:"method"`
	In     string `yaml:"in"`
	Out    string `yaml:"out"`
}

type Route struct {
	Prefix       string `yaml:"prefix"`
	ServiceName  string `yaml:"service"`
	AuthRequired bool   `yaml:"auth_required"`
	TimeoutMs    int    `yaml:"timeout_ms"`
	Rules        []Rule `yaml:"rules"`
}

type Services map[string]Service

type Service struct {
	BaseURL string `yaml:"base_url"`
}

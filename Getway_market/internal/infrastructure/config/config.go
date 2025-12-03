package config

import (
	"os"
	"strings"
	"time"

	"Getway_market/internal/domain"

	"github.com/joho/godotenv"
	"gopkg.in/yaml.v3"
)

type rawConfig struct {
	Services map[string]domain.Service `yaml:"services"`
	Routes   []domain.Route            `yaml:"routes"`
}

type Config struct {
	Addr        string
	JWTSecret   string
	CORSOrigins []string
	Routes      rawConfig
}

func Load() (*Config, error) {
	_ = godotenv.Load()
	addr := getenv("ADDR", ":8080")
	jwt := os.Getenv("JWT_SECRET")
	cors := getenv("CORS_ALLOWED_ORIGINS", "*")
	routesFile := getenv("ROUTES_FILE", "./configs/routes.yaml")

	b, err := os.ReadFile(routesFile)
	if err != nil {
		return nil, err
	}

	// Подставим env-переменные вида ${NAME}
	b = expandEnv(b)
	var rc rawConfig
	if err := yaml.Unmarshal(b, &rc); err != nil {
		return nil, err
	}

	return &Config{
		Addr:        addr,
		JWTSecret:   jwt,
		CORSOrigins: splitCSV(cors),
		Routes:      rc,
	}, nil
}

func splitCSV(s string) []string {
	parts := strings.Split(s, ",")
	out := make([]string, 0, len(parts))
	for _, p := range parts {
		p = strings.TrimSpace(p)
		if p != "" {
			out = append(out, p)
		}
	}
	if len(out) == 0 {
		return []string{"*"}
	}
	return out
}

func getenv(k, def string) string {
	if v := os.Getenv(k); v != "" {
		return v
	}
	return def
}

// very small ${ENV} expansion
func expandEnv(b []byte) []byte {
	s := string(b)
	lookup := func(key string) string {
		if v := os.Getenv(key); v != "" {
			return v
		}
		return "${" + key + "}"
	}
	var out strings.Builder
	for i := 0; i < len(s); {
		if i+2 < len(s) && s[i] == '$' && s[i+1] == '{' {
			j := i + 2
			for j < len(s) && s[j] != '}' {
				j++
			}
			if j < len(s) && s[j] == '}' {
				key := s[i+2 : j]
				out.WriteString(lookup(key))
				i = j + 1
				continue
			}
		}
		out.WriteByte(s[i])
		i++
	}
	return []byte(out.String())
}

// helpers (экспорт если понадобится)
var _ = time.Now

package config

import (
	"fmt"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	HTTP  HTTPConfig
	DB    DBConfig
	Auth  AuthConfig
	Yandex YandexConfig
}

type HTTPConfig struct {
	ListenAddr string
	// Origin for CORS when using cookie-based auth.
	CorsOrigin string
}

type DBConfig struct {
	DSN           string
	MigrationsDir string
}

type AuthConfig struct {
	JWTSecret string
	// TTL for access token. Refresh token rotation can be added later.
	AccessTokenTTLSeconds int
	CookieSecure          bool

	AdminEmail    string
	AdminPassword string
}

type YandexConfig struct {
	GeocoderKey string
	// Base URL for Yandex HTTP Geocoder API.
	GeocoderBaseURL string
}

func Load() (*Config, error) {
	// Load local .env if present.
	_ = godotenv.Load()

	jwtSecret := getenvStr("TRUMPLIN_JWT_SECRET", "dev_jwt_secret_change_me")
	geocoderKey := os.Getenv("YANDEX_GEOCODER_KEY")
	if geocoderKey == "" {
		// Not fatal for compile/dev: public map can be stubbed until geocoder is implemented.
		geocoderKey = "CHANGE_ME"
	}

	httpPort := getenvInt("TRUMPLIN_HTTP_PORT", 8080)

	cfg := &Config{
		HTTP: HTTPConfig{
			ListenAddr: fmt.Sprintf(":%d", httpPort),
			CorsOrigin: getenvStr("TRUMPLIN_CORS_ORIGIN", "http://localhost:3000"),
		},
		DB: DBConfig{
			DSN:            mustEnv("TRUMPLIN_DATABASE_DSN"),
			MigrationsDir:  "migrations",
		},
		Auth: AuthConfig{
			JWTSecret: jwtSecret,
			AccessTokenTTLSeconds: getenvInt("TRUMPLIN_ACCESS_TOKEN_TTL_SECONDS", 900),
			CookieSecure:          getenvInt("TRUMPLIN_COOKIE_SECURE", 0) == 1,
			AdminEmail:           os.Getenv("TRUMPLIN_ADMIN_EMAIL"),
			AdminPassword:        os.Getenv("TRUMPLIN_ADMIN_PASSWORD"),
		},
		Yandex: YandexConfig{
			GeocoderKey:    geocoderKey,
			GeocoderBaseURL: getenvStr("YANDEX_GEOCODER_BASE_URL", "https://geocode-maps.yandex.ru/1.x/"),
		},
	}

	return cfg, nil
}

func getenvInt(key string, def int) int {
	if v := os.Getenv(key); v != "" {
		if i, err := strconv.Atoi(v); err == nil {
			return i
		}
	}
	return def
}

func getenvStr(key, def string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return def
}

func mustEnv(key string) string {
	v := os.Getenv(key)
	if v == "" {
		panic(fmt.Errorf("missing env var: %s", key))
	}
	return v
}

func (c *HTTPConfig) Addr() string { return c.ListenAddr }


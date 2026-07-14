package config

type Config struct {
	Environment string         `mapstructure:"environment" validate:"required,oneof=development production staging"`
	Debug       bool           `mapstructure:"debug"`
	Version     string         `mapstructure:"version" validate:"required"`
	Server      ServerConfig   `mapstructure:"server" validate:"required"`
	Database    DatabaseConfig `mapstructure:"database" validate:"required"`
	Redis       RedisConfig    `mapstructure:"redis" validate:"required"`
	Cors        CorsConfig     `mapstructure:"cors" validate:"required"`
	Auth        AuthConfig     `mapstructure:"auth" validate:"required"`
	Session     SessionConfig  `mapstructure:"session" validate:"required"`
	Frontend    FrontendConfig `mapstructure:"frontend" validate:"required"`
	Compute     ComputeConfig  `mapstructure:"compute" validate:"required"`
	Stripe      StripeConfig   `mapstructure:"stripe" validate:"required"`
	Mail        MailConfig     `mapstructure:"mail" `
}

type ServerConfig struct {
	Host string    `mapstructure:"host" validate:"required"`
	Port int       `mapstructure:"port" validate:"required,min=1,max=65535"`
	TLS  TLSConfig `mapstructure:"tls" validate:"required"`
}

type FrontendConfig struct {
	URL string `mapstructure:"url" validate:"required,url"`
}

type DatabaseConfig struct {
	Host     string `mapstructure:"host" validate:"required"`
	Port     int    `mapstructure:"port" validate:"required,min=1,max=65535"`
	User     string `mapstructure:"user" validate:"required"`
	Password string `mapstructure:"password" validate:"required"`
	Name     string `mapstructure:"name" validate:"required"`
	SSLMode  string `mapstructure:"sslmode" validate:"required,oneof=disable require verify-ca verify-full"`
	Migrate  bool   `mapstructure:"migrate"`
}

type RedisConfig struct {
	Host     string `mapstructure:"host" validate:"required"`
	Port     int    `mapstructure:"port" validate:"required,min=1,max=65535"`
	Username string `mapstructure:"username"`
	Password string `mapstructure:"password"`
	DB       int    `mapstructure:"db" validate:"min=0"`
}

type CorsConfig struct {
	Enabled          bool     `mapstructure:"enabled"`
	AllowedOrigins   []string `mapstructure:"allowed_origins" validate:"required_if=Enabled true"`
	AllowedMethods   []string `mapstructure:"allowed_methods" validate:"required_if=Enabled true"`
	AllowedHeaders   []string `mapstructure:"allowed_headers" validate:"required_if=Enabled true"`
	ExposedHeaders   []string `mapstructure:"exposed_headers"`
	AllowCredentials bool     `mapstructure:"allow_credentials"`
	MaxAge           int      `mapstructure:"max_age" validate:"min=0"`
}

type AuthConfig struct {
	Google        GoogleConfig `mapstructure:"google" validate:"required"`
	SessionSecret string       `mapstructure:"session_secret" validate:"required"`
}

type GoogleConfig struct {
	ClientID     string `mapstructure:"client_id" validate:"required"`
	ClientSecret string `mapstructure:"client_secret" validate:"required"`
	CallbackURL  string `mapstructure:"callback_url" validate:"omitempty,url"`
}

type SessionConfig struct {
	Name     string       `mapstructure:"name" validate:"required"`
	Cookie   CookieConfig `mapstructure:"cookie"`
	Duration int          `mapstructure:"duration" validate:"required,min=1"` // in seconds
}

type CookieConfig struct {
	Path     string `mapstructure:"path"`
	Domain   string `mapstructure:"domain"`
	Secure   bool   `mapstructure:"secure"`
	HTTPOnly bool   `mapstructure:"http_only"`
	SameSite string `mapstructure:"same_site" validate:"oneof=strict lax none"`
}

type TLSConfig struct {
	Enabled bool `mapstructure:"enabled"`
}

type MailConfig struct {
	From     string `mapstructure:"from" `
	Password string `mapstructure:"password" `
	Name     string `mapstructure:"name" `
}

type ComputeConfig struct {
	Price PriceConfig `mapstructure:"price" validate:"required"`
}

type PriceConfig struct {
	CPUHourlyCents int64 `mapstructure:"cpu_hourly_cents" validate:"required,min=1"`
	GPUHourlyCents int64 `mapstructure:"gpu_hourly_cents" validate:"required,min=1"`
	NPUHourlyCents int64 `mapstructure:"npu_hourly_cents" validate:"required,min=1"`
}

type StripeConfig struct {
	SecretKey     string `mapstructure:"secret_key" validate:"required"`
	WebhookSecret string `mapstructure:"webhook_secret" validate:"required"`
	CPUPriceID    string `mapstructure:"cpu_price_id" validate:"required"`
	GPUPriceID    string `mapstructure:"gpu_price_id" validate:"required"`
	NPUPriceID    string `mapstructure:"npu_price_id" validate:"required"`
}

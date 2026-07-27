package config

import (
	"fmt"
	"os"
	"strings"
	"sync"

	"github.com/go-playground/validator/v10"
	"github.com/joho/godotenv"
	"github.com/spf13/viper"
)

var (
	cfg      *Config
	once     sync.Once
	validate *validator.Validate
)

func init() {
	loadEnvFile()
	validate = validator.New()
	if err := validate.RegisterValidation("environment", validateEnvironment); err != nil {
		panic(fmt.Sprintf("Failed to register custom validator: %v", err))
	}
}

func validateEnvironment(fl validator.FieldLevel) bool {
	env := fl.Field().String()
	validEnvs := []string{"development", "staging", "production"}
	for _, validEnv := range validEnvs {
		if env == validEnv {
			return true
		}
	}
	return false
}

func LoadConfig() (*Config, error) {
	var err error
	once.Do(func() {
		err = loadConfig()
	})
	return cfg, err
}

func Get() *Config {
	if cfg == nil {
		var err error
		cfg, err = LoadConfig()
		if err != nil {
			panic(fmt.Sprintf("Failed to load configuration: %v", err))
		}
	}
	return cfg
}

func loadEnvFile() {
	envFiles := []string{".env.local", ".env"}
	for _, file := range envFiles {
		if _, err := os.Stat(file); err == nil {
			if err := godotenv.Load(file); err != nil {
				fmt.Printf("Warning: Could not load %s: %v\n", file, err)
			}
			return
		}
	}
}

func loadConfig() error {
	if err := setupViper(); err != nil {
		return err
	}

	cfg = &Config{}
	if err := viper.Unmarshal(cfg); err != nil {
		return fmt.Errorf("error unmarshaling config: %w", err)
	}

	if err := validate.Struct(cfg); err != nil {
		return fmt.Errorf("configuration validation failed: %w", err)
	}

	return nil
}

func setupViper() error {
	configNames := []string{"config", "dev", "development", "prod", "production"}
	configPaths := []string{".", "./configs", "./config", "./configurations"}

	viper.SetConfigType("toml")

	for _, path := range configPaths {
		viper.AddConfigPath(path)
	}

	var configFound bool
	for _, name := range configNames {
		viper.SetConfigName(name)
		if err := viper.ReadInConfig(); err == nil {
			configFound = true
			break
		}
	}

	if !configFound {
		return fmt.Errorf("no config file found. Expected one of: %v in paths: %v", configNames, configPaths)
	}

	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	viper.AutomaticEnv()

	creditPriceBindings := map[string]string{
		"stripe.secret_key":     "STRIPE_SECRET_KEY",
		"stripe.webhook_secret": "STRIPE_WEBHOOK_SECRET",
		"stripe.pro_price_id":   "STRIPE_PRO_PRICE_ID",

		"stripe.credit_prices.credits_240.price_id":  "STRIPE_CREDITS_240_PRICE_ID",
		"stripe.credit_prices.credits_240.amount":    "STRIPE_CREDITS_240_AMOUNT",
		"stripe.credit_prices.credits_240.currency":  "STRIPE_CREDITS_240_CURRENCY",
		"stripe.credit_prices.credits_2400.price_id": "STRIPE_CREDITS_2400_PRICE_ID",
		"stripe.credit_prices.credits_2400.amount":   "STRIPE_CREDITS_2400_AMOUNT",
		"stripe.credit_prices.credits_2400.currency": "STRIPE_CREDITS_2400_CURRENCY",
		"stripe.credit_prices.credits_5000.price_id": "STRIPE_CREDITS_5000_PRICE_ID",
		"stripe.credit_prices.credits_5000.amount":   "STRIPE_CREDITS_5000_AMOUNT",
		"stripe.credit_prices.credits_5000.currency": "STRIPE_CREDITS_5000_CURRENCY",
	}

	for configKey, environmentVariable := range creditPriceBindings {
		if err := viper.BindEnv(configKey, environmentVariable); err != nil {
			return fmt.Errorf(
				"failed to bind %s: %w",
				configKey,
				err,
			)
		}
	}

	return nil
}

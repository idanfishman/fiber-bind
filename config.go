package validation

import (
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

// ValidationSource indicates the source of the data to validate.
const (
	// Body indicates that the data to validate is in the request body.
	Body = "body"
	// Query indicates that the data to validate is in the query string.
	Query = "query"
	// Params indicates that the data to validate is in the route parameters.
	Params = "params"
)

// Config defines the config for middleware.
type Config struct {
	// Next defines a function to skip this middleware when returned true.
	//
	// Optional. Default: nil
	Next func(c *fiber.Ctx) bool
	// Validator defines the validator instance to use. It is recommended
	// to provide your own instance to avoid import cycles and to be able
	// too add custom validations.
	//
	// Required. Default: validator.New()
	Validator *validator.Validate
	// Source defines the source of the data to validate.
	//
	// Required. Default: body
	Source string
	// FormFiles defines the form files fields of the data to validate.
	// The key is the name of the struct field and the value is the name of the
	// file to upload.
	//
	// Optional. Default: nil
	FormFileFields map[string]string
}

// ConfigDefault is the default config
var ConfigDefault = Config{
	Next:           nil,
	Validator:      validator.New(),
	Source:         Body,
	FormFileFields: nil,
}

// Helper function to set default values
func configDefault(config ...Config) Config {
	// Return default config if no config provided
	if len(config) < 1 {
		return ConfigDefault
	}

	// Override default config
	cfg := config[0]

	// Set default values
	if cfg.Next == nil {
		cfg.Next = ConfigDefault.Next
	}

	if cfg.Validator == nil {
		cfg.Validator = ConfigDefault.Validator
	}

	if cfg.Source == "" {
		cfg.Source = ConfigDefault.Source
	}

	if cfg.FormFileFields == nil {
		cfg.FormFileFields = ConfigDefault.FormFileFields
	}

	return cfg
}

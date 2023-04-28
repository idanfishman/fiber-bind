package validation

import (
	"reflect"

	"github.com/gofiber/fiber/v2"
)

const Version = "0.1.0"

// New creates a new middleware handler
func New(config Config, schema interface{}) fiber.Handler {
	// Set default config
	cfg := configDefault(config)

	return func(c *fiber.Ctx) error {
		// Don't execute middleware if Next returns true
		if cfg.Next != nil && cfg.Next(c) {
			return c.Next()
		}

		var data interface{}
		var err error

		switch cfg.Source {
		case Body:
			data = reflect.New(reflect.TypeOf(schema).Elem()).Interface()
			err = c.BodyParser(data)
		case Query:
			data = reflect.New(reflect.TypeOf(schema).Elem()).Interface()
			err = c.QueryParser(data)
		case Params:
			data = reflect.New(reflect.TypeOf(schema).Elem()).Interface()
			err = c.ParamsParser(data)
		default:
			return fiber.ErrInternalServerError
		}

		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": err.Error(),
			})
		}

		if err := cfg.Validator.Struct(data); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": err.Error(),
			})
		}

		c.Locals(cfg.Source.String(), data)

		return c.Next()
	}
}

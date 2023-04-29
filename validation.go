package validation

import (
	"reflect"
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

const Version = "1.0.0"

// New creates a new middleware handler
func New(config Config, schema interface{}) fiber.Handler {
	// Set default config
	cfg := configDefault(config)

	// Return the middleware handler function
	return func(c *fiber.Ctx) error {
		// Skip middleware execution if Next function returns true
		if cfg.Next != nil && cfg.Next(c) {
			return c.Next()
		}

		// Parse incoming data based on the configured source
		var data interface{}
		var err error
		switch cfg.Source {
		case Body, Form:
			// Parse request body and store it in the data variable
			data = reflect.New(reflect.TypeOf(schema).Elem()).Interface()
			err = c.BodyParser(data)
		case Query:
			// Parse query string parameters and store them in the data variable
			data = reflect.New(reflect.TypeOf(schema).Elem()).Interface()
			err = c.QueryParser(data)
		case Params:
			// Parse route parameters and store them in the data variable
			data = reflect.New(reflect.TypeOf(schema).Elem()).Interface()
			err = c.ParamsParser(data)
		default:
			// Return an internal server error if the source is not recognized
			return fiber.ErrInternalServerError
		}

		// Return a bad request error if the data could not be parsed
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": err.Error(),
			})
		}

		// Validate the data using the configured validator instance and the provided schema
		if err := cfg.Validator.Struct(data); err != nil {
			// Map validation errors to a response object
			response := mapValidationErrors(err, cfg.Source, schema)
			// Return a bad request error with the validation errors
			return c.Status(fiber.StatusBadRequest).JSON(response)
		}

		// Add the validated data to the context locals
		c.Locals(cfg.Source, data)

		// Continue to the next middleware in the chain
		return c.Next()
	}
}

// mapValidationErrors maps validation errors to a response object
func mapValidationErrors(err error, source string, schema interface{}) fiber.Map {
	errors := fiber.Map{}
	// Iterate over each validation error
	for _, err := range err.(validator.ValidationErrors) {
		// Get the name of the field that failed validation
		fieldName := strings.Split(err.StructNamespace(), ".")[1]
		// Get the validation tag for the field
		field, _ := reflect.TypeOf(schema).Elem().FieldByName(fieldName)
		tag := field.Tag.Get(sourceTags[source])
		// Get the validation error message
		errorMessage := err.Tag()
		// Add the error message to the response object
		errors[tag] = errorMessage
	}
	return errors
}

// sourceTags maps data sources (body, query, and params) to validation tags (json, query, and params)
var sourceTags = map[string]string{
	"body":   "json",
	"form":   "form",
	"params": "params",
	"query":  "query",
}

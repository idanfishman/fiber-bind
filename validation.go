package validation

import (
	"fmt"
	"mime/multipart"
	"reflect"
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

const Version = "1.2.0-alpha.4"

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
		case Body:
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
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": fmt.Sprintf("Unrecognized data source: %s", cfg.Source),
			})
		}

		// Return a bad request error if the data could not be parsed
		if err != nil {
			return c.Status(fiber.StatusUnprocessableEntity).JSON(fiber.Map{
				"error": err.Error(),
			})
		}

		// Extract form files from the request body and add them to the data variable
		if cfg.Source == Body && cfg.FormFileFields != nil {
			// Get the multipart form from the request body
			form, err := c.MultipartForm()
			if err != nil {
				return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
					"error": fmt.Sprintf("Failed to parse multipart form: %s", err.Error()),
				})
			}

			// Iterate over each form file field and add the files to the data variable
			dataValue := reflect.ValueOf(data).Elem()

			for _, field := range cfg.FormFileFields {
				formfiles, ok := form.File[field]
				if ok {
					structField := dataValue.FieldByName(field)

					// Check if the field is a slice or a pointer
					switch structField.Kind() {
					case reflect.Ptr:
						if len(formfiles) > 0 {
							structField.Set(reflect.ValueOf(formfiles[0]))
						}
					case reflect.Slice:
						files := make([]*multipart.FileHeader, len(formfiles))
						copy(files, formfiles)
						structField.Set(reflect.ValueOf(files))
					default:
						return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
							"error": fmt.Sprintf("Unsupported field type for %s: %s", field, structField.Kind()),
						})
					}
				}
			}
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
	"query":  "query",
	"body":   "form",
	"params": "params",
}

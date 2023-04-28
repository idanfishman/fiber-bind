package validation_test

import (
	"bytes"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/idan-fishman/validation"
	"github.com/stretchr/testify/assert"
)

type User struct {
	Name  string `json:"name" validate:"required"`
	Email string `json:"email" validate:"required,email"`
	Age   int    `json:"age" validate:"gte=0,lte=130"`
}

func TestValidationMiddleware(t *testing.T) {
	// Initialize a new validator
	v := validator.New()

	// Create a new Fiber instance
	app := fiber.New()

	// Apply the validation middleware
	app.Use(validation.New(validation.Config{
		Validator: v,
		Source:    validation.Body,
	}, &User{}))

	// Define a POST endpoint for testing
	app.Post("/user", func(c *fiber.Ctx) error {
		fmt.Println(c.Locals(validation.Body))
		fmt.Println(c.Locals("body"))
		user := c.Locals(validation.Body).(*User)
		return c.JSON(user)
	})

	// Prepare a valid request
	req := httptest.NewRequest(http.MethodPost, "/user", bytes.NewBufferString(`{"name": "Test User", "email": "test@user.com", "age": 30}`))
	req.Header.Set(fiber.HeaderContentType, fiber.MIMEApplicationJSON)

	resp, _ := app.Test(req)

	// Assert HTTP Status OK for valid request
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	// Prepare an invalid request
	reqInvalid := httptest.NewRequest(http.MethodPost, "/user", bytes.NewBufferString(`{"name": "", "email": "invalid", "age": 200}`))
	reqInvalid.Header.Set(fiber.HeaderContentType, fiber.MIMEApplicationJSON)

	respInvalid, _ := app.Test(reqInvalid)

	// Assert HTTP Status BadRequest for invalid request
	assert.Equal(t, http.StatusBadRequest, respInvalid.StatusCode)
}

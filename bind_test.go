package bind_test

import (
	"bytes"
	"encoding/json"
	"mime/multipart"
	"net/http/httptest"
	"testing"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	bind "github.com/idan-fishman/fiber-bind"
	"github.com/stretchr/testify/assert"
)

type Person struct {
	Name string `json:"name" validate:"required"`
	Age  int    `json:"age" validate:"gte=18"`
}

// Test successful validation
func TestJSONSuccess(t *testing.T) {
	// Arrange
	app := fiber.New()
	app.Post("/person", bind.New(bind.Config{
		Validator: validator.New(),
		Source:    bind.JSON,
	}, &Person{}), func(c *fiber.Ctx) error {
		person := c.Locals(bind.JSON).(*Person)
		return c.JSON(person)
	})

	body, _ := json.Marshal(&Person{Name: "John", Age: 20})
	req := httptest.NewRequest("POST", "/person", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	// Act
	resp, _ := app.Test(req)

	// Assert
	assert.Equal(t, 200, resp.StatusCode, "They should be equal")
}

// Test failed validation
func TestJSONFailure(t *testing.T) {
	// Arrange
	app := fiber.New()
	app.Post("/person", bind.New(bind.Config{
		Validator: validator.New(),
		Source:    bind.JSON,
	}, &Person{}), func(c *fiber.Ctx) error {
		person := c.Locals(bind.JSON).(*Person)
		return c.JSON(person)
	})

	body, _ := json.Marshal(&Person{Name: "", Age: 20})
	req := httptest.NewRequest("POST", "/person", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	// Act
	resp, _ := app.Test(req)

	// Assert
	assert.NotEqual(t, 200, resp.StatusCode, "They should not be equal")
}

type QueryParams struct {
	Page  int `json:"page" validate:"gte=1"`
	Limit int `json:"limit" validate:"required,gte=1,lte=100"`
}

// Test successful validation of query parameters
func TestQueryParamsSuccess(t *testing.T) {
	// Arrange
	app := fiber.New()
	app.Get("/items", bind.New(bind.Config{
		Validator: validator.New(),
		Source:    bind.Query,
	}, &QueryParams{}), func(c *fiber.Ctx) error {
		params := c.Locals(bind.Query).(*QueryParams)
		return c.JSON(params)
	})

	req := httptest.NewRequest("GET", "/items?page=2&limit=50", nil)

	// Act
	resp, _ := app.Test(req)

	// Assert
	assert.Equal(t, 200, resp.StatusCode, "They should be equal")
}

// Test failed validation of query parameters
func TestQueryParamsFailure(t *testing.T) {
	// Arrange
	app := fiber.New()
	app.Get("/items", bind.New(bind.Config{
		Validator: validator.New(),
		Source:    bind.Query,
	}, &QueryParams{}), func(c *fiber.Ctx) error {
		params := c.Locals(bind.Query).(*QueryParams)
		return c.JSON(params)
	})

	req := httptest.NewRequest("GET", "/items?page=0&limit=150", nil)

	// Act
	resp, _ := app.Test(req)

	// Assert
	assert.NotEqual(t, 200, resp.StatusCode, "They should not be equal")
}

type RouteParams struct {
	ID int `json:"id" validate:"required,gte=1"`
}

// Test successful validation of route parameters
func TestRouteParamsSuccess(t *testing.T) {
	// Arrange
	app := fiber.New()
	app.Get("/item/:id", bind.New(bind.Config{
		Validator: validator.New(),
		Source:    bind.Params,
	}, &RouteParams{}), func(c *fiber.Ctx) error {
		params := c.Locals(bind.Params).(*RouteParams)
		return c.JSON(params)
	})

	req := httptest.NewRequest("GET", "/item/123", nil)

	// Act
	resp, _ := app.Test(req)

	// Assert
	assert.Equal(t, 200, resp.StatusCode, "They should be equal")
}

// Test failed validation of route parameters
func TestRouteParamsFailure(t *testing.T) {
	// Arrange
	app := fiber.New()
	app.Get("/item/:id", bind.New(bind.Config{
		Validator: validator.New(),
		Source:    bind.Params,
	}, &RouteParams{}), func(c *fiber.Ctx) error {
		params := c.Locals(bind.Params).(*RouteParams)
		return c.JSON(params)
	})

	req := httptest.NewRequest("GET", "/item/0", nil)

	// Act
	resp, _ := app.Test(req)

	// Assert
	assert.NotEqual(t, 200, resp.StatusCode, "They should not be equal")
}

type FormData struct {
	Title string `form:"title" validate:"required"`
}

// Test successful validation of form data without files
func TestFormDataSuccess(t *testing.T) {
	// Arrange
	app := fiber.New()
	app.Post("/upload", bind.New(bind.Config{
		Validator: validator.New(),
		Source:    bind.Form,
	}, &FormData{}), func(c *fiber.Ctx) error {
		form := c.Locals(bind.Form).(*FormData)
		return c.JSON(form)
	})

	body := new(bytes.Buffer)
	writer := multipart.NewWriter(body)
	err := writer.WriteField("title", "A Great Title")
	assert.NoError(t, err)
	writer.Close()

	req := httptest.NewRequest("POST", "/upload", body)
	req.Header.Set("Content-Type", writer.FormDataContentType())

	// Act
	resp, _ := app.Test(req)

	// Assert
	assert.Equal(t, 200, resp.StatusCode, "They should be equal")
}

// Test failed validation of form data without files
func TestFormDataFailure(t *testing.T) {
	// Arrange
	app := fiber.New()
	app.Post("/upload", bind.New(bind.Config{
		Validator: validator.New(),
		Source:    bind.Form,
	}, &FormData{}), func(c *fiber.Ctx) error {
		form := c.Locals(bind.Form).(*FormData)
		return c.JSON(form)
	})

	body := new(bytes.Buffer)
	writer := multipart.NewWriter(body)
	err := writer.WriteField("title", "")
	assert.NoError(t, err)
	writer.Close()

	req := httptest.NewRequest("POST", "/upload", body)
	req.Header.Set("Content-Type", writer.FormDataContentType())

	// Act
	resp, _ := app.Test(req)

	// Assert
	assert.NotEqual(t, 200, resp.StatusCode, "They should not be equal")
}

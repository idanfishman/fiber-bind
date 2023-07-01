# Bind Middleware for Fiber

[![Mentioned in Awesome Fiber](https://awesome.re/mentioned-badge.svg)](https://github.com/gofiber/awesome-fiber)

Bind is a request schema validator middleware for the [Fiber](https://github.com/gofiber/fiber) web framework. It provides a convenient way to parse and validate data from different sources such as the request body, query string parameters, and route parameters.

## Installation

Use go get to install the middleware:

```bash
go get -u github.com/idan-fishman/fiber-bind
```

## Usage Example

Here are examples of how to use Bind to validate data from the request body and query parameters.

### Body Parameters

```go
package main

import (
    "github.com/gofiber/fiber/v2"
    "github.com/idan-fishman/fiber-bind"
)

type Person struct {
    Name string `json:"name" validate:"required"`
    Age  int    `json:"age" validate:"gte=18"`
}

func main() {
    app := fiber.New()

    app.Post("/person", bind.New(bind.Config{
        Validator: validator.New(),
        Source:    bind.JSON,
    }, &Person{}), func(c *fiber.Ctx) error {
        person := c.Locals(bind.JSON).(*Person)
        return c.JSON(person)
    })

    app.Listen(":3000")
}
```

### Query Parameters

```go
package main

import (
    "github.com/gofiber/fiber/v2"
    "github.com/idan-fishman/fiber-bind"
)

type PersonParams struct {
    Name string `json:"name" validate:"required"`
}

func main() {
    app := fiber.New()

    app.Get("/person", bind.New(bind.Config{
        Validator: validator.New(),
        Source:    bind.Query,
    }, &PersonParams{}), func(c *fiber.Ctx) error {
        params := c.Locals(bind.Query).(*PersonParams)
        return c.JSON(params)
    })

    app.Listen(":3000")
}
```

In these examples, we define a struct that represents the data we want to validate. We then use Bind to validate the data from the request body or query parameters using the `New` function. We specify the source of the data as `bind.Body` or `bind.Query` and provide the struct as the schema to validate against. If the data is valid, it will be available in the context locals for further use.

## Configuration

Bind can be configured using the `Config` struct, which has the following fields:

- `Next` - A function that is called before the middleware is executed. If this function returns `true`, the middleware is skipped.
- `Validator` - A validator instance to use for validating the data. By default, the middleware uses the [go-playground/validator/v10](https://github.com/go-playground/validator/v10) package.
- `Source` - The source of the data to validate. This can be one of `bind.JSON`, `bind.XML`, `bind.Form`, `bind.Query`, or `bind.Params`.
- `FormFileFields` - A map where the key is the name of the file field in the form and the value is the form tag. If a field is specified here, the middleware will check if the file is uploaded in the request.

## License

Bind is licensed under the MIT License. See the [LICENSE](LICENSE) file for details.

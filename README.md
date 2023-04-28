# Validation Middleware for Fiber

The validation middleware is a middleware for [Fiber](https://github.com/gofiber/fiber) web framework that provides a convenient way to parse and validate data from different sources such as the request body, query string parameters, and route parameters.

### Installation

Use go get to install the middleware:

```bash
go get -u github.com/idan-fishman/validation
```

### Usage Example

Here is an example of how to use the middleware to validate data from the request body:

```go
package main

import (
    "github.com/gofiber/fiber/v2"
    "github.com/idan-fishman/validation"
)

type Person struct {
    Name string `json:"name" validate:"required"`
    Age  int    `json:"age" validate:"gte=18"`
}

func main() {
    app := fiber.New()

    app.Post("/person", validation.New(validation.Config{
        Validator: validator.New(),
        Source:    validation.Body,
    }, Person{}), func(c *fiber.Ctx) error {
        person := c.Locals(validation.Body).(Person)
        return c.JSON(person)
    })

    app.Listen(":3000")
}
```

In this example, we define a `Person` struct that represents the data we want to validate. We then use the middleware to validate the data in the request body using the `New` function. We specify the source of the data as `validation.Body` and provide the `Person` struct as the schema to validate against. If the data is valid, it will be available in the context locals for further use.

### Configuration

The middleware can be configured using the `Config` struct, which has the following fields:

- `Next` - A function that is called before the middleware is executed. If this function returns `true`, the middleware is skipped.
- `Validator` - A validator instance to use for validating the data. By default, the middleware uses the [go-playground/validator/v10](https://github.com/go-playground/validator/v10) package.
- `Source` - The source of the data to validate. This can be one of `validation.Body`, `validation.Query`, or `validation.Params`.

## License

This middleware is licensed under the MIT License. See the [LICENSE](LICENSE) file for details.

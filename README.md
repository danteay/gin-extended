# Gin-Gonic Extended Commons

this package stores usefully middlewares and libraries that could be implemented
on a Gin project

## Installation

```bash
go get -u -v github.com/danteay/gin-extended-commons
```

## Content

* [Middlewares](#middlewares)
  * [Authentication](#authentication)
  * [SwaggerValidator](#swagger-validator)
  * [Zerolog](#zerolog)
* [Libraries](#libraries)
  * [Parse](#parse)
    * [ParseYmlFile](#parse-yml-file)
    * [ShouldParseYmlFile](#should-parse-yml-file)

<span id="middlewares"></span>

## Middlewares

<span id="authentication"></span>

### Authentication

This middleware can authenticate request with by two different strategies; as a
`bearer` token or with a `basic` authentication schema.

```go
package main

import (
    "github.com/gin-gonic/gin"
    mw "github.com/danteay/gin-extended-commons/middlewares"
)

func main(){
    app := gin.New()

    app.Use(mw.Authentication())

    // ....

    app.Run()
}
```

You can set custom configuration for this middleware:

```go
package main

import (
    "github.com/gin-gonic/gin"
    mw "github.com/danteay/gin-extended-commons/middlewares"
)

func main(){
    app := gin.New()

    authConf := &mw.AuthorizationConfig{
        Type:   "bearer",
        APIKey: "123456789",
    }

    app.Use(mw.AuthenticationWithConfig())

    // ....

    app.Run()
}
```

Also you can config as basic authentication.

```go
package main

import (
    "github.com/gin-gonic/gin"
    mw "github.com/danteay/gin-extended-commons/middlewares"
)

func main(){
    app := gin.New()

    authConf := &mw.AuthorizationConfig{
        Type:            "basic",
        AuthCredentials: []string{"user", "password"},
    }

    app.Use(mw.AuthenticationWithConfig())

    // ....

    app.Run()
}
```

Or you can define a non static validation for the authentication

```go
package main

import (
    "github.com/gin-gonic/gin"
    mw "github.com/danteay/gin-extended-commons/middlewares"
)

func main(){
    app := gin.New()

    authConf := &mw.AuthorizationConfig{
        Type:      "bearer",
        Validator: func (key string) bool {
            var res bool

            // validate key

            return res
        },
    }

    app.Use(mw.AuthenticationWithConfig())

    // ....

    app.Run()
}
```

<span id="swagger-validator"></span>

### SwaggerValidator

This middleware validate the api request and response schema with the OpenApi
specification. By default the middleware search for a file called `spec.yml`
to load the API specification.

```go
package main

import (
    "github.com/gin-gonic/gin"
    mw "github.com/danteay/gin-extended-commons/middlewares"
)

func main(){
    app := gin.New()

    app.Use(mw.Authentication())

    // ....

    app.Run()
}
```

Also you can specify a route to load the specification.

```go
package main

import (
    "github.com/gin-gonic/gin"
    mw "github.com/danteay/gin-extended-commons/middlewares"
)

func main(){
    app := gin.New()

    app.Use(mw.AuthenticationWithConfig(&mw.AuthorizationConfig{
        Document: "my_spec_file.yml"
    }))

    // ....

    app.Run()
}
```

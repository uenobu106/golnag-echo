package main

import (
	"io"
	"net/http"
	"os"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	e := echo.New()

	type User struct {
		Name  string `json:"name" xml:"name" form:"name" query:"name"`
		Email string `json:"email" xml:"email" form:"email" query:"email"`
	}
	

	// 	$ curl -F "name=Foo bar" -F "email=hoge@huga.com" http://localhost:1323/users
	// => {"name":"Foo bar","email":"hoge@huga.com"}
	e.POST("/users", func(c echo.Context) error {
		u := new(User)
		if err := c.Bind(u); err != nil {
			return err
		}
		return c.JSON(http.StatusCreated, u)
		// or
		// return c.XML(http.StatusCreated, u)
	})
	

	
	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello, World!")
	})

	e.GET("/users/:id", getUser)
	e.GET("/show", show)
  e.POST("/save", save)

	// Root level middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	// Group level middleware
	g := e.Group("/admin")
	g.Use(middleware.BasicAuth(func(username, password string, c echo.Context) (bool, error) {
		if username == "hoge" && password == "hoge123" {
			return true, nil
		}
		return false, nil
	}))

	// Route level middleware
	track := func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			println("request to /users")
			return next(c)
		}
	}
	e.GET("/users", func(c echo.Context) error {
		return c.String(http.StatusOK, "/users")
	}, track)

	e.Logger.Fatal(e.Start(":1323"))

}

// Routing
// e.POST("/users", saveUser)
// e.GET("/users/:id", getUser)
// e.PUT("/users/:id", updateUser)
// e.DELETE("/users/:id", deleteUser)

// e.GET("/users/:id", getUser)
func getUser(c echo.Context) error {
	// User ID from path `users/:id`
	id := c.Param("id")
	return c.String(http.StatusOK, id)
}

//e.GET("/show", show)
func show(c echo.Context) error {
	// Get team and member from the query string
	team := c.QueryParam("team")
	member := c.QueryParam("member")
	return c.String(http.StatusOK, "team:" + team + ", member:" + member)
}

// e.POST("/save", save)
// func save(c echo.Context) error {
// 	// Get name and email
// 	name := c.FormValue("name")
// 	email := c.FormValue("email")
// 	return c.String(http.StatusOK, "name:" + name + ", email:" + email)
// }

// $ curl -F "name=Hoge" -F "avatar=@avatar.png" http://localhost:1323/save
// => <b>Thank you! Hoge</b>%
func save(c echo.Context) error {
	// Get name
	name := c.FormValue("name")
	// Get avatar
  avatar, err := c.FormFile("avatar")
  if err != nil {
 		return err
 	}
 
 	// Source
 	src, err := avatar.Open()
 	if err != nil {
 		return err
 	}
 	defer src.Close()
 
 	// Destination
 	dst, err := os.Create(avatar.Filename)
 	if err != nil {
 		return err
 	}
 	defer dst.Close()
 
 	// Copy
 	if _, err = io.Copy(dst, src); err != nil {
  	return err
  }

	return c.HTML(http.StatusOK, "<b>Thank you! " + name + "</b>")
}


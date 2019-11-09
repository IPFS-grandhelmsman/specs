package main

import (
	"os"

	"github.com/gwaylib/errors"
	"github.com/gwaylib/eweb"
	"github.com/gwaylib/log"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)

var r = eweb.Default()

func init() {
	// init route for static file
	r.Static("/", "./build/website/")
}

func main() {
	r.Debug = os.Getenv("GIN_MODE") != "release"

	// middle ware
	r.Use(middleware.Gzip())

	// 过滤器
	r.Use(func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			req := c.Request()
			uri := req.URL.Path

			// 静态页面过虑
			switch {
			case uri == "/hacheck":
				return c.String(200, "1")
			}

			// 校验通过
			return next(c)
		}
	})

	port := ":1313"
	log.Debugf("start at: %s", port)
	if err := r.Start(port); err != nil {
		log.Exit(1, errors.As(err))
		return
	}
}

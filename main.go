package main

import (
	"fullstack-journal/app/kernel"

	"github.com/labstack/echo/v4"
)

func main() {
	e := echo.New()
	kernel.StartApplication(e)
}

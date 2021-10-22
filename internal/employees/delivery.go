package employees

import "github.com/labstack/echo/v4"

type Handler interface {
	Create()  echo.HandlerFunc
	Update()  echo.HandlerFunc
	Delete()  echo.HandlerFunc
	List()    echo.HandlerFunc
	GetById() echo.HandlerFunc
}

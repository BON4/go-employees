package tasks

import "github.com/labstack/echo/v4"

type Handler interface {
	Create()             echo.HandlerFunc
	Update()			 echo.HandlerFunc
	GetByTaskId()        echo.HandlerFunc
	DeleteByTaskId()     echo.HandlerFunc
	DeleteByEmployeeId() echo.HandlerFunc
	List() 				 echo.HandlerFunc
	GetByEmployeeId()    echo.HandlerFunc
}

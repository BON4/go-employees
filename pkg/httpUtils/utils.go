package httpUtils

import (
	"context"
	"github.com/labstack/echo/v4"
	"github.com/go-playground/validator/v10"
)

var validate *validator.Validate

func init(){
	validate = validator.New()
}

func ValidateStruct(ctx context.Context, s interface{}) error {
	return validate.StructCtx(ctx, s)
}

func ReadRequest(ctx echo.Context, req interface{}) error {
	if err := ctx.Bind(req); err != nil{
		return err
	}
	return ValidateStruct(ctx.Request().Context(), req)
}
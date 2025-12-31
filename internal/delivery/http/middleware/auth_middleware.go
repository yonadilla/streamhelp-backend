package middleware

import (
	"streamhelper-backend/internal/model"
	"streamhelper-backend/internal/usecase"
	"streamhelper-backend/internal/util"

	"github.com/gofiber/fiber/v2"
)

func NewAuth(userUserCase *usecase.UserUseCase, tokenUtli *util.TokenUtil) fiber.Handler{
	return func(ctx *fiber.Ctx) error  {
		request := &model.VerifyUserRequest{Token: ctx.Get("Authorization", "NOT_FOUND")}
		userUserCase.Log.Debugf("Authorization :%s", request.Token )

		auth , err := tokenUtli.ParseToken(ctx.UserContext(), request.Token)
		if err != nil {
			userUserCase.Log.Warnf("Failed find user by token : %+v", err)
			return fiber.ErrUnauthorized
		}

		userUserCase.Log.Debugf("User : %+v", auth.ID)
		ctx.Locals("auth", auth)
		return ctx.Next()
	}
}

func GetUser(ctx *fiber.Ctx) *model.Auth{
	return ctx.Locals("auth").(*model.Auth)
}
package utils

import (
	"strconv"

	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/middleware/jwt"
)

func UserIDMiddleware(ctx iris.Context) {
	params := ctx.Params()
	id := params.Get("id")

	claims := jwt.Get(ctx).(*AccessToken)

	userID := strconv.FormatUint(uint64(claims.ID), 10)

	if userID != id {
		ctx.StatusCode(iris.StatusForbidden)
		return
	}
	ctx.Next()
}

package middleware

import (
	"geektime-basic/webook/internal/web"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"log"
	"net/http"
	"strings"
	"time"
)

type JWTLoginMiddlewareBuilder struct {
	paths []string
}

func NewLoginJWTMiddlewareBuilder() *JWTLoginMiddlewareBuilder {
	return &JWTLoginMiddlewareBuilder{}
}

func (j *JWTLoginMiddlewareBuilder) IgnorePaths(path string) *JWTLoginMiddlewareBuilder {
	j.paths = append(j.paths, path)

	return j
}

func (j *JWTLoginMiddlewareBuilder) Build() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		//if ctx.Request.URL.Path == "/users/signup" ||
		//	ctx.Request.URL.Path == "/users/login" {
		//	return
		//}

		for _, path := range j.paths {
			if ctx.Request.URL.Path == path {
				return
			}
		}

		// 用JWT来校验
		tokenHeader := ctx.GetHeader("Authorization")
		if tokenHeader == "" {
			// 没登录
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		authSegments := strings.SplitN(tokenHeader, " ", 2)
		if len(authSegments) != 2 {
			// 没登录，有人瞎搞
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		tokenStr := authSegments[1]
		uc := web.UserClaims{}
		token, err := jwt.ParseWithClaims(tokenStr, &uc, func(token *jwt.Token) (interface{}, error) {
			return web.JWTKey, nil
		})
		if err != nil || !token.Valid {
			// 没登录
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		expireTime, err := uc.GetExpirationTime()
		if err != nil {
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}
		if expireTime.Before(time.Now()) {
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		if ctx.GetHeader("User-Agent") != uc.UserAgent {
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		if expireTime.Sub(time.Now()) < time.Second*50 {
			uc.ExpiresAt = jwt.NewNumericDate(time.Now().Add(time.Minute))
			newToken, err := token.SignedString(web.JWTKey)
			if err != nil {
				log.Println(err)
			} else {
				ctx.Header("x-jwt-token", newToken)
			}
		}

		ctx.Set("user", uc)

	}
}

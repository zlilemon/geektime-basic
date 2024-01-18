package ioc

import (
	"geektime-basic/webook/internal/web"
	"geektime-basic/webook/internal/web/middleware"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"strings"
	"time"
)

func InitGin(mdls []gin.HandlerFunc, hdl *web.UserHandler) *gin.Engine {
	server := gin.Default()
	server.Use(mdls...)
	hdl.RegisterRoutes(server)

	return server
}

func InitMiddlewares() []gin.HandlerFunc {
	return []gin.HandlerFunc{
		cors.New(cors.Config{
			// 是否允许带cookie 之类的东西
			AllowCredentials: true,
			AllowHeaders:     []string{"Content-Type", "Authorization"},
			// 不加这个 ExposeHeaders，前端是那不到 x-jwt-token的
			// ExposeHeaders: []string{"x-jwt-token"},
			AllowOriginFunc: func(origin string) bool {
				if strings.HasPrefix(origin, "http://localhosts") {
					return true
				}
				return strings.Contains(origin, "your_company.com")
			},
			MaxAge: 12 * time.Hour,
		}),
		middleware.NewLoginJWTMiddlewareBuilder().
			IgnorePaths("/users/signup").
			IgnorePaths("/users/login").
			IgnorePaths("/users/loginJwt").
			IgnorePaths("/users/profile").
			IgnorePaths("/users/login_sms/code/send").
			Build(),
	}
}

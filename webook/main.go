package main

import (
	"geektime-basic/webook/config"
	"geektime-basic/webook/internal/repository"
	"geektime-basic/webook/internal/repository/cache"
	"geektime-basic/webook/internal/repository/dao"
	"geektime-basic/webook/internal/service"
	"geektime-basic/webook/internal/web"
	"geektime-basic/webook/internal/web/middleware"
	"github.com/gin-contrib/cors"
	//"github.com/gin-contrib/sessions/redis"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"strings"
	"time"
)

func main() {

	db := initDB()

	server := initWebServer()

	rdb := initRedis()

	initUser(server, db, rdb)

	//u := &web.UserHandler{}
	//u.RegisterRoutes(server)

	//server := gin.Default()
	//server.GET("/hello", func(ctx *gin.Context) {
	//	ctx.String(http.StatusOK, "你来了")
	//})
	server.Run(":8080")

}

func initDB() *gorm.DB {
	db, err := gorm.Open(mysql.Open("root:root@tcp(localhost:13316)/webook"))
	if err != nil {
		panic(err)
	}

	err = dao.InitTables(db)
	if err != nil {
		panic(err)
	}
	return db
}

func initWebServer() *gin.Engine {
	server := gin.Default()

	server.Use(func(ctx *gin.Context) {
		println("这是第一个 middleware for test")
	})

	server.Use(func(ctx *gin.Context) {
		println("这是第二个 middleware for test")
	})

	server.Use(func(ctx *gin.Context) {
		println("这是第三个 middleware for test")
	})
	/*
		cmd := redis.NewClient(&redis.Options{
			Addr:     "localhost:6379",
			Password: "",
			DB:       1,
		})

		server.Use(ratelimit.NewBuilder(cmd, time.Minute, 100).Build())
	*/

	server.Use(cors.New(cors.Config{
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
	}))

	/*

		store, err := redis.NewStore(16,
			"tcp", "localhost:6379", "",
			[]byte("95osj3fUD7fo0mlYdDbncXz4VD2igvf0"), []byte("0Pf2r0wZBpXVXlQNdpwCXN4ncnlnZSc3"))
		// usingJWT(server)

		if err != nil {
			panic(err)
		}
	*/

	// server.Use(sessions.Sessions("mysession", store))

	server.Use(middleware.NewLoginJWTMiddlewareBuilder().
		IgnorePaths("/users/signup").
		IgnorePaths("/users/login").
		IgnorePaths("/users/loginJwt").
		IgnorePaths("/users/profile").
		Build())
	//server.Use(middleware.)

	return server
}

func usingJWT(server *gin.Engine) {
	mldBd := &middleware.JWTLoginMiddlewareBuilder{}
	server.Use(mldBd.Build())
}

func initRedis() redis.Cmdable {
	redisClient := redis.NewClient(&redis.Options{
		Addr: config.Config.Redis.Addr,
	})

	return redisClient
}

func initUser(server *gin.Engine, db *gorm.DB, rdb redis.Cmdable) {
	ud := dao.NewUserDAO(db)
	uc := cache.NewUserCache(rdb)

	ur := repository.NewUserRepostiry(ud, uc)
	us := service.NewUserService(ur)
	c := web.NewUserHandler(us)

	c.RegisterRoutes(server)
}

func hello() {

}

//go:build wireinject

package main

import (
	"geektime-basic/webook/internal/repository"
	"geektime-basic/webook/internal/repository/cache"
	"geektime-basic/webook/internal/repository/dao"
	"geektime-basic/webook/internal/service"
	"geektime-basic/webook/internal/web"
	"geektime-basic/webook/ioc"
	"github.com/gin-gonic/gin"
	"github.com/google/wire"
)

func InitWebServer() *gin.Engine {
	wire.Build(ioc.InitDB, ioc.InitRedis,
		dao.NewUserDAO,
		cache.NewUserCache, cache.NewCodeCache,
		repository.NewCodeRepository, repository.NewUserRepostiry,
		service.NewUserService, service.NewCodeService,
		ioc.InitSMSService,
		web.NewUserHandler,
		ioc.InitGin,
		ioc.InitMiddlewares,
	)
	return new(gin.Engine)
}

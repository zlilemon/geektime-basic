//go:build wireinject

package main

import (
	"geektime-basic/wire/repository"
	"geektime-basic/wire/repository/dao"
	"github.com/google/wire"
)

func InitRepository() *repository.UserRepository {
	// 这个方法里面，传入各个组件的初始化方法
	wire.Build(repository.NewUserRepository,
		dao.NewUserDAO,
		InitDB)

	return new(repository.UserRepository)
}

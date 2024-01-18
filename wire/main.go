package main

import (
	"fmt"
)

func main() {
	/*
		db, err := gorm.Open(mysql.Open("root:root@tcp(localhost:13316)/webook"))
		if err != nil {
			panic(err)
		}

		ud := dao.NewUserDAO(db)
		repo := repository.NewUserRepository(ud)
	*/
	repo := InitRepository()
	fmt.Println(repo)
	fmt.Println("hello")
}

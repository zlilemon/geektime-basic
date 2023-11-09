package repository

import (
	"context"
	"geektime-basic/webook/internal/domain"
	"geektime-basic/webook/internal/repository/cache"
	"geektime-basic/webook/internal/repository/dao"
)

var ErrUserDuplicateEmail = dao.ErrUserDuplicateEmail
var ErrUserNotFound = dao.ErrUserNotFound

type UserRepository struct {
	dao   *dao.UserDAO
	cache *cache.UserCache
}

func NewUserRepostiry(d *dao.UserDAO, c *cache.UserCache) *UserRepository {
	return &UserRepository{
		dao:   d,
		cache: c,
	}
}

func (ur *UserRepository) Create(ctx context.Context, u domain.User) error {
	err := ur.dao.Insert(ctx, dao.User{
		Email:    u.Email,
		Password: u.Password,
	})

	return err
}

func (ur *UserRepository) FindById(ctx context.Context, id int64) (domain.User, error) {

	// 从cache中查询
	u, err := ur.cache.Get(ctx, id)
	if err == nil {
		return u, nil
	}

	// 没这个数据
	// 从db中查询
	// u, err := ur.dao.FindById(ctx, id)
	u = domain.User{
		Email:    u.Email,
		Password: u.Password,
	}

	go func() {
		err = ur.cache.Set(ctx, u)
		if err != nil {
			// 打日志监控
		}
	}()

	return u, err
}

func (ur *UserRepository) FindByEmail(ctx context.Context, email string) (domain.User, error) {
	u, err := ur.dao.FindByEmail(ctx, email)
	if err != nil {
		return domain.User{}, err
	}

	return domain.User{
		Id:       u.Id,
		Email:    u.Email,
		Password: u.Password,
	}, nil
}

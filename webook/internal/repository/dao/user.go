package dao

import (
	"context"
	"errors"
	"github.com/go-sql-driver/mysql"
	"gorm.io/gorm"
	"time"
)

var ErrUserDuplicateEmail = errors.New("邮件冲突")
var ErrUserNotFound = gorm.ErrRecordNotFound

type User struct {
	Id       int64  `gorm:"primaryKey,autoIncrement"`
	Email    string `gorm:"unique"`
	Password string

	Ctime int64
	Utime int64
}

type UserDAO interface {
	Insert(ctx context.Context, u User) error
	FindById(ctx context.Context, id int64) (User, error)
	FindByEmail(ctx context.Context, email string) (User, error)
}

type GORMUserDAO struct {
	db *gorm.DB
}

func NewUserDAO(db *gorm.DB) UserDAO {
	return &GORMUserDAO{
		db: db,
	}
}

func (ud *GORMUserDAO) Insert(ctx context.Context, u User) error {
	now := time.Now().UnixMilli()
	u.Ctime = now
	u.Utime = now

	err := ud.db.Create(&u).Error
	if me, ok := err.(*mysql.MySQLError); ok {
		const uniqueIndexErrNo uint16 = 1062
		if me.Number == uniqueIndexErrNo {
			return ErrUserDuplicateEmail
		}
	}

	return err
}

func (ud *GORMUserDAO) FindById(ctx context.Context, id int64) (User, error) {
	var u User
	err := ud.db.First(&u, "id = ?", id).Error
	return u, err
}

func (ud *GORMUserDAO) FindByEmail(ctx context.Context, email string) (User, error) {
	var u User
	err := ud.db.WithContext(ctx).Where("email = ?", email).First(&u).Error
	return u, err
}

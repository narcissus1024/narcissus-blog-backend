package dao

import (
	"github.com/narcissus1949/narcissus-blog/internal/database/mysql"
	"github.com/narcissus1949/narcissus-blog/internal/model"
)

var UserDaoInstance = &userDao{}

type userDao struct {
}

func (d *userDao) QueryByUsername(username string) (*model.User, error) {
	var user model.User
	res := mysql.Client.Where("username = ?", username).Find(&user)
	if res.Error != nil {
		return nil, res.Error
	}
	if res.RowsAffected == 0 {
		return nil, nil
	}
	return &user, res.Error
}

func (d *userDao) InsertUser(user *model.User) error {
	res := mysql.Client.Create(user)
	return res.Error
}

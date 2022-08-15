package user

import "shopping/utils/hash"

// 用户service结构体
type Service struct {
	r Repository
}

// 实例化service
func NewUserService(r Repository) *Service {
	r.Migration()
	r.InsertSampleData() //测试数据
	return &Service{
		r: r,
	}
}

// 创建用户
func (c *Service) Create(user *User) error {
	if user.Password != user.Password2 {
		return ErrMismatchedPasswords
	}
	// 用户名存在
	_, err := c.r.GetByName(user.Username)
	if err == nil {
		return ErrUserExistWithName
	}
	err = c.r.Create(user)
	return err
}

// 查询用户
func (c *Service) GetUser(username string, password string) (User, error) {
	user, err := c.r.GetByName(username)
	if err != nil {
		return User{}, ErrUserNotFound
	}
	match := hash.CheckPasswordHash(password+user.Salt, user.Password)
	if !match {
		return User{}, ErrUserNotFound
	}
	return user, nil
}

// 更新用户
func (c *Service) UpdateUser(user *User) error {
	return c.r.Update(user)
}

package storage

import (
	"context"
	userpkg "github.com/opendigitalpay-io/open-balance/internal/user"
	"time"
)

type userModel struct {
	ID         uint64 `gorm:"primary_key"`
	Email      string
	Phone      string
	ExternalID string
	Metadata   []byte
	CreatedAt  int64
	UpdatedAt  int64
}

func (u *userModel) TableName() string {
	return "users"
}

func (u *userModel) model(user userpkg.User) {
	u.ID = user.ID
	u.Email = user.Email
	u.Phone = user.Phone
	u.ExternalID = user.ExternalID
	u.Metadata = user.Metadata
}

func (u *userModel) domain() userpkg.User {
	return userpkg.User{
		ID:         u.ID,
		Email:      u.Email,
		Phone:      u.Phone,
		ExternalID: u.ExternalID,
		Metadata:   u.Metadata,
		CreatedAt:  u.CreatedAt,
		UpdatedAt:  u.UpdatedAt,
	}
}

func (r *Repository) AddUser(ctx context.Context, user userpkg.User) (userpkg.User, error) {
	db := r.DB(ctx)

	var u userModel
	u.model(user)

	now := time.Now().Unix()
	u.CreatedAt = now
	u.UpdatedAt = now

	err := db.Create(&u).Error
	if err != nil {
		return userpkg.User{}, wrapDBError(err, "user")
	}

	user.CreatedAt = u.CreatedAt
	user.UpdatedAt = u.UpdatedAt

	return user, nil
}

func (r *Repository) GetUser(ctx context.Context, userID uint64) (userpkg.User, error) {
	db := r.DB(ctx)

	var u userModel
	err := db.Unscoped().First(&u, userID).Error
	if err != nil {
		return userpkg.User{}, wrapDBError(err, "user")
	}

	user := u.domain()

	return user, nil
}

func (r *Repository) UpdateUser(ctx context.Context, user userpkg.User) (userpkg.User, error) {
	db := r.DB(ctx)

	var u userModel
	u.model(user)

	u.UpdatedAt = time.Now().Unix()

	err := db.Model(&u).Updates(&u).Error
	if err != nil {
		return userpkg.User{}, wrapDBError(err, "user")
	}

	user.UpdatedAt = u.UpdatedAt

	return user, nil
}

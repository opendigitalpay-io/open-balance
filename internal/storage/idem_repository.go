package storage

import (
	"context"
	"encoding/json"
	"time"
)

type idemModel struct {
	ID          uint64 `gorm:"primary_key"`
	IdemId      string
	IsCompleted bool
	Response    []byte
	CreatedAt   int64
	UpdatedAt   int64
}

func (i *idemModel) TableName() string {
	return "idempotency"
}

func (r *Repository) AddIdemEntry(ctx context.Context, id uint64, idemId string)  error {
	db := r.DB(ctx)

	now := time.Now().Unix()
	u := idemModel{
		ID: id,
		IdemId:     idemId,
		CreatedAt:  now,
		UpdatedAt:  now,
	}

	err := db.Create(&u).Error
	if err != nil {
		return wrapDBError(err, "idem")
	}

	return nil
}

func (r *Repository) UpdateIdemEntry(ctx context.Context, idemId string, response interface{}) error {
	db := r.DB(ctx)

	idem := idemModel{
		IsCompleted: true,
		Response: 	 marshallResponse(response),
		UpdatedAt:   time.Now().Unix(),
	}

	err := db.Model(&idem).Where("idem_id = ?", idemId).Updates(&idem).Error
	if err != nil {
		return wrapDBError(err, "idem")
	}

	return nil
}

//todo will use util
func marshallResponse(u interface{}) []byte {
	response, err := json.Marshal(u)
	if err != nil {
		panic(err)
	}
	return response
}

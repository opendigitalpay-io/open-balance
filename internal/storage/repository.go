package storage

import (
	"context"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type Repository struct {
	client *client
}

type client struct {
	db *gorm.DB
}

// NewRepository sets up the database connections using the configuration in the
// process's environment variables. This should be called just once per port
// instance.
func NewRepository(ctx context.Context, config *Config) (*Repository, error) {
	client, err := newClient(ctx, config)
	if err != nil {
		return nil, err
	}

	return &Repository{
		client: client,
	}, nil
}

func newClient(ctx context.Context, config *Config) (*client, error) {
	dsn := "root:root@tcp(127.0.0.1:3306)/open_balance?charset=utf8mb4&parseTime=True" // TODO: move config outside

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	err = db.AutoMigrate(&userModel{}, &rootAccountModel{}, &balanceAccountModel{}, &transactionModel{})
	if err != nil {
		return nil, err
	}

	return &client{
		db: db,
	}, nil
}

func (r *Repository) TxnExec(ctx context.Context, f func(ctxWithTxn context.Context) (interface{}, error)) (interface{}, error) {
	tx := r.client.db.Begin()
	defer func() {
		r := recover()
		if r != nil {
			tx.Rollback()
		}
	}()

	ctxWithTxn := context.WithValue(ctx, "tx", tx)

	res, err := f(ctxWithTxn)
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	tx.Commit()
	return res, nil
}

func (r *Repository) DB(ctx context.Context) *gorm.DB {
	// retrieve tx from ctx, if nil -> return r.db
	tx, ok := ctx.Value("tx").(*gorm.DB)
	if ok {
		return tx
	}
	return r.client.db
}

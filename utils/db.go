package utils

import (
	"context"

	"github.com/pkg/errors"

	"github.com/mrlyc/cmdr/model"
)

func WithTx(ctx context.Context, client *model.Client, fn func(client *model.Client) error) error {
	tx, err := client.Tx(ctx)
	if err != nil {
		return errors.Wrapf(err, "create transaction failed")
	}

	defer func() {
		if v := recover(); v != nil {
			tx.Rollback()
			panic(v)
		}
	}()

	if err := fn(tx.Client()); err != nil {
		rerr := tx.Rollback()
		if rerr != nil {
			err = errors.Wrapf(err, "rolling back transaction: %v", rerr)
		}
		return err
	}

	err = tx.Commit()
	if err != nil {
		return errors.Wrapf(err, "committing transaction")
	}
	return nil
}

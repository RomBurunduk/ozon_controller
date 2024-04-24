package repository

import (
	"errors"
)

var ErrObjectNotFound = errors.New("not found")

type PVZDbId int

type PvzDb struct {
	Id      int64  `db:"id"`
	Name    string `db:"name"`
	Address string `db:"address"`
	Contact string `db:"contact"`
}

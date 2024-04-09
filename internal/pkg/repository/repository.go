//go:generate mockgen -source ./repository.go -destination=./mocks/repository.go -package=mock_repository
package repository

import (
	"context"

	"pvz_controller/internal/model"
)

type PVZRepo interface {
	Add(ctx context.Context, point *model.Pickups) (int64, error)
	GetByID(ctx context.Context, id int64) (PvzDb, error)
	DeleteByID(ctx context.Context, id int64) error
	Update(ctx context.Context, point *model.Pickups, id int64) error
	ListAll(ctx context.Context) ([]PvzDb, error)
}

package postgresql

import (
	"context"
	"database/sql"
	"errors"

	"pvz_controller/internal/model"
	"pvz_controller/internal/pkg/db"
	"pvz_controller/internal/pkg/repository"
	"pvz_controller/internal/pkg/repository/transaction_manager"
)

type PVZRepo struct {
	db db.PGX
	tx transaction_manager.QueryEngineProvider
}

func NewArticles(database db.PGX) *PVZRepo {
	tx := transaction_manager.NewTransactionManager(database.GetPool(context.TODO()))
	return &PVZRepo{
		db: database,
		tx: tx,
	}
}

func (r *PVZRepo) Add(ctx context.Context, point *model.Pickups) (int64, error) {
	tx, err := r.db.BeginTx(ctx)
	if err != nil {
		return 0, err
	}
	defer tx.Rollback(ctx)

	var id int64
	err = r.db.ExecQueryRow(ctx, `INSERT INTO pickpoints(name,address,contact) VALUES ($1,$2,$3) RETURNING id;`,
		point.Name, point.Address, point.Contact).Scan(&id)

	if err = tx.Commit(ctx); err != nil {
		return 0, err
	}

	return id, err
}

func (r *PVZRepo) GetByID(ctx context.Context, id int64) (repository.PvzDb, error) {
	tx, err := r.db.BeginTx(ctx)
	if err != nil {
		return repository.PvzDb{}, err
	}
	defer tx.Rollback(ctx)

	var a repository.PvzDb
	err = r.db.Get(ctx, &a, "SELECT id, name,address,contact FROM pickpoints where id=$1", id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return repository.PvzDb{}, repository.ErrObjectNotFound
		}
		return repository.PvzDb{}, err
	}

	if err = tx.Commit(ctx); err != nil {
		return repository.PvzDb{}, err
	}

	return a, nil
}

func (r *PVZRepo) DeleteByID(ctx context.Context, id int64) error {
	tx, err := r.db.BeginTx(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	_, err = r.db.Exec(ctx, "DELETE FROM pickpoints WHERE id=$1", id)
	if err != nil {
		return err
	}

	if err = tx.Commit(ctx); err != nil {
		return err
	}
	return nil
}

func (r *PVZRepo) Update(ctx context.Context, point *model.Pickups, id int64) error {
	tx, err := r.db.BeginTx(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	_, err = r.db.Exec(ctx, `
        UPDATE pickpoints 
        SET name=$1, address=$2, contact=$3 
        WHERE id=$4
    `, point.Name, point.Address, point.Contact, id)

	if err = tx.Commit(ctx); err != nil {
		return err
	}
	return err
}

func (r *PVZRepo) ListAll(ctx context.Context) ([]repository.PvzDb, error) {
	tx, err := r.db.BeginTx(ctx)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback(ctx)
	var pickups []repository.PvzDb
	err = r.db.Select(ctx, &pickups, "SELECT id, name, address, contact FROM pickpoints")
	if err != nil {
		return nil, err
	}

	if err = tx.Commit(ctx); err != nil {
		return nil, err
	}
	return pickups, nil
}

package postgresql

import (
	"context"
	"fmt"
	"strings"
	"testing"

	"pvz_controller/internal/pkg/db"
)

type TDB struct {
	DB db.DBops
}

func NewFromEnv() *TDB {
	database, err := db.NewDb(context.Background())
	if err != nil {
		panic(err)
	}
	return &TDB{DB: database}
}

func (d *TDB) SetUp(t *testing.T, tableName ...string) {
	t.Helper()
	d.truncateTable(context.Background(), tableName...)
}

func (d *TDB) TearDown() {

}

func (d *TDB) truncate(ctx context.Context) {
	var tables []string
	err := d.DB.Select(ctx, &tables, "SELECT table_name FROM information_schema.tables WHERE table_schema='public' AND table_type='BASE TABLE' AND table_name != 'goose_db_version'")
	if err != nil {
		panic(err)
	}
	if len(tables) == 0 {
		panic(
			"run migration first")
	}
	q := fmt.Sprintf("TRUNCATE table %s RESTART IDENTITY", strings.Join(tables, ","))
	if _, err := d.DB.Exec(ctx, q); err != nil {
		panic(err)
	}
}

func (d *TDB) truncateTable(ctx context.Context, tableName ...string) {

	q := fmt.Sprintf("TRUNCATE table %s RESTART IDENTITY", strings.Join(tableName, ","))
	if _, err := d.DB.Exec(ctx, q); err != nil {
		panic(err)
	}
}

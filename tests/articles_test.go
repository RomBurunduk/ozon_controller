//go:build integration

package tests

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"pvz_controller/internal/pkg/repository"
	"pvz_controller/internal/pkg/repository/postgresql"
	"pvz_controller/tests/fixtures"
)

func TestCreateArticle(t *testing.T) {
	var (
		ctx = context.Background()
	)

	t.Run("smoke test", func(t *testing.T) {

		db.SetUp(t, "pickpoints")
		defer db.TearDown()
		// arrange
		repo := postgresql.NewPVZRepo(db.DB)
		//act
		resp, err := repo.Add(ctx, fixtures.Article().Valid().P())
		//assert
		require.NoError(t, err)
		assert.Equal(t, resp, int64(1))
	})
}

func TestGetArticle(t *testing.T) {
	var (
		ctx = context.Background()
	)

	db.SetUp(t, "pickpoints")
	defer db.TearDown()
	// arrange
	repo := postgresql.NewPVZRepo(db.DB)
	respAdd, err := repo.Add(ctx, fixtures.Article().Valid().P())
	require.NoError(t, err)

	//act
	resp, err := repo.GetByID(ctx, respAdd)
	//assert
	require.NoError(t, err)
	assert.Equal(t, resp, repository.PvzDb{
		Id:      respAdd,
		Name:    "some",
		Address: "some",
		Contact: "some",
	})
}

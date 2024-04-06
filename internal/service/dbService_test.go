package service

import (
	"context"
	"net/http"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"pvz_controller/internal/model"
	"pvz_controller/internal/pkg/repository"
)

func TestGetByID(t *testing.T) {
	t.Parallel()
	var (
		ctx = context.Background()
		id  = int64(1)
	)
	t.Run("smoke test", func(t *testing.T) {
		t.Parallel()
		// arrange
		s := setUp(t)
		defer s.tearDown()
		s.mockArticles.EXPECT().GetByID(gomock.Any(), id).Return(repository.PvzDb{
			Id:      1,
			Name:    "1",
			Address: "1",
			Contact: "1",
		}, nil)
		//act
		result, status := s.srv.get(ctx, id)
		//assert
		require.Equal(t, http.StatusOK, status)
		assert.Equal(t, "{\"Id\":1,\"Name\":\"1\",\"Address\":\"1\",\"Contact\":\"1\"}", string(result))
	})
}

func Test_DeleteById(t *testing.T) {
	t.Parallel()
	var (
		ctx = context.Background()
		id  = int64(1)
	)
	t.Run("smoke test", func(t *testing.T) {
		t.Parallel()
		//arrange
		s := setUp(t)
		defer s.tearDown()
		s.mockArticles.EXPECT().DeleteByID(gomock.Any(), id).Return(nil)
		//act
		err, status := s.srv.delete(ctx, id)
		//assert
		require.Equal(t, http.StatusOK, status)
		require.NoError(t, err)
	})
}

func Test_listAll(t *testing.T) {
	t.Parallel()
	var (
		ctx = context.Background()
	)
	t.Run("smoke test", func(t *testing.T) {
		t.Parallel()
		//arrange
		s := setUp(t)
		defer s.tearDown()
		s.mockArticles.EXPECT().ListAll(gomock.Any()).Return([]repository.PvzDb{{
			Id:      1,
			Name:    "1",
			Address: "1",
			Contact: "1",
		}, {
			Id:      2,
			Name:    "2",
			Address: "2",
			Contact: "2",
		}}, nil)
		//act
		data, status := s.srv.listAll(ctx)
		//assert
		require.Equal(t, "[{\"Id\":1,\"Name\":\"1\",\"Address\":\"1\",\"Contact\":\"1\"},{\"Id\":2,\"Name\":\"2\",\"Address\":\"2\",\"Contact\":\"2\"}]", string(data))
		require.Equal(t, http.StatusOK, status)
	})
}

func Test_update(t *testing.T) {
	t.Parallel()
	var (
		ctx    = context.Background()
		id     = int64(1)
		pickup = model.Pickups{
			Name:    "1",
			Address: "1",
			Contact: "1",
		}
	)
	t.Run("smoke test", func(t *testing.T) {
		t.Parallel()
		//arrange
		s := setUp(t)
		defer s.tearDown()
		s.mockArticles.EXPECT().Update(gomock.Any(), &pickup, id).Return(nil)
		//act
		err, status := s.srv.update(ctx, pickup, id)
		//assert
		require.NoError(t, err)
		require.Equal(t, http.StatusOK, status)
	})

}

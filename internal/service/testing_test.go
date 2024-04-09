package service

import (
	"testing"

	"github.com/golang/mock/gomock"

	mock_repository "pvz_controller/internal/pkg/repository/mocks"
	mock_pickupStorage "pvz_controller/internal/storage/mocks"
)

type pvzRepoFixtures struct {
	ctrl         *gomock.Controller
	srv          ServerService
	mockArticles *mock_repository.MockPVZRepo
}

func setUp(t *testing.T) pvzRepoFixtures {
	ctrl := gomock.NewController(t)
	mockPvz := mock_repository.NewMockPVZRepo(ctrl)
	srv := ServerService{Repo: mockPvz}
	return pvzRepoFixtures{
		ctrl:         ctrl,
		mockArticles: mockPvz,
		srv:          srv,
	}
}

func (a *pvzRepoFixtures) tearDown() {
	a.ctrl.Finish()
}

type pickupStorageFixtures struct {
	ctrl       *gomock.Controller
	srv        PickupService
	mockPickup *mock_pickupStorage.MockDBops
}

func setStorageUp(t *testing.T) pickupStorageFixtures {
	ctrl := gomock.NewController(t)
	mockStorage := mock_pickupStorage.NewMockDBops(ctrl)
	srv := PickupService{s: mockStorage}
	return pickupStorageFixtures{
		ctrl:       ctrl,
		srv:        srv,
		mockPickup: mockStorage,
	}
}

func (p *pickupStorageFixtures) tearDown() {
	p.ctrl.Finish()
}

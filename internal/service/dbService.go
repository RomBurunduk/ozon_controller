package service

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
	"github.com/jackc/pgx/v4/pgxpool"

	"pvz_controller/internal/model"
	"pvz_controller/internal/pkg/repository"
	"pvz_controller/internal/pkg/repository/in_memory"
	"pvz_controller/internal/pkg/repository/redis"
	"pvz_controller/internal/pkg/repository/transaction_manager"
)

type ServerInterface interface {
	Create(w http.ResponseWriter, req *http.Request)
	GetByID(w http.ResponseWriter, req *http.Request)
	DeleteById(w http.ResponseWriter, req *http.Request)
	Update(w http.ResponseWriter, req *http.Request)
	ListAll(w http.ResponseWriter, req *http.Request)
	DeleteList(w http.ResponseWriter, req *http.Request)
}

type Redis interface {
	Set(key string, value interface{}, expiration time.Duration) error
	Get(ctx context.Context, key string) (string, error)
	Del(ctx context.Context, key string) error
}

type InMemory interface {
	Set(id repository.PVZDbId, item repository.PvzDb, expiration time.Duration)
	Get(id repository.PVZDbId) (repository.PvzDb, error)
	Delete(id repository.PVZDbId)
}

type TransactionManager interface {
	RunSerializable(ctx context.Context, f func(ctxTX context.Context) error) error
}

const QueryParamKey = "key"

type addPvzRequest struct {
	Name    string `json:"name"`
	Address string `json:"address"`
	Contact string `json:"contact"`
}

type addPvzResponse struct {
	Id      int64  `json:"id"`
	Name    string `json:"name"`
	Address string `json:"address"`
	Contact string `json:"contact"`
}

func NewServerService(repo repository.PVZRepo, db *pgxpool.Pool) ServerInterface {
	c := in_memory.NewInMemoryCache()
	r := redis.NewRedis()
	tx := transaction_manager.NewTransactionManager(db)
	return &ServerService{
		Repo:      repo,
		cache:     c,
		redis:     r,
		txManager: tx}
}

type ServerService struct {
	Repo      repository.PVZRepo
	cache     InMemory
	redis     Redis
	txManager TransactionManager
}

func (s *ServerService) Create(w http.ResponseWriter, req *http.Request) {
	body, err := io.ReadAll(req.Body)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	var unm addPvzRequest
	if err = json.Unmarshal(body, &unm); err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	PvzRepo := &model.Pickups{
		Name:    unm.Name,
		Address: unm.Address,
		Contact: unm.Contact,
	}
	id, err := s.Repo.Add(req.Context(), PvzRepo)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	resp := addPvzResponse{
		Id:      id,
		Name:    PvzRepo.Name,
		Address: PvzRepo.Address,
		Contact: PvzRepo.Contact,
	}
	articleJson, _ := json.Marshal(resp)
	_, err = w.Write(articleJson)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func (s *ServerService) GetByID(w http.ResponseWriter, req *http.Request) {
	key, ok := mux.Vars(req)[QueryParamKey]
	if !ok {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	keyInt, err := strconv.ParseInt(key, 10, 64)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	hash := s.reqHash(req)
	get, err := s.redis.Get(req.Context(), hash)
	if err != nil {
		articleJson, status := s.get(req.Context(), keyInt)
		_, err = w.Write(articleJson)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.WriteHeader(status)
		s.redis.Set(hash, string(articleJson), time.Hour*3)
		return
	}
	_, err = w.Write([]byte(get))
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)

}

func (s *ServerService) get(ctx context.Context, keyInt int64) ([]byte, int) {
	get, err := s.cache.Get(repository.PVZDbId(keyInt))
	if err != nil {
		article, err := s.Repo.GetByID(ctx, keyInt)
		if err != nil {
			if errors.Is(err, repository.ErrObjectNotFound) {
				return nil, http.StatusNotFound
			}
			return nil, http.StatusInternalServerError
		}
		s.cache.Set(repository.PVZDbId(keyInt), article, 12*time.Hour)
		articleJson, _ := json.Marshal(article)
		return articleJson, http.StatusOK
	}
	articleJson, _ := json.Marshal(get)
	return articleJson, http.StatusOK
}

func (s *ServerService) DeleteById(w http.ResponseWriter, req *http.Request) {
	key, ok := mux.Vars(req)[QueryParamKey]
	if !ok {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	keyInt, err := strconv.ParseInt(key, 10, 64)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	s.redis.Del(req.Context(), s.reqHash(req))

	err, status := s.delete(req.Context(), keyInt)
	_, err = w.Write([]byte("success"))
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(status)
}

func (s *ServerService) delete(ctx context.Context, keyInt int64) (error, int) {
	err := s.Repo.DeleteByID(ctx, keyInt)
	if err != nil {
		if errors.Is(err, repository.ErrObjectNotFound) {
			return nil, http.StatusNotFound
		}
		return nil, http.StatusInternalServerError
	}
	s.cache.Delete(repository.PVZDbId(keyInt))
	return err, http.StatusOK
}

func (s *ServerService) Update(w http.ResponseWriter, req *http.Request) {
	key, ok := mux.Vars(req)[QueryParamKey]
	if !ok {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	keyInt, err := strconv.ParseInt(key, 10, 64)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	s.redis.Del(req.Context(), s.reqHash(req))
	body, err := io.ReadAll(req.Body)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	var unm addPvzRequest
	if err := json.Unmarshal(body, &unm); err != nil {
		return
	}
	PvzRepo := &model.Pickups{
		Name:    unm.Name,
		Address: unm.Address,
		Contact: unm.Contact,
	}
	err, status := s.update(req.Context(), *PvzRepo, keyInt)
	_, err = w.Write([]byte("success"))
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(status)
}

func (s *ServerService) update(ctx context.Context, PvzRepo model.Pickups, keyInt int64) (error, int) {
	err := s.Repo.Update(ctx, &PvzRepo, keyInt)
	if err != nil {
		return nil, http.StatusInternalServerError
	}
	s.cache.Set(repository.PVZDbId(keyInt), repository.PvzDb{
		Id:      keyInt,
		Name:    PvzRepo.Name,
		Address: PvzRepo.Address,
		Contact: PvzRepo.Contact,
	}, 12*time.Hour)
	return err, http.StatusOK
}

func (s *ServerService) ListAll(w http.ResponseWriter, req *http.Request) {
	data, status := s.listAll(req.Context())
	_, err := w.Write(data)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(status)
}

func (s *ServerService) listAll(ctx context.Context) ([]byte, int) {
	all, err := s.Repo.ListAll(ctx)
	if err != nil {
		return nil, http.StatusInternalServerError
	}
	data, err := json.Marshal(all)
	if err != nil {
		return nil, http.StatusInternalServerError
	}
	return data, http.StatusOK
}

// DeleteList - удаляет массив ПВЗ
func (s *ServerService) DeleteList(w http.ResponseWriter, req *http.Request) {
	body, err := io.ReadAll(req.Body)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	var nums []int
	if err = json.Unmarshal(body, &nums); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	err = s.txManager.RunSerializable(req.Context(), func(ctx context.Context) error {
		for _, num := range nums {
			err = s.Repo.DeleteByID(req.Context(), int64(num))
			if err != nil {
				return err
			}
		}
		return nil
	})
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Println(err)
		return
	}
}

func (s *ServerService) reqHash(req *http.Request) string {
	var hash string
	hash = req.Header.Get("X-Hash")
	return hash
}

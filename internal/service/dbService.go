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

	"pvz_controller/internal/model"
	"pvz_controller/internal/pkg/repository"
	"pvz_controller/internal/pkg/repository/in_memory"
	"pvz_controller/internal/pkg/repository/redis"
)

type ServerInterface interface {
	Create(w http.ResponseWriter, req *http.Request)
	GetByID(w http.ResponseWriter, req *http.Request)
	DeleteById(w http.ResponseWriter, req *http.Request)
	Update(w http.ResponseWriter, req *http.Request)
	ListAll(w http.ResponseWriter, req *http.Request)
}

type Redis interface {
	Set(key string, value interface{}, expiration time.Duration) error
	Get(ctx context.Context, key string) (string, error)
	Del(ctx context.Context, key string) error
}

type InMemory interface {
	Set(id repository.PVZDbId, item repository.PvzDb)
	Get(id repository.PVZDbId) (repository.PvzDb, error)
	Delete(id repository.PVZDbId)
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

func NewServerService(repo repository.PVZRepo) ServerInterface {
	c := in_memory.NewInMemoryCache()
	r := redis.NewRedis()
	return &ServerService{
		Repo:  repo,
		cache: c,
		redis: r}
}

type ServerService struct {
	Repo  repository.PVZRepo
	cache InMemory
	redis Redis
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
		s.cache.Set(repository.PVZDbId(keyInt), article)
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
	})
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

func (s *ServerService) reqHash(req *http.Request) string {
	var hash string
	hash = req.Header.Get("X-Hash")
	return hash
}

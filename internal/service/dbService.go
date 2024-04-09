package service

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"

	"pvz_controller/internal/model"
	"pvz_controller/internal/pkg/repository"
)

type ServerInterface interface {
	Create(w http.ResponseWriter, req *http.Request)
	GetByID(w http.ResponseWriter, req *http.Request)
	DeleteById(w http.ResponseWriter, req *http.Request)
	Update(w http.ResponseWriter, req *http.Request)
	ListAll(w http.ResponseWriter, req *http.Request)
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
	return &ServerService{Repo: repo}
}

type ServerService struct {
	Repo repository.PVZRepo
}

func (s *ServerService) Create(w http.ResponseWriter, req *http.Request) {
	body, err := io.ReadAll(req.Body)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	var unm addPvzRequest
	if err = json.Unmarshal(body, &unm); err != nil {
		fmt.Println(err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	PvzRepo := &model.Pickups{
		Name:    unm.Name,
		Address: unm.Address,
		Contact: unm.Contact,
	}
	id, err := s.Repo.Add(req.Context(), PvzRepo)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
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
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}

func (s *ServerService) GetByID(w http.ResponseWriter, req *http.Request) {
	key, ok := mux.Vars(req)[QueryParamKey]
	if !ok {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}
	keyInt, err := strconv.ParseInt(key, 10, 64)
	if err != nil {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	articleJson, status := s.get(req.Context(), keyInt)
	_, err = w.Write(articleJson)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(status)
}

func (s *ServerService) get(ctx context.Context, keyInt int64) ([]byte, int) {
	article, err := s.Repo.GetByID(ctx, keyInt)
	if err != nil {
		if errors.Is(err, repository.ErrObjectNotFound) {
			return nil, http.StatusNotFound
		}
		return nil, http.StatusInternalServerError
	}
	articleJson, _ := json.Marshal(article)

	return articleJson, http.StatusOK
}

func (s *ServerService) DeleteById(w http.ResponseWriter, req *http.Request) {
	key, ok := mux.Vars(req)[QueryParamKey]
	if !ok {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}
	keyInt, err := strconv.ParseInt(key, 10, 64)
	if err != nil {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}
	err, status := s.delete(req.Context(), keyInt)
	_, err = w.Write([]byte("success"))
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
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
	return err, http.StatusOK
}

func (s *ServerService) Update(w http.ResponseWriter, req *http.Request) {
	key, ok := mux.Vars(req)[QueryParamKey]
	if !ok {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}
	keyInt, err := strconv.ParseInt(key, 10, 64)
	if err != nil {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	body, err := io.ReadAll(req.Body)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
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
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(status)
}

func (s *ServerService) update(ctx context.Context, PvzRepo model.Pickups, keyInt int64) (error, int) {
	err := s.Repo.Update(ctx, &PvzRepo, keyInt)
	if err != nil {
		return nil, http.StatusInternalServerError
	}
	return err, http.StatusOK
}

func (s *ServerService) ListAll(w http.ResponseWriter, req *http.Request) {
	data, status := s.listAll(req.Context())
	_, err := w.Write(data)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
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

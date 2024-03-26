package service

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"

	"pvz_controller/internal/model"
	"pvz_controller/internal/pkg/repository"
	"pvz_controller/internal/pkg/repository/postgresql"
)

const QueryParamKey = "key"

type ServerService interface {
	Create(w http.ResponseWriter, req *http.Request)
	GetByID(w http.ResponseWriter, req *http.Request)
	DeleteById(w http.ResponseWriter, req *http.Request)
	Update(w http.ResponseWriter, req *http.Request)
	ListAll(w http.ResponseWriter, req *http.Request)
}

type serverService struct {
	repo postgresql.PVZRepo
}

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

func NewServerService(repo *postgresql.PVZRepo) ServerService {
	return &serverService{repo: *repo}
}

func (s *serverService) Create(w http.ResponseWriter, req *http.Request) {
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
	id, err := s.repo.Add(req.Context(), PvzRepo)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	resp := &addPvzResponse{
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

func (s *serverService) GetByID(w http.ResponseWriter, req *http.Request) {
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

	pvz, err := s.repo.GetByID(req.Context(), keyInt)
	if err != nil {
		if errors.Is(err, repository.ErrObjectNotFound) {
			http.Error(w, "Not Found", http.StatusNotFound)
			return
		}
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	articleJson, _ := json.Marshal(pvz)
	_, err = w.Write(articleJson)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}

func (s *serverService) DeleteById(w http.ResponseWriter, req *http.Request) {
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
	err = s.repo.DeleteByID(req.Context(), keyInt)
	if err != nil {
		if errors.Is(err, repository.ErrObjectNotFound) {
			http.Error(w, "Not Found", http.StatusNotFound)
			return
		}
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	_, err = w.Write([]byte("success"))
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}

func (s *serverService) Update(w http.ResponseWriter, req *http.Request) {
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
	err = s.repo.Update(req.Context(), PvzRepo, keyInt)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	_, err = w.Write([]byte("success"))
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}

func (s *serverService) ListAll(w http.ResponseWriter, req *http.Request) {
	all, err := s.repo.ListAll(req.Context())
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	data, err := json.Marshal(all)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	_, err = w.Write(data)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}

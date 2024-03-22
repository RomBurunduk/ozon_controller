package storage

import (
	"bufio"
	"encoding/json"
	"io"
	"io/fs"
	"os"
	"sync"

	"pvz_controller/internal/model"
)

type PVZStorage struct {
	file  *os.File
	mutex sync.RWMutex
}

const pickupStorageName = "pickupStorage"

func NewPickupStorage() (PVZStorage, error) {
	file, err := os.OpenFile(pickupStorageName, os.O_CREATE, fs.ModePerm)
	if err != nil {
		return PVZStorage{}, err
	}
	return PVZStorage{file: file, mutex: sync.RWMutex{}}, err
}

func (s *PVZStorage) WritePVZ(pvz model.Pickups) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	existingData, err := s.readDataFromFile()
	if err != nil {
		return err
	}

	existingData = append(existingData, pvz)

	return s.writeBytes(existingData)
}

func (s *PVZStorage) ReadPVZ() ([]model.Pickups, error) {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	return s.readDataFromFile()
}

func (s *PVZStorage) readDataFromFile() ([]model.Pickups, error) {
	reader := bufio.NewReader(s.file)
	RawBytes, err := io.ReadAll(reader)
	if err != nil {
		return nil, err
	}
	_, err = s.file.Seek(0, 0)
	if err != nil {
		return nil, err
	}

	var orders []model.Pickups
	if len(RawBytes) == 0 {
		return orders, nil
	}
	err = json.Unmarshal(RawBytes, &orders)
	if err != nil {
		return nil, err
	}
	return orders, nil
}

func (s *PVZStorage) writeBytes(all []model.Pickups) error {
	bytes, err := json.Marshal(all)
	if err != nil {
		return err
	}
	err = os.WriteFile(pickupStorageName, bytes, fs.ModePerm)
	if err != nil {
		return err
	}
	return nil
}

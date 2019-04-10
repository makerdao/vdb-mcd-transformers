package test_helpers

import (
	"github.com/vulcanize/vulcanizedb/pkg/datastore/postgres"

	"github.com/vulcanize/mcd_transformers/transformers/storage"
)

type MockMakerStorageRepository struct {
	DaiKeys             []string
	GemKeys             []storage.Urn
	GetDaiKeysCalled    bool
	GetDaiKeysError     error
	GetGemKeysCalled    bool
	GetGemKeysError     error
	GetIlksCalled       bool
	GetIlksError        error
	GetMaxFlipCalled    bool
	GetMaxFlipError     error
	GetVatSinKeysCalled bool
	GetVatSinKeysError  error
	GetVowSinKeysCalled bool
	GetVowSinKeysError  error
	GetUrnsCalled       bool
	GetUrnsError        error
	Ilks                []string
	MaxFlip             int64
	SinKeys             []string
	Urns                []storage.Urn
}

func (repository *MockMakerStorageRepository) GetDaiKeys() ([]string, error) {
	repository.GetDaiKeysCalled = true
	return repository.DaiKeys, repository.GetDaiKeysError
}

func (repository *MockMakerStorageRepository) GetGemKeys() ([]storage.Urn, error) {
	repository.GetGemKeysCalled = true
	return repository.GemKeys, repository.GetGemKeysError
}

func (repository *MockMakerStorageRepository) GetIlks() ([]string, error) {
	repository.GetIlksCalled = true
	return repository.Ilks, repository.GetIlksError
}

func (repository *MockMakerStorageRepository) GetMaxFlip() (int64, error) {
	repository.GetMaxFlipCalled = true
	return repository.MaxFlip, repository.GetMaxFlipError
}

func (repository *MockMakerStorageRepository) GetVatSinKeys() ([]string, error) {
	repository.GetVatSinKeysCalled = true
	return repository.SinKeys, repository.GetVatSinKeysError
}

func (repository *MockMakerStorageRepository) GetVowSinKeys() ([]string, error) {
	repository.GetVowSinKeysCalled = true
	return repository.SinKeys, repository.GetVowSinKeysError
}

func (repository *MockMakerStorageRepository) GetUrns() ([]storage.Urn, error) {
	repository.GetUrnsCalled = true
	return repository.Urns, repository.GetUrnsError
}

func (repository *MockMakerStorageRepository) SetDB(db *postgres.DB) {}

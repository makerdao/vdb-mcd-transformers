package mocks

import (
	"github.com/ethereum/go-ethereum/common"

	"github.com/vulcanize/vulcanizedb/libraries/shared/storage/utils"
	"github.com/vulcanize/vulcanizedb/pkg/datastore/postgres"
)

type MockMappings struct {
	Metadata     utils.StorageValueMetadata
	LookupCalled bool
	LookupErr    error
}

func (mappings *MockMappings) Lookup(key common.Hash) (utils.StorageValueMetadata, error) {
	mappings.LookupCalled = true
	return mappings.Metadata, mappings.LookupErr
}

func (*MockMappings) SetDB(db *postgres.DB) {
	panic("implement me")
}

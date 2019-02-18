package mocks

import (
	"github.com/ethereum/go-ethereum/core/types"

	"github.com/vulcanize/vulcanizedb/libraries/shared/constants"
	"github.com/vulcanize/vulcanizedb/libraries/shared/transformer"
	"github.com/vulcanize/vulcanizedb/pkg/core"
	"github.com/vulcanize/vulcanizedb/pkg/datastore/postgres"
)

type MockTransformer struct {
	ExecuteWasCalled bool
	ExecuteError     error
	PassedLogs       []types.Log
	PassedHeader     core.Header
	config           transformer.TransformerConfig
}

func (mh *MockTransformer) Execute(logs []types.Log, header core.Header, recheckHeaders constants.TransformerExecution) error {
	if mh.ExecuteError != nil {
		return mh.ExecuteError
	}
	mh.ExecuteWasCalled = true
	mh.PassedLogs = logs
	mh.PassedHeader = header
	return nil
}

func (mh *MockTransformer) GetConfig() transformer.TransformerConfig {
	return mh.config
}

func (mh *MockTransformer) SetTransformerConfig(config transformer.TransformerConfig) {
	mh.config = config
}

func (mh *MockTransformer) FakeTransformerInitializer(db *postgres.DB) transformer.EventTransformer {
	return mh
}

var FakeTransformerConfig = transformer.TransformerConfig{
	TransformerName:   "FakeTransformer",
	ContractAddresses: []string{"FakeAddress"},
	Topic:             "FakeTopic",
}

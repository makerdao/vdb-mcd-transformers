package mocks

import (
	"github.com/makerdao/vdb-mcd-transformers/backfill/shared"
)

type MockDartDinkRetriever struct {
	RetrieveDartDinkDiffsError error
	PassedDartDinks            []shared.DartDink
}

func (mock *MockDartDinkRetriever) RetrieveDartDinkDiffs(dartDink shared.DartDink) error {
	mock.PassedDartDinks = append(mock.PassedDartDinks, dartDink)
	return mock.RetrieveDartDinkDiffsError
}

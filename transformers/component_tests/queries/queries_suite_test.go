package queries

import (
	"testing"

	"github.com/makerdao/vdb-mcd-transformers/test_config"
	"github.com/makerdao/vulcanizedb/pkg/datastore/postgres"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var db *postgres.DB

func TestQueries(t *testing.T) {
	RegisterFailHandler(Fail)
	db = test_config.NewTestDB(test_config.NewTestNode())
	RunSpecs(t, "Queries Component Test Suite")
}

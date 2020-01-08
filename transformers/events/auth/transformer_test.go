package auth_test

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/makerdao/vdb-mcd-transformers/test_config"
	"github.com/makerdao/vdb-mcd-transformers/transformers/events/auth"
	"github.com/makerdao/vdb-mcd-transformers/transformers/shared"
	"github.com/makerdao/vdb-mcd-transformers/transformers/shared/constants"
	"github.com/makerdao/vdb-mcd-transformers/transformers/test_data"
	"github.com/makerdao/vulcanizedb/libraries/shared/factories/event"
	"github.com/makerdao/vulcanizedb/pkg/core"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Deny Transformer", func() {
	var db = test_config.NewTestDB(test_config.NewTestNode())

	BeforeEach(func() {
		test_config.CleanTestDB(db)
	})

	It("converts normal rely logs to models", func() {
		transformer := auth.Transformer{TableName: constants.RelyTable}
		models, err := transformer.ToModels(constants.CatABI(), []core.EventLog{test_data.RelyEventLog}, db)
		Expect(err).NotTo(HaveOccurred())

		var contractAddressID int64
		contractAddressErr := db.Get(&contractAddressID, `SELECT id FROM addresses WHERE address = $1`,
			test_data.RelyEventLog.Log.Address.String())
		Expect(contractAddressErr).NotTo(HaveOccurred())

		var usrAddressID int64
		usrAddressErr := db.Get(&usrAddressID, `SELECT id FROM addresses WHERE address = $1`,
			common.HexToAddress(test_data.RelyEventLog.Log.Topics[2].Hex()).Hex())
		Expect(usrAddressErr).NotTo(HaveOccurred())

		expectedModel := test_data.RelyModel()
		expectedModel.ColumnValues[event.AddressFK] = contractAddressID
		expectedModel.ColumnValues[constants.UsrColumn] = usrAddressID

		Expect(models).To(Equal([]event.InsertionModel{expectedModel}))
	})

	It("converts Vat rely logs to models", func() {
		transformer := auth.Transformer{TableName: constants.RelyTable, LogNoteArgumentOffset: -1}
		models, err := transformer.ToModels(constants.VatABI(), []core.EventLog{test_data.VatRelyEventLog}, db)
		Expect(err).NotTo(HaveOccurred())

		var contractAddressID int64
		contractAddressErr := db.Get(&contractAddressID, `SELECT id FROM addresses WHERE address = $1`,
			test_data.VatRelyEventLog.Log.Address.String())
		Expect(contractAddressErr).NotTo(HaveOccurred())

		var usrAddressID int64
		usrAddressErr := db.Get(&usrAddressID, `SELECT id FROM addresses WHERE address = $1`,
			common.HexToAddress(test_data.VatRelyEventLog.Log.Topics[1].Hex()).Hex())
		Expect(usrAddressErr).NotTo(HaveOccurred())

		expectedModel := test_data.VatRelyModel()
		expectedModel.ColumnValues[event.AddressFK] = contractAddressID
		expectedModel.ColumnValues[constants.UsrColumn] = usrAddressID

		Expect(models).To(Equal([]event.InsertionModel{expectedModel}))
	})

	It("converts normal deny logs to models", func() {
		transformer := auth.Transformer{TableName: constants.DenyTable}
		models, err := transformer.ToModels(constants.CatABI(), []core.EventLog{test_data.DenyEventLog}, db)
		Expect(err).NotTo(HaveOccurred())

		var contractAddressID int64
		contractAddressErr := db.Get(&contractAddressID, `SELECT id FROM addresses WHERE address = $1`,
			test_data.DenyEventLog.Log.Address.String())
		Expect(contractAddressErr).NotTo(HaveOccurred())

		var usrAddressID int64
		usrAddressErr := db.Get(&usrAddressID, `SELECT id FROM addresses WHERE address = $1`,
			common.HexToAddress(test_data.DenyEventLog.Log.Topics[2].Hex()).Hex())
		Expect(usrAddressErr).NotTo(HaveOccurred())

		expectedModel := test_data.DenyModel()
		expectedModel.ColumnValues[event.AddressFK] = contractAddressID
		expectedModel.ColumnValues[constants.UsrColumn] = usrAddressID

		Expect(models).To(Equal([]event.InsertionModel{expectedModel}))
	})

	It("converts Vat deny logs to models", func() {
		transformer := auth.Transformer{TableName: constants.DenyTable, LogNoteArgumentOffset: -1}
		models, err := transformer.ToModels(constants.VatABI(), []core.EventLog{test_data.VatDenyEventLog}, db)
		Expect(err).NotTo(HaveOccurred())

		var contractAddressID int64
		contractAddressErr := db.Get(&contractAddressID, `SELECT id FROM addresses WHERE address = $1`,
			test_data.VatDenyEventLog.Log.Address.String())
		Expect(contractAddressErr).NotTo(HaveOccurred())

		var usrAddressID int64
		usrAddressErr := db.Get(&usrAddressID, `SELECT id FROM addresses WHERE address = $1`,
			common.HexToAddress(test_data.VatDenyEventLog.Log.Topics[1].Hex()).Hex())
		Expect(usrAddressErr).NotTo(HaveOccurred())

		expectedModel := test_data.VatDenyModel()
		expectedModel.ColumnValues[event.AddressFK] = contractAddressID
		expectedModel.ColumnValues[constants.UsrColumn] = usrAddressID

		Expect(models).To(Equal([]event.InsertionModel{expectedModel}))
	})

	It("returns an error if the expected amount of topics aren't in the log", func() {
		transformer := auth.Transformer{}
		invalidLog := test_data.DenyEventLog
		invalidLog.Log.Topics = []common.Hash{}

		_, err := transformer.ToModels(constants.CatABI(), []core.EventLog{invalidLog}, db)

		Expect(err).To(HaveOccurred())
		Expect(err).To(MatchError(shared.ErrLogMissingTopics(3, 0)))
	})
})

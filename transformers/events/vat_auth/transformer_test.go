package vat_auth_test

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/makerdao/vdb-mcd-transformers/test_config"
	"github.com/makerdao/vdb-mcd-transformers/transformers/events/vat_auth"
	"github.com/makerdao/vdb-mcd-transformers/transformers/shared"
	"github.com/makerdao/vdb-mcd-transformers/transformers/shared/constants"
	"github.com/makerdao/vdb-mcd-transformers/transformers/test_data"
	"github.com/makerdao/vulcanizedb/libraries/shared/factories/event"
	"github.com/makerdao/vulcanizedb/pkg/core"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Vat Auth Transformer", func() {
	var db = test_config.NewTestDB(test_config.NewTestNode())

	BeforeEach(func() {
		test_config.CleanTestDB(db)
	})

	It("converts Vat rely logs to models", func() {
		converter := vat_auth.Transformer{TableName: constants.VatRelyTable}
		models, err := converter.ToModels("", []core.EventLog{test_data.VatRelyEventLog}, db)
		Expect(err).NotTo(HaveOccurred())

		usrAddressID, usrAddressErr := shared.GetOrCreateAddress(test_data.VatRelyEventLog.Log.Topics[1].Hex(), db)
		Expect(usrAddressErr).NotTo(HaveOccurred())

		expectedModel := test_data.VatRelyModel()
		expectedModel.ColumnValues[constants.UsrColumn] = usrAddressID

		Expect(models).To(Equal([]event.InsertionModel{expectedModel}))
	})

	It("converts Vat deny logs to models", func() {
		converter := vat_auth.Transformer{TableName: constants.VatDenyTable}
		models, err := converter.ToModels("", []core.EventLog{test_data.VatDenyEventLog}, db)
		Expect(err).NotTo(HaveOccurred())

		usrAddressID, usrAddressErr := shared.GetOrCreateAddress(test_data.VatDenyEventLog.Log.Topics[1].Hex(), db)
		Expect(usrAddressErr).NotTo(HaveOccurred())

		expectedModel := test_data.VatDenyModel()
		expectedModel.ColumnValues[constants.UsrColumn] = usrAddressID

		Expect(models).To(Equal([]event.InsertionModel{expectedModel}))
	})

	It("converts Vat hope logs to models", func() {
		converter := vat_auth.Transformer{TableName: constants.VatHopeTable}
		models, err := converter.ToModels("", []core.EventLog{test_data.VatHopeEventLog}, db)
		Expect(err).NotTo(HaveOccurred())

		usrAddressID, usrAddressErr := shared.GetOrCreateAddress(test_data.VatHopeEventLog.Log.Topics[1].Hex(), db)
		Expect(usrAddressErr).NotTo(HaveOccurred())

		expectedModel := test_data.VatHopeModel()
		expectedModel.ColumnValues[constants.UsrColumn] = usrAddressID

		Expect(models).To(Equal([]event.InsertionModel{expectedModel}))
	})

	It("converts Vat nope logs to models", func() {
		converter := vat_auth.Transformer{TableName: constants.VatNopeTable}
		models, err := converter.ToModels("", []core.EventLog{test_data.VatNopeEventLog}, db)
		Expect(err).NotTo(HaveOccurred())

		usrAddressID, usrAddressErr := shared.GetOrCreateAddress(test_data.VatNopeEventLog.Log.Topics[1].Hex(), db)
		Expect(usrAddressErr).NotTo(HaveOccurred())

		expectedModel := test_data.VatNopeModel()
		expectedModel.ColumnValues[constants.UsrColumn] = usrAddressID

		Expect(models).To(Equal([]event.InsertionModel{expectedModel}))
	})

	It("returns an error if the expected amount of topics aren't in the log", func() {
		converter := vat_auth.Transformer{}
		invalidLog := test_data.VatDenyEventLog
		invalidLog.Log.Topics = []common.Hash{}

		_, err := converter.ToModels(constants.Cat100ABI(), []core.EventLog{invalidLog}, db)

		Expect(err).To(HaveOccurred())
		Expect(err).To(MatchError(shared.ErrLogMissingTopics(2, 0)))
	})
})

package vat_grab_test

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/vulcanize/mcd_transformers/transformers/events/vat_grab"
	"github.com/vulcanize/mcd_transformers/transformers/test_data"
)

var _ = Describe("Vat grab converter", func() {
	var converter vat_grab.VatGrabConverter

	BeforeEach(func() {
		converter = vat_grab.VatGrabConverter{}
	})

	It("returns err if log is missing topics", func() {
		badLog := types.Log{
			Data: []byte{1, 1, 1, 1, 1},
		}

		_, err := converter.ToModels([]types.Log{badLog})

		Expect(err).To(HaveOccurred())
	})

	It("returns err if log is missing data", func() {
		badLog := types.Log{
			Topics: []common.Hash{{}, {}, {}, {}},
		}

		_, err := converter.ToModels([]types.Log{badLog})

		Expect(err).To(HaveOccurred())
	})

	It("converts a log with positive dink to a model", func() {
		models, err := converter.ToModels([]types.Log{test_data.EthVatGrabLogWithPositiveDink})

		Expect(err).NotTo(HaveOccurred())
		Expect(len(models)).To(Equal(1))
		Expect(models[0]).To(Equal(test_data.VatGrabModelWithPositiveDink))
	})

	It("converts a log with negative dink to a model", func() {
		models, err := converter.ToModels([]types.Log{test_data.EthVatGrabLogWithNegativeDink})

		Expect(err).NotTo(HaveOccurred())
		Expect(len(models)).To(Equal(1))
		Expect(models[0]).To(Equal(test_data.VatGrabModelWithNegativeDink))
	})
})

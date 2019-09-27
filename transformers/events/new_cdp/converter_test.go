// VulcanizeDB
// Copyright © 2018 Vulcanize

// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU Affero General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.

// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU Affero General Public License for more details.

// You should have received a copy of the GNU Affero General Public License
// along with this program.  If not, see <http://www.gnu.org/licenses/>.

package new_cdp_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/vulcanize/vulcanizedb/pkg/core"

	"github.com/vulcanize/mcd_transformers/transformers/events/new_cdp"
	"github.com/vulcanize/mcd_transformers/transformers/shared/constants"
	"github.com/vulcanize/mcd_transformers/transformers/test_data"
)

var _ = Describe("NewCdp Converter", func() {
	var converter = new_cdp.NewCdpConverter{}

	It("converts a log to a Model", func() {
		models, err := converter.ToModels(constants.CdpManagerABI(), []core.HeaderSyncLog{test_data.NewCdpHeaderSyncLog})

		Expect(err).NotTo(HaveOccurred())
		Expect(len(models)).To(Equal(1))
		Expect(models[0]).To(Equal(test_data.NewCdpModel()))
	})

	It("returns an error if converting log to entity fails", func() {
		_, err := converter.ToModels("error abi", []core.HeaderSyncLog{test_data.NewCdpHeaderSyncLog})
		Expect(err).To(HaveOccurred())
	})
})

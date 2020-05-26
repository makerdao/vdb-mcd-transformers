package backfill_test

import (
	"github.com/makerdao/vdb-mcd-transformers/backfill"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("BackFill CLI Argument Validator", func() {
	Describe("ValidateArgs", func() {
		It("returns error if passed 0 events", func() {
			err := backfill.ValidateArgs([]string{})

			Expect(err).To(HaveOccurred())
			Expect(err).To(MatchError(backfill.ErrTooFewEvents))
		})

		It("returns error if passed more events than maximum", func() {
			var eventsToBackFill []string
			for i := 0; i < backfill.MaxEvents; i++ {
				eventsToBackFill = append(eventsToBackFill, "s")
			}
			eventsToBackFill = append(eventsToBackFill, "s")

			err := backfill.ValidateArgs(eventsToBackFill)

			Expect(err).To(HaveOccurred())
			Expect(err).To(MatchError(backfill.ErrTooManyEvents))
		})

		It("returns error if passed unsupported event", func() {
			err := backfill.ValidateArgs([]string{"unsupported"})

			Expect(err).To(HaveOccurred())
			Expect(err).To(MatchError(backfill.ErrUnsupportedEvent))
		})

		It("returns nil for error if passed supported event", func() {
			err := backfill.ValidateArgs([]string{backfill.ForkEvent})
			Expect(err).NotTo(HaveOccurred())
		})

		It("returns error if any of two events are unsupported", func() {
			eventsToBackFill := []string{backfill.ForkEvent, "unsupported"}

			err := backfill.ValidateArgs(eventsToBackFill)

			Expect(err).To(HaveOccurred())
			Expect(err).To(MatchError(backfill.ErrUnsupportedEvent))
		})

		It("returns error if two events are duplicated", func() {
			eventsToBackFill := []string{backfill.ForkEvent, backfill.ForkEvent}

			err := backfill.ValidateArgs(eventsToBackFill)

			Expect(err).To(HaveOccurred())
			Expect(err).To(MatchError(backfill.ErrDuplicateEvent))
		})

		It("returns nil for error if passed two distinct and supported events", func() {
			eventsToBackFill := []string{backfill.ForkEvent, backfill.FrobEvent}

			err := backfill.ValidateArgs(eventsToBackFill)

			Expect(err).NotTo(HaveOccurred())
		})

		It("returns error if any of three events are unsupported", func() {
			eventsToBackFill := []string{backfill.ForkEvent, backfill.FrobEvent, "unsupported"}

			err := backfill.ValidateArgs(eventsToBackFill)

			Expect(err).To(HaveOccurred())
			Expect(err).To(MatchError(backfill.ErrUnsupportedEvent))
		})

		It("returns error if subset of three events are duplicated", func() {
			eventsToBackFill := []string{backfill.ForkEvent, backfill.FrobEvent, backfill.FrobEvent}

			err := backfill.ValidateArgs(eventsToBackFill)

			Expect(err).To(HaveOccurred())
			Expect(err).To(MatchError(backfill.ErrDuplicateEvent))
		})

		It("returns nil for error if passed three distinct and supported events", func() {
			eventsToBackFill := []string{backfill.ForkEvent, backfill.FrobEvent, backfill.GrabEvent}

			err := backfill.ValidateArgs(eventsToBackFill)

			Expect(err).NotTo(HaveOccurred())
		})
	})
})

package prompts_test

import (
	"github.com/AlecAivazis/survey/v2"
	"github.com/makerdao/vdb-mcd-transformers/generators/new_collateral/prompts"
	"github.com/makerdao/vulcanizedb/pkg/fakes"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("PromptsWrapper", func() {
	var (
		prompter                prompts.Prompter
		SurveyQuestionsPassedIn [][]*survey.Question
		AskErr                  error
	)

	BeforeEach(func() {
		AskErr = nil
		SurveyQuestionsPassedIn = [][]*survey.Question{}
		var (
			MockAskFunction = func(questions []*survey.Question, response interface{}, opts ...survey.AskOpt) error {
				SurveyQuestionsPassedIn = append(SurveyQuestionsPassedIn, questions)
				return AskErr
			}
		)

		prompter = prompts.Prompter{
			PromptAsker: MockAskFunction,
		}
	})

	Context("GetCollateralDetails", func() {
		It("gets collateral details", func() {
			_, err := prompter.GetCollateralDetails()
			Expect(err).NotTo(HaveOccurred())
			Expect(SurveyQuestionsPassedIn).To(ContainElement(prompts.CollateralQuestions))
		})

		It("returns an error if getting collateral details fails", func() {
			AskErr = fakes.FakeError
			_, err := prompter.GetCollateralDetails()
			Expect(err).To(HaveOccurred())
			Expect(err).To(MatchError(fakes.FakeError))
		})
	})

	Context("GetContractDetails", func() {
		It("gets flip contract details", func() {
			contracts, err := prompter.GetContractDetails()
			Expect(err).NotTo(HaveOccurred())
			Expect(SurveyQuestionsPassedIn).To(ContainElement(prompts.FlipContractQuestions))
			Expect(contracts).To(HaveKey("flip"))
		})

		It("asks if median contract is required", func() {
			_, err := prompter.GetContractDetails()
			Expect(err).NotTo(HaveOccurred())
			Expect(SurveyQuestionsPassedIn).To(ContainElement(prompts.MedianContractRequiredQuestion))
		})

		It("gets median contract details if it is required", func() {
			prompter.MedianRequired = true
			contracts, err := prompter.GetContractDetails()
			Expect(err).NotTo(HaveOccurred())
			Expect(SurveyQuestionsPassedIn).To(ContainElement(prompts.MedianContractQuestions))
			Expect(contracts).To(HaveKey("median"))
		})

		It("doesn't ask for median contract details if it's not required", func() {
			prompter.MedianRequired = false
			contracts, err := prompter.GetContractDetails()
			Expect(err).NotTo(HaveOccurred())
			Expect(contracts).NotTo(HaveKey("median"))
		})

		It("asks if osm contract is required", func() {
			_, err := prompter.GetContractDetails()
			Expect(err).NotTo(HaveOccurred())
			Expect(SurveyQuestionsPassedIn).To(ContainElement(prompts.OsmContractRequiredQuestion))
		})

		It("gets osm contract details if it is required", func() {
			prompter.OsmRequired = true
			contracts, err := prompter.GetContractDetails()
			Expect(err).NotTo(HaveOccurred())
			Expect(SurveyQuestionsPassedIn).To(ContainElement(prompts.OsmContractQuestions))
			Expect(contracts).To(HaveKey("osm"))
		})

		It("doesn't ask for osm contract details if it's not required", func() {
			prompter.OsmRequired = false
			contracts, err := prompter.GetContractDetails()
			Expect(err).NotTo(HaveOccurred())
			Expect(contracts).NotTo(HaveKey("osm"))
		})

		It("returns an error if getting contract details fails", func() {
			AskErr = fakes.FakeError
			contracts, err := prompter.GetContractDetails()
			Expect(err).To(HaveOccurred())
			Expect(err).To(MatchError(fakes.FakeError))
			Expect(contracts).To(BeEmpty())
		})
	})
})

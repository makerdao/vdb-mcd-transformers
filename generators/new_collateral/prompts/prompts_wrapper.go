package prompts

import (
	"fmt"

	"github.com/AlecAivazis/survey/v2"
	"github.com/makerdao/vdb-mcd-transformers/generators/new_collateral/types"
)

type PromptAsker func(questions []*survey.Question, response interface{}, opts ...survey.AskOpt) error

type Prompter struct {
	PromptAsker           PromptAsker
	FlipContractAnswers   ContractAnswers
	MedianRequired        bool
	MedianContractAnswers ContractAnswers
	OsmRequired           bool
	OsmContractAnswers    ContractAnswers
}

func NewPrompter() Prompter {
	return Prompter{
		PromptAsker: survey.Ask,
	}
}

var (
	CollateralQuestions = []*survey.Question{
		{
			Name:     "name",
			Prompt:   &survey.Input{Message: "Collateral Name:"},
			Validate: survey.Required,
		},
		{
			Name:     "version",
			Prompt:   &survey.Input{Message: "Collateral Version:"},
			Validate: survey.Required,
		},
	}
	FlipContractQuestions          = createContractQuestions("Flip")
	MedianContractRequiredQuestion = createContractRequiredQuestions("median")
	MedianContractQuestions        = createContractQuestions("Median")
	OsmContractRequiredQuestion    = createContractRequiredQuestions("osm")
	OsmContractQuestions           = createContractQuestions("Osm")
)

type CollateralAnswers struct {
	Name    string
	Version string
}

type ContractAnswers struct {
	Address string
	Abi     string
	Block   int
}

type ContractRequiredAnswer struct {
	Required bool
}

func (p *Prompter) GetCollateralDetails() (types.Collateral, error) {
	var answers CollateralAnswers
	err := p.PromptAsker(CollateralQuestions, &answers)
	if err != nil {
		return types.Collateral{}, err
	}

	return types.Collateral{
		Name:    answers.Name,
		Version: answers.Version,
	}, nil
}

func (p *Prompter) GetContractDetails() (types.Contracts, error) {
	contracts := make(map[string]types.Contract)

	flipErr := p.PromptAsker(FlipContractQuestions, &p.FlipContractAnswers)
	if flipErr != nil {
		return types.Contracts{}, flipErr
	}
	contracts["flip"] = types.Contract{
		Address:  p.FlipContractAnswers.Address,
		Abi:      p.FlipContractAnswers.Abi,
		Deployed: p.FlipContractAnswers.Block,
	}

	medianRequiredErr := p.PromptAsker(MedianContractRequiredQuestion, &p.MedianRequired)
	if medianRequiredErr != nil {
		return types.Contracts{}, medianRequiredErr
	}
	if p.MedianRequired {
		medianErr := p.PromptAsker(MedianContractQuestions, &p.MedianContractAnswers)
		if medianErr != nil {
			return types.Contracts{}, medianErr
		}
		contracts["median"] = types.Contract{
			Address:  p.MedianContractAnswers.Address,
			Abi:      p.MedianContractAnswers.Abi,
			Deployed: p.MedianContractAnswers.Block,
		}
	}

	osmRequiredErr := p.PromptAsker(OsmContractRequiredQuestion, &p.OsmRequired)
	if osmRequiredErr != nil {
		return types.Contracts{}, osmRequiredErr
	}
	if p.OsmRequired {
		osmErr := p.PromptAsker(OsmContractQuestions, &p.OsmContractAnswers)
		if osmErr != nil {
			return types.Contracts{}, osmErr
		}
		contracts["osm"] = types.Contract{
			Address:  p.OsmContractAnswers.Address,
			Abi:      p.OsmContractAnswers.Abi,
			Deployed: p.OsmContractAnswers.Block,
		}
	}

	return contracts, nil
}

func createContractQuestions(contractType string) []*survey.Question {
	return []*survey.Question{
		{
			Name:     "address",
			Prompt:   &survey.Input{Message: fmt.Sprintf("%s Contract Address:", contractType)},
			Validate: survey.Required,
		},
		{
			Name:     "abi",
			Prompt:   &survey.Multiline{Message: fmt.Sprintf("%s Contract Abi:", contractType)},
			Validate: survey.Required,
		},
		{
			Name:     "block",
			Prompt:   &survey.Input{Message: fmt.Sprintf("%s Contract Deployment Block:", contractType)},
			Validate: survey.Required,
		},
	}
}

func createContractRequiredQuestions(contractType string) []*survey.Question {
	return []*survey.Question{
		{
			Name: "required",
			Prompt: &survey.Confirm{
				Message: fmt.Sprintf("Is a %s contract required?", contractType),
			},
			Validate: survey.Required,
		},
	}
}

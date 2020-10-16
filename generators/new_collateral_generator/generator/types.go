package generator

type Contract struct {
	Address  string
	Abi      string
	Deployed int
}

type Contracts map[string]Contract

type TransformersConfig struct {
	ExporterMetadata     ExporterMetaData `toml:"exporter"`
	TransformerExporters TransformerExporters
	Contracts            Contracts `toml:"contract"`
}

type ExporterMetaData struct {
	Home             string
	Name             string
	Save             bool
	Schema           string
	TransformerNames []string
}

type TransformerExporter struct {
	Path       string
	Type       string `toml:"type"`
	Repository string
	Migrations string
	Contracts  []string
	Rank       string
}

type TransformerExporters map[string]TransformerExporter


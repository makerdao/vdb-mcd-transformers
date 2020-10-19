package config

type Contract struct {
	Address  string `toml:"address"`
	Abi      string `toml:"abi"`
	Deployed int    `toml:"deployed"`
}

type Contracts map[string]Contract

type TransformersConfig struct {
	ExporterMetadata     ExporterMetaData `toml:"exporter"`
	TransformerExporters TransformerExporters
	Contracts            Contracts `toml:"contract"`
}

type TransformersConfigForTomlEncoding struct {
	ExporterMetadata     ExporterMetaData     `toml:"exporter"`
	TransformerExporters TransformerExporters `toml:"exporter"`
	Contracts            Contracts            `toml:"contract"`
}

type ExporterMetaData struct {
	Home             string   `toml:"home"`
	Name             string   `toml:"name"`
	Save             bool     `toml:"save"`
	Schema           string   `toml:"schema"`
	TransformerNames []string `toml:"transformerNames"`
}

type TransformerExporter struct {
	Path       string   `toml:"path"`
	Type       string   `toml:"type"`
	Repository string   `toml:"repository"`
	Migrations string   `toml:"migrations"`
	Contracts  []string `toml:"contracts"`
	Rank       string   `toml:"rank"`
}

type TransformerExporters map[string]TransformerExporter

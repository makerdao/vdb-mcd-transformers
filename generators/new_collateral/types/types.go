package types

// TransformersConfigForToml is used as the struct to decode and encode the toml config file
type TransformersConfigForToml struct {
	Exporter  map[string]interface{} `toml:"exporter"`
	Contracts Contracts              `toml:"contract"`
}

// TransformersConfig is used as a domain object while updating different pieces in the config in the updater
type TransformersConfig struct {
	ExporterMetadata     ExporterMetaData `toml:"exporter"`
	TransformerExporters TransformerExporters
	Contracts            Contracts `toml:"contract"`
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

type Contract struct {
	Address  string `toml:"address"`
	Abi      string `toml:"abi"`
	Deployed int    `toml:"deployed"`
}

type Contracts map[string]Contract

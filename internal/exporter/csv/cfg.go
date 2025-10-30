package csv

type Specification struct {
	Delimiter string                              `yaml:"delimiter" json:"delimiter"`
	Transform map[string][]SpecificationTransform `yaml:"transform" json:"transform"`
}

type SpecificationTransform struct {
	OnlyColumns   []string          `yaml:"only_columns" json:"only_columns"`
	SkipColumns   []string          `yaml:"skip_columns" json:"skip_columns"`
	RenameColumns map[string]string `yaml:"rename_columns" json:"rename_columns"`
}

package laravel

type ModelsSpecification struct {
	Namespace string `yaml:"namespace" json:"namespace"`
	TimeAs    string `yaml:"time_as" json:"time_as"` // datetime, carbon
}

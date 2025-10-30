package goentity

type EntitySpecification struct {
	GoModule string `yaml:"go_module" json:"go_module"`
	Package  string `yaml:"package" json:"package"` // default: entities
}

type RepositorySpecRepoInterfacesPlace string

const (
	RepositorySpecRepoInterfacesPlaceUnspecified    RepositorySpecRepoInterfacesPlace = ""
	RepositorySpecRepoInterfacesPlaceWithEntity     RepositorySpecRepoInterfacesPlace = "with_entity"
	RepositorySpecRepoInterfacesPlaceWithRepository RepositorySpecRepoInterfacesPlace = "with_repository"
	RepositorySpecRepoInterfacesPlaceEntity         RepositorySpecRepoInterfacesPlace = "entity"
)

type EntityRepositorySpecification struct {
	GoModule string `yaml:"go_module" json:"go_module"`
	Entities struct {
		Package string `yaml:"package" json:"package"`
	} `yaml:"entities" json:"entities"`
	Repositories struct {
		Package   string `yaml:"package" json:"package"`
		Container struct {
			StructName string `yaml:"struct_name" json:"struct_name"`
		} `yaml:"container" json:"container"`
		Interfaces struct {
			Place     RepositorySpecRepoInterfacesPlace `yaml:"place" json:"place"`
			WithMocks bool                              `yaml:"with_mocks" json:"with_mocks"`
		} `yaml:"interfaces" json:"interfaces"`
	} `yaml:"repositories" json:"repositories"`
}

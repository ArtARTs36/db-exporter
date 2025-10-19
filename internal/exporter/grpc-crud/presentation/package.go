package presentation

type Package struct {
	name string

	files []*File

	cfg *config

	enumLocations map[string]string
}

func NewPackage(name string, configurators ...Configurator) *Package {
	return &Package{
		name:          name,
		files:         make([]*File, 0),
		cfg:           newConfig(configurators),
		enumLocations: make(map[string]string),
	}
}

func (p *Package) CreateFile(name string) *File {
	file := NewFile(p, name)

	p.files = append(p.files, file)

	return file
}

func (p *Package) LocateEnum(enumName string) (string, bool) {
	filepath, ok := p.enumLocations[enumName]
	return filepath, ok
}

func (p *Package) registerEnumLocation(enumName string, filename string) {
	p.enumLocations[enumName] = filename
}

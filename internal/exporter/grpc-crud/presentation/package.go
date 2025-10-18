package presentation

type Package struct {
	name string

	files []*File

	cfg *config
}

func NewPackage(name string, configurators ...Configurator) *Package {
	return &Package{
		name:  name,
		files: make([]*File, 0),
		cfg:   newConfig(configurators),
	}
}

func (p *Package) CreateFile(name string) *File {
	file := NewFile(p, name)

	p.files = append(p.files, file)

	return file
}

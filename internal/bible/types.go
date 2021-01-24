package bible

type Characters map[string][]Character

type Character struct {
	Name      string      `yaml:"name"`
	Filename  string      `yaml:"-"`
	Meaning   string      `yaml:"meaning,omitempty"`
	Sex       string      `yaml:"sex,omitempty"`
	Locations []Reference `yaml:"locations,omitempty"`
	Parents   []Reference `yaml:"parents,omitempty"`
	Spouse    []Reference `yaml:"spouse,omitempty"`
	Children  []Reference `yaml:"children,omitempty"`
	Info      []Note      `yaml:"info,omitempty"`
}

type Reference struct {
	Name      string `yaml:"name"`
	Reference string `yaml:"ref"`
}

type Note struct {
	Note       string `yaml:"note"`
	Reference  string `yaml:"ref"`
	Commentary string `yaml:"commentary,omitempty"`
}

type Locations map[string][]Location

type Location struct {
	Name     string `yaml:"name"`
	Filename string `yaml:"-"`
	Info     []Note `yaml:"info,omitempty"`
}

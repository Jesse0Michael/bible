package bible

import "time"

type Characters map[string][]Character

type Character struct {
	Name       string      `yaml:"name" json:"name"`
	Filename   string      `yaml:"-" json:"-"`
	CreateTime time.Time   `yaml:"created" json:"created"`
	UpdateTime time.Time   `yaml:"updated" json:"updated"`
	Alias      string      `yaml:"alias,omitempty" json:"alias,omitempty"`
	Meaning    string      `yaml:"meaning,omitempty" json:"meaning,omitempty"`
	Sex        string      `yaml:"sex,omitempty" json:"sex,omitempty"`
	Group      string      `yaml:"group,omitempty" json:"group,omitempty"`
	Locations  []Reference `yaml:"locations,omitempty" json:"locations,omitempty"`
	Parents    []Reference `yaml:"parents,omitempty" json:"parents,omitempty"`
	Spouse     []Reference `yaml:"spouse,omitempty" json:"spouse,omitempty" `
	Children   []Reference `yaml:"children,omitempty" json:"children,omitempty"`
	Associates []Reference `yaml:"associates,omitempty" json:"associates,omitempty"`
	Info       []Note      `yaml:"info,omitempty" json:"info,omitempty"`
}

type Reference struct {
	Name      string `yaml:"name" json:"name"`
	Reference string `yaml:"ref" json:"ref"`
}

type Note struct {
	Note       string `yaml:"note" json:"note"`
	Reference  string `yaml:"ref" json:"ref"`
	Commentary string `yaml:"commentary,omitempty" json:"commentary,omitempty"`
}

type Locations map[string][]Location

type Location struct {
	Name       string    `yaml:"name" json:"name"`
	Filename   string    `yaml:"-" json:"-"`
	Meaning    string    `yaml:"meaning,omitempty" json:"meaning,omitempty"`
	CreateTime time.Time `yaml:"created" json:"created"`
	UpdateTime time.Time `yaml:"updated" json:"updated"`
	Info       []Note    `yaml:"info,omitempty" json:"info,omitempty"`
	Alias      string    `yaml:"alias,omitempty" json:"alias,omitempty"`
}

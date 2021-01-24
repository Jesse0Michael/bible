package bible

import (
	"errors"
	"io/ioutil"
	"path/filepath"
	"strings"

	"github.com/urfave/cli/v2"
	"gopkg.in/yaml.v2"
)

const (
	locationDir = "locations"
)

func CreateLocation(c *cli.Context) error {
	name := strings.Title(c.Args().First())
	if name == "" {
		return errors.New("Required arg NAME not set")
	}

	f, err := NewFile(locationDir, name)
	if err != nil {
		return err
	}
	defer f.Close()

	location, err := processLocation(c, Location{
		Name:     name,
		Filename: filepath.Base(f.Name()),
	})
	if err != nil {
		return err
	}

	b, err := yaml.Marshal(location)
	if err != nil {
		return err
	}

	_, err = f.Write(b)
	return err
}

func UpdateLocation(c *cli.Context) error {
	name := strings.Title(c.Args().First())
	if name == "" {
		return errors.New("Required arg NAME not set")
	}

	ref, err := FindReference(locationDir, name)
	if err != nil {
		return err
	}
	b, _ := ioutil.ReadFile(filepath.Join(locationDir, ref))
	var location Location
	err = yaml.Unmarshal(b, &location)
	if err != nil {
		return err
	}

	location.Filename = ref
	location, err = processLocation(c, location)
	if err != nil {
		return err
	}

	b, err = yaml.Marshal(location)
	if err != nil {
		return err
	}

	return ioutil.WriteFile(filepath.Join(locationDir, ref), b, 0)
}

func processLocation(c *cli.Context, location Location) (Location, error) {
	if c.String("note") != "" {
		note := Note{
			Note:       c.String("note"),
			Reference:  c.String("reference"),
			Commentary: c.String("commentary"),
		}
		location.Info = append(location.Info, note)
	}

	return location, nil
}

package bible

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/urfave/cli/v2"
	"gopkg.in/yaml.v2"
)

const (
	locationDir = "assets/locations"
)

func GetLocation(c *cli.Context) error {
	name := c.Args().First()
	if name == "" {
		return errors.New("required arg NAME not set")
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

	switch c.String("output") {
	case "json":
		b, err = json.Marshal(location)
		if err != nil {
			return err
		}
		fmt.Println(string(b))
	default:
		b, err = yaml.Marshal(location)
		if err != nil {
			return err
		}
		fmt.Println(string(b))
	}

	return nil
}

func CreateLocation(c *cli.Context) error {
	name := strings.Title(c.Args().First())
	if name == "" {
		return errors.New("required arg NAME not set")
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
		_ = os.Remove(f.Name())
		return err
	}

	location.CreateTime = time.Now()
	location.UpdateTime = location.CreateTime
	b, err := yaml.Marshal(location)
	if err != nil {
		_ = os.Remove(f.Name())
		return err
	}

	_, err = f.Write(b)
	return err
}

func UpdateLocation(c *cli.Context) error {
	name := strings.Title(c.Args().First())
	if name == "" {
		return errors.New("required arg NAME not set")
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

	location.UpdateTime = time.Now()
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

	if c.String("meaning") != "" {
		location.Meaning = c.String("meaning")
	}

	return location, nil
}

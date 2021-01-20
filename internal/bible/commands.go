package bible

import (
	"errors"
	"io/ioutil"
	"path/filepath"

	"github.com/urfave/cli/v2"
	"gopkg.in/yaml.v2"
)

const (
	characterDir = "characters"
	// locationDir  = "locations"
)

func CreateCharacter(c *cli.Context) error {
	name := c.Args().First()
	if name == "" {
		return errors.New("Required arg NAME not set")
	}

	f, err := NewFile(characterDir, name)
	if err != nil {
		return err
	}
	defer f.Close()

	character := Character{
		Name: name,
		Sex:  c.String("sex"),
	}
	parent := c.String("parent")
	if parent != "" {
		ref, err := FindReference(characterDir, parent)
		if err != nil {
			return err
		}
		b, _ := ioutil.ReadFile(filepath.Join(characterDir, ref))
		var parentCharacter Character
		err = yaml.Unmarshal(b, &parentCharacter)
		if err != nil {
			return err
		}
		character.Parents = []Reference{
			{Name: parentCharacter.Name, Reference: ref},
		}

		// Update parent reference
		parentCharacter.Children = append(parentCharacter.Children, Reference{Name: name, Reference: f.Name()})
		b, err = yaml.Marshal(parentCharacter)
		if err != nil {
			return err
		}

		if err = ioutil.WriteFile(filepath.Join(characterDir, ref), b, 0); err != nil {
			return err
		}
	}

	b, err := yaml.Marshal(character)
	if err != nil {
		return err
	}

	_, err = f.Write(b)
	return err
}

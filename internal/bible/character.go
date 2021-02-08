package bible

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"path/filepath"
	"strings"

	"github.com/urfave/cli/v2"
	"gopkg.in/yaml.v2"
)

const (
	characterDir = "characters"
)

func GetCharacter(c *cli.Context) error {
	name := c.Args().First()
	if name == "" {
		return errors.New("Required arg NAME not set")
	}

	ref, err := FindReference(characterDir, name)
	if err != nil {
		return err
	}
	b, _ := ioutil.ReadFile(filepath.Join(characterDir, ref))
	var character Character
	err = yaml.Unmarshal(b, &character)
	if err != nil {
		return err
	}

	switch c.String("output") {
	case "json":
		b, err = json.Marshal(character)
		if err != nil {
			return err
		}
		fmt.Println(string(b))
	default:
		b, err = yaml.Marshal(character)
		if err != nil {
			return err
		}
		fmt.Println(string(b))
	}

	return nil
}

func CreateCharacter(c *cli.Context) error {
	name := strings.Title(c.Args().First())
	if name == "" {
		return errors.New("Required arg NAME not set")
	}

	f, err := NewFile(characterDir, name)
	if err != nil {
		return err
	}
	defer f.Close()

	character, err := processCharacter(c, Character{
		Name:     name,
		Filename: filepath.Base(f.Name()),
	})
	if err != nil {
		return err
	}

	b, err := yaml.Marshal(character)
	if err != nil {
		return err
	}

	_, err = f.Write(b)
	return err
}

func UpdateCharacter(c *cli.Context) error {
	name := strings.Title(c.Args().First())
	if name == "" {
		return errors.New("Required arg NAME not set")
	}

	ref, err := FindReference(characterDir, name)
	if err != nil {
		return err
	}
	b, _ := ioutil.ReadFile(filepath.Join(characterDir, ref))
	var character Character
	err = yaml.Unmarshal(b, &character)
	if err != nil {
		return err
	}

	character.Filename = ref
	character, err = processCharacter(c, character)
	if err != nil {
		return err
	}

	b, err = yaml.Marshal(character)
	if err != nil {
		return err
	}

	return ioutil.WriteFile(filepath.Join(characterDir, ref), b, 0)
}

func processCharacter(c *cli.Context, character Character) (Character, error) {
	ref := Reference{Name: character.Name, Reference: character.Filename}
	if c.String("sex") != "" {
		character.Sex = c.String("sex")
	}

	parent, err := addParent(c.String("parent"), ref)
	if err != nil {
		return character, err
	}
	if parent != nil {
		character.Parents = append(character.Parents, *parent)
	}

	spouse, err := addSpouse(c.String("spouse"), ref)
	if err != nil {
		return character, err
	}
	if parent != nil {
		character.Spouse = append(character.Spouse, *spouse)
	}

	associate, err := addAssociate(c.String("associate"), ref)
	if err != nil {
		return character, err
	}
	if associate != nil {
		character.Associates = append(character.Associates, *associate)
	}

	if c.String("note") != "" {
		note := Note{
			Note:       c.String("note"),
			Reference:  c.String("reference"),
			Commentary: c.String("commentary"),
		}
		character.Info = append(character.Info, note)
	}

	return character, nil
}

func addParent(parent string, character Reference) (*Reference, error) {
	if parent == "" {
		return nil, nil
	}
	ref, err := FindReference(characterDir, parent)
	if err != nil {
		return nil, err
	}
	b, _ := ioutil.ReadFile(filepath.Join(characterDir, ref))
	var parentCharacter Character
	err = yaml.Unmarshal(b, &parentCharacter)
	if err != nil {
		return nil, err
	}

	// Update parent reference
	parentCharacter.Children = append(parentCharacter.Children, character)
	b, err = yaml.Marshal(parentCharacter)
	if err != nil {
		return nil, err
	}

	if err = ioutil.WriteFile(filepath.Join(characterDir, ref), b, 0); err != nil {
		return nil, err
	}
	return &Reference{Name: parentCharacter.Name, Reference: ref}, nil
}

func addSpouse(spouse string, character Reference) (*Reference, error) {
	if spouse == "" {
		return nil, nil
	}
	ref, err := FindReference(characterDir, spouse)
	if err != nil {
		return nil, err
	}
	b, _ := ioutil.ReadFile(filepath.Join(characterDir, ref))
	var spouseCharacter Character
	err = yaml.Unmarshal(b, &spouseCharacter)
	if err != nil {
		return nil, err
	}

	// Update spouse reference
	spouseCharacter.Spouse = append(spouseCharacter.Spouse, character)
	b, err = yaml.Marshal(spouseCharacter)
	if err != nil {
		return nil, err
	}

	if err = ioutil.WriteFile(filepath.Join(characterDir, ref), b, 0); err != nil {
		return nil, err
	}
	return &Reference{Name: spouseCharacter.Name, Reference: ref}, nil
}

func addAssociate(associate string, character Reference) (*Reference, error) {
	if associate == "" {
		return nil, nil
	}
	ref, err := FindReference(characterDir, associate)
	if err != nil {
		return nil, err
	}
	b, _ := ioutil.ReadFile(filepath.Join(characterDir, ref))
	var associateCharacter Character
	err = yaml.Unmarshal(b, &associateCharacter)
	if err != nil {
		return nil, err
	}

	// Update associate reference
	associateCharacter.Associates = append(associateCharacter.Associates, character)
	b, err = yaml.Marshal(associateCharacter)
	if err != nil {
		return nil, err
	}

	if err = ioutil.WriteFile(filepath.Join(characterDir, ref), b, 0); err != nil {
		return nil, err
	}
	return &Reference{Name: associateCharacter.Name, Reference: ref}, nil
}

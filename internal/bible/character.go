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
	characterDir = "assets/characters"
)

func GetCharacter(c *cli.Context) error {
	name := c.Args().First()
	if name == "" {
		return errors.New("required arg NAME not set")
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
		return errors.New("required arg NAME not set")
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
		_ = os.Remove(f.Name())
		return err
	}

	character.CreateTime = time.Now()
	character.UpdateTime = character.CreateTime
	b, err := yaml.Marshal(character)
	if err != nil {
		_ = os.Remove(f.Name())
		return err
	}

	_, err = f.Write(b)
	return err
}

func UpdateCharacter(c *cli.Context) error {
	name := strings.Title(c.Args().First())
	if name == "" {
		return errors.New("required arg NAME not set")
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

	character.UpdateTime = time.Now()
	b, err = yaml.Marshal(character)
	if err != nil {
		return err
	}

	return ioutil.WriteFile(filepath.Join(characterDir, ref), b, 0)
}

func AuditCharacters(c *cli.Context) error {
	files, err := ioutil.ReadDir(characterDir)
	if err != nil {
		return err
	}
	for _, f := range files {
		b, _ := ioutil.ReadFile(filepath.Join(characterDir, f.Name()))
		var character Character
		err = yaml.Unmarshal(b, &character)
		if err != nil {
			return err
		}

		// Audit here!

		b, err = yaml.Marshal(character)
		if err != nil {
			return err
		}

		_ = ioutil.WriteFile(filepath.Join(characterDir, f.Name()), b, 0)
	}
	return nil
}

func processCharacter(c *cli.Context, character Character) (Character, error) {
	ref := Reference{Name: character.Name, Reference: character.Filename}
	if c.String("sex") != "" {
		switch c.String("sex") {
		case "female":
			character.Sex = "female"
		case "male":
			character.Sex = "male"
		default:
			return character, fmt.Errorf("invalid sex: \"%s\" valid values: (male,female)", c.String("sex"))
		}
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
	if spouse != nil {
		character.Spouse = append(character.Spouse, *spouse)
	}

	associate, err := addAssociate(c.String("associate"), ref)
	if err != nil {
		return character, err
	}
	if associate != nil {
		character.Associates = append(character.Associates, *associate)
	}

	location, err := addLocation(c.String("location"), ref)
	if err != nil {
		return character, err
	}
	if location != nil {
		character.Locations = append(character.Locations, *location)
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

func addLocation(location string, character Reference) (*Reference, error) {
	if location == "" {
		return nil, nil
	}
	ref, err := FindReference(locationDir, location)
	if err != nil {
		return nil, err
	}
	b, _ := ioutil.ReadFile(filepath.Join(locationDir, ref))
	var loc Location
	err = yaml.Unmarshal(b, &loc)
	if err != nil {
		return nil, err
	}

	return &Reference{Name: loc.Name, Reference: fmt.Sprintf("../locations/%s", ref)}, nil
}

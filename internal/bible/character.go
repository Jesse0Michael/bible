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

	"github.com/spf13/viper"
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
	b, _ := ioutil.ReadFile(filepath.Join(viper.GetString("DIR"), characterDir, ref))
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
	if len(c.Args().Slice()) == 0 {
		return errors.New("required arg NAME not set")
	}
	for _, arg := range c.Args().Slice() {
		name := strings.Title(arg)
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

		if _, err = f.Write(b); err != nil {
			return err
		}
		fmt.Printf("created character: %s\n", filepath.Base(f.Name()))
	}
	return nil
}

func UpdateCharacter(c *cli.Context) error {
	if len(c.Args().Slice()) == 0 {
		return errors.New("required arg NAME not set")
	}
	for _, arg := range c.Args().Slice() {
		name := strings.Title(arg)
		ref, err := FindReference(characterDir, name)
		if err != nil {
			return err
		}
		b, _ := ioutil.ReadFile(filepath.Join(viper.GetString("DIR"), characterDir, ref))
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

		if err := ioutil.WriteFile(filepath.Join(viper.GetString("DIR"), characterDir, ref), b, 0); err != nil {
			return err
		}
		fmt.Printf("updated character: %s\n", ref)
	}
	return nil
}

func AuditCharacters(c *cli.Context) error {
	files, err := ioutil.ReadDir(characterDir)
	if err != nil {
		return err
	}
	for _, f := range files {
		b, _ := ioutil.ReadFile(filepath.Join(viper.GetString("DIR"), characterDir, f.Name()))
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

		_ = ioutil.WriteFile(filepath.Join(viper.GetString("DIR"), characterDir, f.Name()), b, 0)
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

	for _, arg := range c.StringSlice("parent") {
		parent, err := addParent(arg, ref)
		if err != nil {
			return character, err
		}
		if parent != nil {
			if !hasReference(parent.Reference, character.Parents) {
				character.Parents = append(character.Parents, *parent)
			}
		}
	}

	for _, arg := range c.StringSlice("spouse") {
		spouse, err := addSpouse(arg, ref)
		if err != nil {
			return character, err
		}
		if spouse != nil {
			if !hasReference(spouse.Reference, character.Spouse) {
				character.Spouse = append(character.Spouse, *spouse)
			}
		}
	}

	for _, arg := range c.StringSlice("associate") {
		associate, err := addAssociate(arg, ref)
		if err != nil {
			return character, err
		}
		if associate != nil {
			if !hasReference(associate.Reference, character.Associates) {
				character.Associates = append(character.Associates, *associate)
			}
		}
	}

	for _, arg := range c.StringSlice("location") {
		location, err := addLocation(arg, ref)
		if err != nil {
			return character, err
		}
		if location != nil {
			if !hasReference(location.Reference, character.Locations) {
				character.Locations = append(character.Locations, *location)
			}
		}
	}

	if c.String("note") != "" {
		note := Note{
			Note:       c.String("note"),
			Reference:  c.String("reference"),
			Commentary: c.String("commentary"),
		}
		character.Info = append(character.Info, note)
	}

	if c.String("meaning") != "" {
		character.Meaning = c.String("meaning")
	}

	if c.String("group") != "" {
		character.Group = c.String("group")
	}

	if c.String("alias") != "" {
		character.Alias = c.String("alias")
	}

	return character, nil
}

func addParent(parent string, character Reference) (*Reference, error) {
	ref, err := FindReference(characterDir, parent)
	if err != nil {
		return nil, err
	}
	b, _ := ioutil.ReadFile(filepath.Join(viper.GetString("DIR"), characterDir, ref))
	var parentCharacter Character
	err = yaml.Unmarshal(b, &parentCharacter)
	if err != nil {
		return nil, err
	}

	// Update parent reference
	if !hasReference(character.Reference, parentCharacter.Children) {
		parentCharacter.Children = append(parentCharacter.Children, character)
		b, err = yaml.Marshal(parentCharacter)
		if err != nil {
			return nil, err
		}

		if err = ioutil.WriteFile(filepath.Join(viper.GetString("DIR"), characterDir, ref), b, 0); err != nil {
			return nil, err
		}
	}
	return &Reference{Name: parentCharacter.Name, Reference: ref}, nil
}

func addSpouse(spouse string, character Reference) (*Reference, error) {
	ref, err := FindReference(characterDir, spouse)
	if err != nil {
		return nil, err
	}
	b, _ := ioutil.ReadFile(filepath.Join(viper.GetString("DIR"), characterDir, ref))
	var spouseCharacter Character
	err = yaml.Unmarshal(b, &spouseCharacter)
	if err != nil {
		return nil, err
	}

	// Update spouse reference
	if !hasReference(character.Reference, spouseCharacter.Spouse) {
		spouseCharacter.Spouse = append(spouseCharacter.Spouse, character)
		b, err = yaml.Marshal(spouseCharacter)
		if err != nil {
			return nil, err
		}

		if err = ioutil.WriteFile(filepath.Join(viper.GetString("DIR"), characterDir, ref), b, 0); err != nil {
			return nil, err
		}
	}
	return &Reference{Name: spouseCharacter.Name, Reference: ref}, nil
}

func addAssociate(associate string, character Reference) (*Reference, error) {
	ref, err := FindReference(characterDir, associate)
	if err != nil {
		return nil, err
	}
	b, _ := ioutil.ReadFile(filepath.Join(viper.GetString("DIR"), characterDir, ref))
	var associateCharacter Character
	err = yaml.Unmarshal(b, &associateCharacter)
	if err != nil {
		return nil, err
	}

	// Update associate reference
	if !hasReference(character.Reference, associateCharacter.Associates) {
		associateCharacter.Associates = append(associateCharacter.Associates, character)
		b, err = yaml.Marshal(associateCharacter)
		if err != nil {
			return nil, err
		}

		if err = ioutil.WriteFile(filepath.Join(viper.GetString("DIR"), characterDir, ref), b, 0); err != nil {
			return nil, err
		}
	}
	return &Reference{Name: associateCharacter.Name, Reference: ref}, nil
}

func addLocation(location string, character Reference) (*Reference, error) {
	ref, err := FindReference(locationDir, location)
	if err != nil {
		return nil, err
	}
	b, _ := ioutil.ReadFile(filepath.Join(viper.GetString("DIR"), locationDir, ref))
	var loc Location
	err = yaml.Unmarshal(b, &loc)
	if err != nil {
		return nil, err
	}

	return &Reference{Name: loc.Name, Reference: fmt.Sprintf("../locations/%s", ref)}, nil
}

func hasReference(ref string, references []Reference) bool {
	for _, s := range references {
		if s.Reference == ref {
			return true
		}
	}
	return false
}

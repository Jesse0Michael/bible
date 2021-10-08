package bible

import (
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/spf13/viper"
	"github.com/urfave/cli/v2"
)

func GetConfig(c *cli.Context) error {
	if viper.GetString("DIR") != "" {
		fmt.Printf("DIR=%s\n", viper.Get("DIR"))
	}
	return nil
}

func SetConfig(c *cli.Context) error {
	key := c.Args().First()
	if key == "" {
		return errors.New("required arg KEY not set")
	}
	switch key {
	case "DIR":
	default:
		return errors.New("invalid KEY. supported config keys: DIR")
	}

	value := c.Args().Get(1)
	if value == "" {
		return errors.New("required arg VALUE not set")
	}

	switch key {
	case "DIR":
		if !strings.HasSuffix(value, "/") {
			value = value + "/"
		}
		if _, err := os.ReadDir(value); err != nil {
			return errors.New("required valid directory for DIR config")
		}
	}

	viper.Set(key, value)

	if err := viper.WriteConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			if err = viper.SafeWriteConfig(); err != nil {
				return err
			}
		} else {
			return err
		}
	}

	return nil
}

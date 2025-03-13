package lib

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/m-tsuru/trans/structs"
	"gopkg.in/yaml.v2"
)

func DecideConfigPath(config string) (*string, error) {
	//
	// Priority of Reading Config Path.
	// 1. Option `--config`
	// 2. `./.trans/config.yaml`
	// 3. config.yaml in the same folder as the executable file
	// 4. $HOME/.trans/config.yaml
	//
	// If nothing, raise error.
	//

	var ConfigNotFound = errors.New("config file was not found")

	if config != "" {
		return &config, nil
	}

	CurrentDirectoryConfigPath := filepath.Join("./", ".trans", "config.yaml")
	if Exists(CurrentDirectoryConfigPath) {
		return &CurrentDirectoryConfigPath, nil
	}

	execPath, err := GetExecutablePath()
	if err == nil {
		execConfigPath := filepath.Join(filepath.Dir(execPath), "config.yaml")
		if Exists(execConfigPath) {
			return &execConfigPath, nil
		}
	}

	homeConfigPath := filepath.Join(os.Getenv("HOME"), ".trans", "config.yaml")
	if Exists(homeConfigPath) {
		return &homeConfigPath, nil
	}

	return nil, ConfigNotFound
}

func ReadConfig(fn string) (*structs.Conf, error) {
	var c structs.Conf
	b, err := os.ReadFile(fn)
	if err != nil {
		return nil, err
	}
	yaml.Unmarshal(b, &c)
	return &c, nil
}

func FindImportKey(profileName string, config structs.Conf) (*structs.ImportProfile, error) {
	var profile structs.ImportProfile
	NoDefaultProfileError := errors.New("Default Profile is not set.")
	NoProfileError := fmt.Errorf("Profile `%s` is not found.", *&profileName)

	if *&profileName == "" {
		for _, v := range config.Import {
			if v.Default == true {
				profile = v
				break
			} else {
				return nil, NoDefaultProfileError
			}
		}
	} else {
		for _, v := range config.Import {
			if v.Name == *&profileName {
				profile = v
				break
			} else {
				return nil, NoProfileError
			}
		}
	}

	return &profile, nil
}

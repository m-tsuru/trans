package main

import (
	_ "embed"
	"flag"
	"fmt"
	"os"

	"github.com/m-tsuru/trans/lib"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

//go:embed VERSION
var version string

var h string = fmt.Sprintf(`
trans %s
Usage: trans <command> [<args>]

Commands:
  help    Show this help
  import  Copy from original to target
  check   Verify that the files are copied correctly

Example (import with profile in config.yaml):
  $ trans import [--profile <profile>]

Example (Verify with profile in config.yaml):
  $ trans check [--profile <profile>]

If you want to know about placing config files,
please access to https://github.com/m-tsuru/trans/README.md .

`, version)

func main() {
	os.Exit(_main())
}

func _main() int {
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stdout}).Level(zerolog.InfoLevel)

	const (
		help    = "help"
		_import = "import"
		check   = "check"
	)

	command := help
	options := []string{}

	if len(os.Args) > 1 {
		command = os.Args[1]
		options = os.Args[2:]
	}

	switch command {
	case help:
		fmt.Print(h)
	case _import:
		res := __Import(options)
		return res
	case check:
		res := __Check(options)
		return res
	default:
		fmt.Print(h)
		return 1
	}
	return 0
}

func __Import(options []string) int {
	cmd := flag.NewFlagSet("trans import [--profile <profile>]", flag.ExitOnError)
	configPath := cmd.String("config", "", "config file")
	profileName := cmd.String("profile", "", "profile name")

	err := cmd.Parse(options)

	fmt.Println(*profileName)

	configPath, err = lib.DecideConfigPath(*configPath)
	if err != nil {
		log.Error().Msgf("Failed to read config file: %s", err)
		return 1
	} else {
		log.Info().Msgf("Read Configuration File Path: %s", *configPath)
	}

	config, err := lib.ReadConfig(*configPath)
	if err != nil {
		log.Error().Msgf("Failed to read config file: %s", err)
		return 1
	} else {
		log.Info().Msgf("Read Configuration Successfully: %s", *configPath)
	}

	profile, err := lib.FindImportKey(*profileName, *config)
	if err != nil {
		log.Error().Msgf("Failed to read profile: %s", err)
		return 1
	}

	log.Info().Msgf("[Profile: %s] - Base Directory: %s -> %s", profile.Name, profile.Original, profile.Target)

	log.Info().Msgf("Enumerate Files: %s", profile.Original)
	err = lib.Import(profile.Original, *profile)
	if err != nil {
		log.Error().Msgf("Failed to import files: %s", err)
	}
	return 0
}

func __Check(options []string) int {
	cmd := flag.NewFlagSet("trans check [--profile <profile>]", flag.ExitOnError)
	configPath := cmd.String("config", "", "config file")
	profileName := cmd.String("profile", "", "profile name")

	err := cmd.Parse(options)

	fmt.Println(*profileName)

	configPath, err = lib.DecideConfigPath(*configPath)
	if err != nil {
		log.Error().Msgf("Failed to read config file: %s", err)
		return 1
	} else {
		log.Info().Msgf("Read Configuration File Path: %s", *configPath)
	}

	config, err := lib.ReadConfig(*configPath)
	if err != nil {
		log.Error().Msgf("Failed to read config file: %s", err)
		return 1
	} else {
		log.Info().Msgf("Read Configuration Successfully: %s", *configPath)
	}

	profile, err := lib.FindImportKey(*profileName, *config)
	if err != nil {
		log.Error().Msgf("Failed to read profile: %s", err)
		return 1
	}

	log.Info().Msgf("[Profile: %s] - Base Directory: %s -> %s", profile.Name, profile.Original, profile.Target)

	log.Info().Msgf("Enumerate Files: %s", profile.Original)

	err = lib.Check(profile.Original, *profile)
	if err != nil {
		log.Error().Msgf("Failed to check files: %s", err)
		return 1
	}
	return 0
}

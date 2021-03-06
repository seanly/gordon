package mainutil

import (
	"io"
	"os"
	"path/filepath"

	"github.com/alecthomas/kingpin"
	"github.com/kyoh86/gordon/internal/command"
	"github.com/kyoh86/gordon/internal/env"
	"github.com/kyoh86/xdg"
)

func setConfigFlag(cmd *kingpin.CmdClause, configFile *string) {
	cmd.Flag("config", "configuration file").
		Default(filepath.Join(xdg.ConfigHome(), "gordon", "config.yaml")).
		Envar("GORDON_CONFIG").
		StringVar(configFile)
}

func openYAML(filename string) (io.Reader, func() error, error) {
	var reader io.Reader
	var teardown func() error
	file, err := os.Open(filename)
	switch {
	case err == nil:
		teardown = file.Close
		reader = file
	case os.IsNotExist(err):
		reader = env.EmptyYAMLReader
		teardown = func() error { return nil }
	default:
		return nil, nil, err
	}
	return reader, teardown, nil
}

func WrapCommand(cmd *kingpin.CmdClause, f func(command.Env) error) (string, func() error) {
	var configFile string
	setConfigFlag(cmd, &configFile)
	return cmd.FullCommand(), func() (retErr error) {
		reader, teardown, err := openYAML(configFile)
		if err != nil {
			return err
		}
		defer func() {
			if err := teardown(); err != nil && retErr == nil {
				retErr = err
				return
			}
		}()

		access, err := env.GetAccess(reader, env.EnvarPrefix)
		if err != nil {
			return err
		}

		return f(&access)
	}
}

func WrapConfigurableCommand(cmd *kingpin.CmdClause, f func(command.Env, *env.Config) error) (string, func() error) {
	var configFile string
	setConfigFlag(cmd, &configFile)
	return cmd.FullCommand(), func() (retErr error) {
		reader, teardown, err := openYAML(configFile)
		if err != nil {
			return err
		}
		defer func() {
			if err := teardown(); err != nil && retErr == nil {
				retErr = err
				return
			}
		}()

		config, access, err := env.GetAppenv(reader, env.EnvarPrefix)
		if err != nil {
			return err
		}

		if err = f(&access, &config); err != nil {
			return err
		}

		if err := os.MkdirAll(filepath.Dir(configFile), 0744); err != nil {
			return err
		}
		file, err := os.OpenFile(configFile, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0644)
		if err != nil {
			return err
		}
		defer file.Close()
		return config.Save(file)
	}
}

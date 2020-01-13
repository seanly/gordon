package context

import (
	"errors"
	"github.com/zalando/go-keyring"
)

var (
	ErrEmptyValue        = errors.New("empty value")
	ErrInvalidOptionName = errors.New("invalid option name")
)

type OptionAccessor struct {
	optionName string
	getter     func(cfg *Config) string
	putter     func(cfg *Config, value string) error
	unsetter   func(cfg *Config) error
}

func (a OptionAccessor) Get(cfg *Config) string              { return a.getter(cfg) }
func (a OptionAccessor) Put(cfg *Config, value string) error { return a.putter(cfg, value) }
func (a OptionAccessor) Unset(cfg *Config) error             { return a.unsetter(cfg) }

var (
	configAccessor  map[string]OptionAccessor
	optionNames     []string
	optionAccessors = []OptionAccessor{
		rootOptionAccessor,
		archOptionAccessor,
		osOptionAccessor,
		gitHubHostOptionAccessor,
		gitHubUserOptionAccessor,
		gitHubTokenOptionAccessor,
		logLevelOptionAccessor,
		logDateOptionAccessor,
		logTimeOptionAccessor,
		logMicroSecondsOptionAccessor,
		logLongFileOptionAccessor,
		logShortFileOptionAccessor,
		logUTCOptionAccessor,
	}
)

func init() {
	m := map[string]OptionAccessor{}
	n := make([]string, 0, len(optionAccessors))
	for _, a := range optionAccessors {
		n = append(n, a.optionName)
		m[a.optionName] = a
	}
	configAccessor = m
	optionNames = n
}

func Option(optionName string) (*OptionAccessor, error) {
	a, ok := configAccessor[optionName]
	if !ok {
		return nil, ErrInvalidOptionName
	}
	return &a, nil
}

func OptionNames() []string {
	return optionNames
}

var (
	gitHubUserOptionAccessor = OptionAccessor{
		optionName: "github.user",
		getter: func(cfg *Config) string {
			return cfg.GitHubUser()
		},
		putter: func(cfg *Config, value string) error {
			if value == "" {
				return ErrEmptyValue
			}
			if err := ValidateOwner(value); err != nil {
				return err
			}
			cfg.GitHub.User = value
			return nil
		},
		unsetter: func(cfg *Config) error {
			cfg.GitHub.User = ""
			return nil
		},
	}

	gitHubTokenOptionAccessor = OptionAccessor{
		optionName: "github.token",
		getter: func(cfg *Config) string {
			if cfg.GitHubToken() == "" {
				return ""
			}
			return "*****"
		},
		putter: func(cfg *Config, value string) error {
			if value == "" {
				return ErrEmptyValue
			}
			return keyring.Set(keyGordonServiceName, keyGordonGitHubToken, value)
		},
		unsetter: func(cfg *Config) error {
			return keyring.Delete(keyGordonServiceName, keyGordonGitHubToken)
		},
	}

	gitHubHostOptionAccessor = OptionAccessor{
		optionName: "github.host",
		getter: func(cfg *Config) string {
			return cfg.GitHubHost()
		},
		putter: func(cfg *Config, value string) error {
			if value == "" {
				return ErrEmptyValue
			}
			cfg.GitHub.Host = value
			return nil
		},
		unsetter: func(cfg *Config) error {
			cfg.GitHub.Host = ""
			return nil
		},
	}

	logLevelOptionAccessor = OptionAccessor{
		optionName: "log.level",
		getter: func(cfg *Config) string {
			return cfg.LogLevel()
		},
		putter: func(cfg *Config, value string) error {
			if value == "" {
				return ErrEmptyValue
			}
			if err := ValidateLogLevel(value); err != nil {
				return err
			}
			cfg.Log.Level = value
			return nil
		},
		unsetter: func(cfg *Config) error {
			cfg.Log.Level = ""
			return nil
		},
	}

	logDateOptionAccessor = OptionAccessor{
		optionName: "log.date",
		getter: func(cfg *Config) string {
			return cfg.Log.Date.String()
		},
		putter: func(cfg *Config, value string) error {
			if value == "" {
				return ErrEmptyValue
			}
			return cfg.Log.Date.Decode(value)
		},
		unsetter: func(cfg *Config) error {
			cfg.Log.Date = EmptyBoolOption
			return nil
		},
	}

	logTimeOptionAccessor = OptionAccessor{
		optionName: "log.time",
		getter: func(cfg *Config) string {
			return cfg.Log.Time.String()
		},
		putter: func(cfg *Config, value string) error {
			if value == "" {
				return ErrEmptyValue
			}
			return cfg.Log.Time.Decode(value)
		},
		unsetter: func(cfg *Config) error {
			cfg.Log.Time = EmptyBoolOption
			return nil
		},
	}

	logMicroSecondsOptionAccessor = OptionAccessor{
		optionName: "log.microseconds",
		getter: func(cfg *Config) string {
			return cfg.Log.MicroSeconds.String()
		},
		putter: func(cfg *Config, value string) error {
			if value == "" {
				return ErrEmptyValue
			}
			return cfg.Log.MicroSeconds.Decode(value)
		},
		unsetter: func(cfg *Config) error {
			cfg.Log.MicroSeconds = EmptyBoolOption
			return nil
		},
	}

	logLongFileOptionAccessor = OptionAccessor{
		optionName: "log.longfile",
		getter: func(cfg *Config) string {
			return cfg.Log.LongFile.String()
		},
		putter: func(cfg *Config, value string) error {
			if value == "" {
				return ErrEmptyValue
			}
			return cfg.Log.LongFile.Decode(value)
		},
		unsetter: func(cfg *Config) error {
			cfg.Log.LongFile = EmptyBoolOption
			return nil
		},
	}

	logShortFileOptionAccessor = OptionAccessor{
		optionName: "log.shortfile",
		getter: func(cfg *Config) string {
			return cfg.Log.ShortFile.String()
		},
		putter: func(cfg *Config, value string) error {
			if value == "" {
				return ErrEmptyValue
			}
			return cfg.Log.ShortFile.Decode(value)
		},
		unsetter: func(cfg *Config) error {
			cfg.Log.ShortFile = EmptyBoolOption
			return nil
		},
	}

	logUTCOptionAccessor = OptionAccessor{
		optionName: "log.utc",
		getter: func(cfg *Config) string {
			return cfg.Log.UTC.String()
		},
		putter: func(cfg *Config, value string) error {
			if value == "" {
				return ErrEmptyValue
			}
			return cfg.Log.UTC.Decode(value)
		},
		unsetter: func(cfg *Config) error {
			cfg.Log.UTC = EmptyBoolOption
			return nil
		},
	}

	rootOptionAccessor = OptionAccessor{
		optionName: "root",
		getter: func(cfg *Config) string {
			return cfg.Root()
		},
		putter: func(cfg *Config, value string) error {
			if value == "" {
				return ErrEmptyValue
			}
			if err := ValidateRoot(value); err != nil {
				return err
			}
			cfg.VRoot = value
			return nil
		},
		unsetter: func(cfg *Config) error {
			cfg.VRoot = ""
			return nil
		},
	}

	archOptionAccessor = OptionAccessor{
		optionName: "arch",
		getter: func(cfg *Config) string {
			return cfg.Arch()
		},
		putter: func(cfg *Config, value string) error {
			if value == "" {
				return ErrEmptyValue
			}
			// TODO: if err := ValidateArch(value); err != nil {
			// 	return err
			// }
			cfg.VArch = value
			return nil
		},
		unsetter: func(cfg *Config) error {
			cfg.VArch = ""
			return nil
		},
	}

	osOptionAccessor = OptionAccessor{
		optionName: "os",
		getter: func(cfg *Config) string {
			return cfg.OS()
		},
		putter: func(cfg *Config, value string) error {
			if value == "" {
				return ErrEmptyValue
			}
			// TODO: if err := ValidateOS(value); err != nil {
			// return err
			// }
			cfg.VOS = value
			return nil
		},
		unsetter: func(cfg *Config) error {
			cfg.VOS = ""
			return nil
		},
	}
)
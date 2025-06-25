package bark

import (
	"github.com/Anniext/Arkitektur/system/config"
	"github.com/jzksnsjswkw/go-bark"
)

type BarkConfig struct {
	Host  string
	Token string
	Group string
	Copy  string
	Auto  bool
	Sound string
	Icon  string
}
type Option func(*BarkConfig)

func WithHostOption(host string) Option {
	return func(c *BarkConfig) {
		c.Host = host
	}
}

func WithTokenOption(token string) Option {
	return func(c *BarkConfig) {
		c.Token = token
	}
}

func WithGroupOption(group string) Option {
	return func(c *BarkConfig) {
		c.Group = group
	}
}

func WithCopyOption(copy string) Option {
	return func(c *BarkConfig) {
		c.Copy = copy
	}
}

func WithAutoCopyOption(auto bool) Option {
	return func(c *BarkConfig) {
		c.Auto = auto
	}
}

func WithSoundOption(sound string) Option {
	return func(c *BarkConfig) {
		c.Sound = sound
	}
}

func WithIconOption(icon string) Option {
	return func(c *BarkConfig) {
		c.Icon = icon
	}
}
func NewBarkOption(options ...Option) {
	defaultBarkConfig = &BarkConfig{}
	for _, option := range options {
		option(defaultBarkConfig)
	}
}

var defaultBarkConfig *BarkConfig

func GetDefaultBarkConfig() *BarkConfig {
	return defaultBarkConfig
}

func PushBark(title, msg, icon, group string) error {
	cnf := config.GetBarkInfo()
	if icon == "" {
		icon = cnf.Icon
	}
	if group != "" {
		group = cnf.Group
	}
	if defaultBark != nil {
		err := defaultBark.Push(&bark.Options{
			Msg:      msg,
			Token:    cnf.Token,
			Title:    title,
			Group:    group,
			Copy:     msg,
			AutoCopy: cnf.AutoCopy,
			Sound:    cnf.Sound,
			Icon:     icon,
		})
		if err != nil {
			return err
		}
	}
	return nil
}

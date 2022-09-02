package config

import (
	"fmt"
	"log"
	"text/template"

	"github.com/elseano/dollop/internal/templating"
	"github.com/spf13/viper"
)

type Config struct {
	LevelField     string        `yaml:"levelField"`
	MessageField   string        `yaml:"messageField"`
	TimestampField string        `yaml:"timestampField"`
	Groups         []*GroupSpec  `yaml:"groups"`
	Statuses       []*StatusSpec `yaml:"statuses"`
	Tags           []*TagSpec    `yaml:"tags"`

	LevelTmpl     *template.Template
	MessageTmpl   *template.Template
	TimestampTmpl *template.Template
}

type GroupSpec struct {
	ValueField string     `yaml:"valueField"`
	TitleField string     `yaml:"titleField"`
	Tags       []*TagSpec `yaml:"tags"`
	Name       string     `yaml:"name"`

	TitleTmpl *template.Template
	ValueTmpl *template.Template
}

type StatusSpec struct {
	Display string `yaml:"display"`

	DisplayTmpl *template.Template
}

type TagSpec struct {
	Value string `yaml:"source"`
	Key   string `yaml:"name"`

	ValueTmpl *template.Template
	KeyTmpl   *template.Template
}

func (c *Config) PrepareTemplates() {
	c.MessageTmpl = templating.BuildTemplate(c.MessageField)
	c.TimestampTmpl = templating.BuildTemplate(c.TimestampField)
	c.LevelTmpl = templating.BuildTemplate(c.LevelField)

	if c.Tags == nil {
		c.Tags = []*TagSpec{}
	}

	for _, t := range c.Tags {
		t.PrepareTemplates()
	}

	for _, g := range c.Groups {
		g.ValueTmpl = templating.BuildTemplate(g.ValueField)
		g.TitleTmpl = templating.BuildTemplate(g.TitleField)
		if g.Tags == nil {
			g.Tags = []*TagSpec{}
		}

		for _, t := range g.Tags {
			t.PrepareTemplates()
		}
	}

	for _, s := range c.Statuses {
		s.DisplayTmpl = templating.BuildTemplate(s.Display)
	}
}

func (t *TagSpec) PrepareTemplates() {
	t.KeyTmpl = templating.BuildTemplateText(t.Key)
	if t.Value != "" {
		t.ValueTmpl = templating.BuildTemplate(t.Value)
	}
}

func Get() (config Config) {
	config = Config{
		MessageField:   "msg",
		TimestampField: "time",
		LevelField:     "level",
		Groups: []*GroupSpec{
			{ValueField: "request_id", TitleField: "msg", Name: "Request"},
			{ValueField: "category", TitleField: "category", Name: "Category"},
		},
	}

	if err := viper.Unmarshal(&config); err != nil {
		log.Fatal("Error parsing config", err)
	}

	if err := config.Validate(); err != nil {
		log.Fatalf("Config error: %s", err.Error())
	}

	config.PrepareTemplates()

	return
}

func (c Config) Validate() error {
	if c.LevelField == "" {
		return fmt.Errorf("'levelField' cannot be blank")
	}

	if c.MessageField == "" {
		return fmt.Errorf("'messageField' cannot be blank")
	}

	if c.TimestampField == "" {
		return fmt.Errorf("'timestampField' cannot be blank")
	}

	for i, g := range c.Groups {
		if err := g.Validate(); err != nil {
			return fmt.Errorf("group entry #%d: %w", i+1, err)
		}
	}

	for i, s := range c.Statuses {
		if err := s.Validate(); err != nil {
			return fmt.Errorf("status entry #%d: %w", i+1, err)
		}
	}

	return nil
}

func (g GroupSpec) Validate() error {
	if g.TitleField == "" {
		return fmt.Errorf("'titleField' cannot be blank")
	}

	if g.ValueField == "" {
		return fmt.Errorf("'valueField' cannot be blank")
	}

	if g.Name == "" {
		return fmt.Errorf("'name' cannot be blank")
	}

	for i, t := range g.Tags {
		if err := t.Validate(); err != nil {
			return fmt.Errorf("tag %d: %w", i, err)
		}
	}

	return nil
}

func (s StatusSpec) Validate() error {
	if s.Display == "" {
		return fmt.Errorf("'display' cannot be blank")
	}

	return nil
}

func (s TagSpec) Validate() error {
	if s.Key == "" {
		return fmt.Errorf("'key' must be specified")
	}

	return nil
}

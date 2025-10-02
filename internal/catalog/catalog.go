package catalog

import (
	"errors"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

type Strategy struct {
	Brew    string            `yaml:"brew,omitempty"`
	Apt     string            `yaml:"apt,omitempty"`
	Dnf     string            `yaml:"dnf,omitempty"`
	Pacman  string            `yaml:"pacman,omitempty"`
	Zypper  string            `yaml:"zypper,omitempty"`
	Pipx    string            `yaml:"pipx,omitempty"`
	Uv      string            `yaml:"uv,omitempty"`
	Npm     string            `yaml:"npm,omitempty"`
	Volta   string            `yaml:"volta,omitempty"`
	Release map[string]string `yaml:"release,omitempty"` // url, bin
}

type Item struct {
	ID          string   `yaml:"id"`
	Name        string   `yaml:"name"`
	Description string   `yaml:"description,omitempty"`
	Verify      string   `yaml:"verify"` // e.g., "git-town --version"
	Strategy    Strategy `yaml:"strategies"`
}

type Config struct {
	Items  []Item   `yaml:"items"`
	Curate []string `yaml:"curate,omitempty"`
}

func Load() (*Config, error) {
	path := filepath.Join("config", "catalog.yaml")
	b, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var c Config
	if err := yaml.Unmarshal(b, &c); err != nil {
		return nil, err
	}
	return &c, nil
}

func (c *Config) ByIDs(ids []string) []Item {
	set := map[string]struct{}{}
	for _, id := range ids {
		set[id] = struct{}{}
	}
	var out []Item
	for _, it := range c.Items {
		if _, ok := set[it.ID]; ok {
			out = append(out, it)
		}
	}
	return out
}

func (c *Config) Curated() []Item {
	if len(c.Curate) == 0 {
		return c.Items
	}
	return c.ByIDs(c.Curate)
}

func (i Item) Validate() error {
	if i.ID == "" || i.Name == "" {
		return errors.New("invalid item: id/name required")
	}
	return nil
}

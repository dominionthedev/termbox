package registry

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"gopkg.in/yaml.v3"
)

var mu sync.Mutex

// ConfigKind defines how a config entry is handled by termbox.
//
//   main     — termbox manages this config. "termbox use app" copies it to Target.
//   template — a named variant. "termbox use app.name" swaps it in, backs up current.
//   addition — termbox is sourced from the user's config. No file operations.
//              "termbox setup" prints the source lines to add manually.
type ConfigKind string

const (
	KindMain     ConfigKind = "main"
	KindTemplate ConfigKind = "template"
	KindAddition ConfigKind = "addition"
)

// Item is a single component registered in the registry.
type Item struct {
	Name        string     `yaml:"name"`
	Kind        string     `yaml:"kind"`                  // tool | script | config
	ConfigKind  ConfigKind `yaml:"config_kind,omitempty"` // main | template | addition
	App         string     `yaml:"app,omitempty"`         // owning app: nvim | starship | tmux | zsh ...
	Path        string     `yaml:"path"`                  // relative to TERMBOX_HOME
	Target      string     `yaml:"target,omitempty"`      // copy destination (main/template only, no symlinks)
	Description string     `yaml:"description"`
	Active      bool       `yaml:"active,omitempty"`
}

// Registry is the full component index.
type Registry struct {
	Tools   []Item `yaml:"tools"`
	Scripts []Item `yaml:"scripts"`
	Configs []Item `yaml:"configs"`
}

var ErrDuplicateName = errors.New("item with this name already exists")

// FindHome resolves TERMBOX_HOME in order:
//  1. $TERMBOX_HOME env var — set by wizard in termbox.env, most reliable
//  2. Parent of the executable's bin/ directory
//  3. Current working directory — least reliable fallback
func FindHome() (string, error) {
	if h := os.Getenv("TERMBOX_HOME"); h != "" {
		clean := filepath.Clean(h)
		if _, err := os.Stat(filepath.Join(clean, "config", "registry.yaml")); err == nil {
			return clean, nil
		}
		return "", fmt.Errorf("TERMBOX_HOME=%q is set but config/registry.yaml not found there", h)
	}

	if exe, err := os.Executable(); err == nil {
		root := filepath.Dir(filepath.Dir(filepath.Clean(exe)))
		if _, err := os.Stat(filepath.Join(root, "config", "registry.yaml")); err == nil {
			return root, nil
		}
	}

	cwd, err := os.Getwd()
	if err != nil {
		return "", fmt.Errorf("cannot determine working directory: %w", err)
	}
	if _, err := os.Stat(filepath.Join(cwd, "config", "registry.yaml")); err == nil {
		return cwd, nil
	}

	return "", fmt.Errorf(
		"TERMBOX_HOME not found — run 'termbox setup --wizard' and follow the instructions",
	)
}

func LoadRegistry(customPath string) (*Registry, error) {
	mu.Lock()
	defer mu.Unlock()
	path, err := resolvePath(customPath)
	if err != nil {
		return nil, err
	}
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("reading registry %q: %w", path, err)
	}
	var reg Registry
	if err := yaml.Unmarshal(data, &reg); err != nil {
		return nil, fmt.Errorf("parsing registry: %w", err)
	}
	return &reg, nil
}

func SaveRegistry(reg *Registry, customPath string) error {
	mu.Lock()
	defer mu.Unlock()
	path, err := resolvePath(customPath)
	if err != nil {
		return err
	}
	data, err := yaml.Marshal(reg)
	if err != nil {
		return fmt.Errorf("serialising registry: %w", err)
	}
	tmp := path + ".tmp"
	if err := os.WriteFile(tmp, data, 0644); err != nil {
		return fmt.Errorf("writing temp: %w", err)
	}
	if err := os.Rename(tmp, path); err != nil {
		_ = os.Remove(tmp)
		return fmt.Errorf("replacing registry: %w", err)
	}
	return nil
}

// FindItem returns a pointer to the actual slice element so mutations stick.
func (r *Registry) FindItem(name string) *Item {
	for i := range r.Tools {
		if r.Tools[i].Name == name {
			return &r.Tools[i]
		}
	}
	for i := range r.Scripts {
		if r.Scripts[i].Name == name {
			return &r.Scripts[i]
		}
	}
	for i := range r.Configs {
		if r.Configs[i].Name == name {
			return &r.Configs[i]
		}
	}
	return nil
}

// ConfigsForApp returns all configs whose App field matches.
func (r *Registry) ConfigsForApp(app string) []Item {
	var out []Item
	for _, c := range r.Configs {
		if strings.EqualFold(c.App, app) {
			out = append(out, c)
		}
	}
	return out
}

func (r *Registry) ValidateNewItem(item Item) error {
	if strings.TrimSpace(item.Name) == "" {
		return errors.New("name cannot be empty")
	}
	if strings.TrimSpace(item.Path) == "" {
		return errors.New("path cannot be empty")
	}
	if r.FindItem(item.Name) != nil {
		return fmt.Errorf("%w: %q", ErrDuplicateName, item.Name)
	}
	return nil
}

// ExpandHome expands a leading ~/ in a path.
func ExpandHome(path string) string {
	if strings.HasPrefix(path, "~/") {
		if home, err := os.UserHomeDir(); err == nil {
			return filepath.Join(home, path[2:])
		}
	}
	return path
}

func resolvePath(customPath string) (string, error) {
	if customPath != "" {
		return filepath.Clean(customPath), nil
	}
	home, err := FindHome()
	if err != nil {
		return "", err
	}
	return filepath.Join(home, "config", "registry.yaml"), nil
}

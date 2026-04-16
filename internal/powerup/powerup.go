package powerup

import (
	"os"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v3"
)

// Rules defines when a powerup should be activated.
type Rules struct {
	Always   bool     `yaml:"always"`
	OS       string   `yaml:"os,omitempty"`
	Requires []string `yaml:"requires,omitempty"` // binaries that must be in PATH
	Env      []EnvRule `yaml:"env,omitempty"`     // environment criteria
}

// EnvRule defines an environment-based activation criterion.
// Kind is one of: "file_exists", "dir_exists", "env_set", "env_value".
// Value is the path, env var name, or expected value.
//
// Examples:
//   kind: file_exists  value: go.mod        → activate if ./go.mod exists
//   kind: file_exists  value: Cargo.toml    → activate if ./Cargo.toml exists
//   kind: dir_exists   value: .git          → activate if in a git repo
//   kind: env_set      value: GOPATH        → activate if $GOPATH is set
//   kind: env_value    value: TERMBOX_POWERUP=rust → activate if env matches
type EnvRule struct {
	Kind  string `yaml:"kind"`
	Value string `yaml:"value"`
}

// Powerup represents a loaded powerup definition.
type Powerup struct {
	Name        string   `yaml:"name"`
	Description string   `yaml:"description"`
	Version     string   `yaml:"version"`
	Scripts     []string `yaml:"scripts"`
	Tools       []string `yaml:"tools"`
	Rules       Rules    `yaml:"rules"`
}

// Load reads a powerup YAML file by name from the powerups/ directory.
func Load(tbHome, name string) (*Powerup, error) {
	path := filepath.Join(tbHome, "powerups", name+".yaml")
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var p Powerup
	if err := yaml.Unmarshal(data, &p); err != nil {
		return nil, err
	}
	return &p, nil
}

// LoadAll reads every *.yaml file in the powerups/ directory.
func LoadAll(tbHome string) ([]*Powerup, error) {
	dir := filepath.Join(tbHome, "powerups")
	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
	}
	var out []*Powerup
	for _, e := range entries {
		if e.IsDir() || !strings.HasSuffix(e.Name(), ".yaml") {
			continue
		}
		name := strings.TrimSuffix(e.Name(), ".yaml")
		p, err := Load(tbHome, name)
		if err != nil {
			continue
		}
		out = append(out, p)
	}
	return out, nil
}

// MeetsEnvCriteria returns true if all env rules in p.Rules.Env pass
// for the current working directory and process environment.
func (p *Powerup) MeetsEnvCriteria() bool {
	cwd, _ := os.Getwd()
	for _, rule := range p.Rules.Env {
		if !evalEnvRule(rule, cwd) {
			return false
		}
	}
	return true
}

// MeetsRequires returns true if every binary listed in p.Rules.Requires
// is available in $PATH.
func (p *Powerup) MeetsRequires() (bool, string) {
	for _, bin := range p.Rules.Requires {
		if _, err := lookPath(bin); err != nil {
			return false, bin
		}
	}
	return true, ""
}

// ShouldAutoActivate returns true if this powerup should activate automatically
// given the current environment (used when settings.powerup.auto_detect = true).
func (p *Powerup) ShouldAutoActivate() bool {
	if p.Rules.Always {
		return true
	}
	ok, _ := p.MeetsRequires()
	return ok && p.MeetsEnvCriteria()
}

// evalEnvRule evaluates a single EnvRule against the current environment.
func evalEnvRule(rule EnvRule, cwd string) bool {
	switch rule.Kind {
	case "file_exists":
		_, err := os.Stat(filepath.Join(cwd, rule.Value))
		return err == nil
	case "dir_exists":
		info, err := os.Stat(filepath.Join(cwd, rule.Value))
		return err == nil && info.IsDir()
	case "env_set":
		return os.Getenv(rule.Value) != ""
	case "env_value":
		// format: "KEY=VALUE"
		if i := strings.IndexByte(rule.Value, '='); i > 0 {
			key := rule.Value[:i]
			val := rule.Value[i+1:]
			return os.Getenv(key) == val
		}
		return false
	}
	return false
}

// lookPath is a thin wrapper so tests can override it.
var lookPath = func(name string) (string, error) {
	return findInPath(name)
}

func findInPath(name string) (string, error) {
	pathEnv := os.Getenv("PATH")
	for _, dir := range filepath.SplitList(pathEnv) {
		full := filepath.Join(dir, name)
		if info, err := os.Stat(full); err == nil && !info.IsDir() {
			return full, nil
		}
	}
	return "", os.ErrNotExist
}

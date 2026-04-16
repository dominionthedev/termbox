package settings

import (
	"fmt"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

// Settings represents the full termbox settings file (config/settings.yaml).
// This is the application config — separate from the component registry.
type Settings struct {
	Display struct {
		Color bool `yaml:"color"`
		Icons bool `yaml:"icons"`
		Pager bool `yaml:"pager"`
	} `yaml:"display"`

	Banner struct {
		Enabled bool   `yaml:"enabled"`
		File    string `yaml:"file"`
		Color   bool   `yaml:"color"`
	} `yaml:"banner"`

	Notes struct {
		Folder    string `yaml:"folder"`
		Extension string `yaml:"extension"`
		Editor    string `yaml:"editor"`
	} `yaml:"notes"`

	Theme struct {
		Active    string `yaml:"active"`
		AutoApply bool   `yaml:"auto_apply"`
	} `yaml:"theme"`

	Powerup struct {
		Active     string `yaml:"active"`
		AutoDetect bool   `yaml:"auto_detect"`
	} `yaml:"powerup"`

	Build struct {
		OutputDir string `yaml:"output_dir"`
		GoFlags   string `yaml:"go_flags"`
	} `yaml:"build"`
}

// Defaults returns a Settings struct populated with sensible defaults.
func Defaults() *Settings {
	s := &Settings{}
	s.Display.Color = true
	s.Display.Icons = true
	s.Banner.Enabled = true
	s.Banner.Color = true
	s.Notes.Extension = ".md"
	s.Theme.Active = "catppuccin_mocha"
	s.Theme.AutoApply = true
	s.Powerup.Active = "core"
	s.Powerup.AutoDetect = true
	s.Build.OutputDir = "bin"
	return s
}

// Load reads settings.yaml from the termbox home directory.
// If the file does not exist, Defaults() is returned without error.
func Load(tbHome string) (*Settings, error) {
	path := filepath.Join(tbHome, "config", "settings.yaml")

	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return Defaults(), nil
		}
		return nil, fmt.Errorf("reading settings %q: %w", path, err)
	}

	s := Defaults() // start from defaults so unset keys keep sensible values
	if err := yaml.Unmarshal(data, s); err != nil {
		return nil, fmt.Errorf("parsing settings: %w", err)
	}
	return s, nil
}

// Save writes settings back to settings.yaml.
func Save(s *Settings, tbHome string) error {
	path := filepath.Join(tbHome, "config", "settings.yaml")

	data, err := yaml.Marshal(s)
	if err != nil {
		return fmt.Errorf("serialising settings: %w", err)
	}

	tmp := path + ".tmp"
	if err := os.WriteFile(tmp, data, 0644); err != nil {
		return fmt.Errorf("writing settings temp: %w", err)
	}
	if err := os.Rename(tmp, path); err != nil {
		_ = os.Remove(tmp)
		return fmt.Errorf("replacing settings: %w", err)
	}
	return nil
}

// NoteFolder resolves the effective notes directory, checking settings then env.
func (s *Settings) NoteFolder() string {
	if s.Notes.Folder != "" {
		return s.Notes.Folder
	}
	if env := os.Getenv("NOTE_FOLDER"); env != "" {
		return env
	}
	home, _ := os.UserHomeDir()
	return filepath.Join(home, "Developer", "notes")
}

// BannerFile resolves the effective banner file path.
func (s *Settings) BannerFile(tbHome string) string {
	if s.Banner.File != "" {
		return s.Banner.File
	}
	if env := os.Getenv("BANNER"); env != "" {
		return env
	}
	return filepath.Join(tbHome, "assets", "dominiondev.banner")
}

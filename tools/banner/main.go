// banner — standalone banner tool
// Built into bin/banner by `termbox build`.
// Reads BANNER and TERMBOX_SHOW_BANNER from the environment (set in termbox.env).
// To change the banner file, edit BANNER in config/termbox.env directly,
// or use: termbox config set banner.file <path>
//
// Usage:
//   banner           show the configured banner
//   banner show      same as above
//   banner path      print the resolved banner file path
package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
)

func main() {
	if err := run(os.Args[1:]); err != nil {
		fmt.Fprintf(os.Stderr, "banner: %v\n", err)
		os.Exit(1)
	}
}

func run(args []string) error {
	action := "show"
	if len(args) > 0 {
		action = args[0]
	}

	tbHome := os.Getenv("TERMBOX_HOME")
	if tbHome == "" {
		tbHome = filepath.Join(os.Getenv("HOME"), "Developer", "termbox")
	}

	switch action {
	case "show", "":
		return showBanner(tbHome)
	case "path":
		fmt.Println(resolveBannerFile(tbHome))
		return nil
	default:
		return fmt.Errorf("unknown action %q — use: show, path", action)
	}
}

func resolveBannerFile(tbHome string) string {
	if f := os.Getenv("BANNER"); f != "" {
		return f
	}
	return filepath.Join(tbHome, "assets", "dominiondev.banner")
}

func showBanner(tbHome string) error {
	bannerFile := resolveBannerFile(tbHome)
	data, err := os.ReadFile(bannerFile)
	if err != nil {
		return fmt.Errorf("banner file not found at %q\nSet BANNER in config/termbox.env", bannerFile)
	}

	if _, err := exec.LookPath("lolcat"); err == nil {
		r, w, _ := os.Pipe()
		go func() { _, _ = w.Write(data); w.Close() }()
		c := exec.Command("lolcat")
		c.Stdin = r
		c.Stdout = os.Stdout
		c.Stderr = os.Stderr
		if c.Run() == nil {
			return nil
		}
	}

	fmt.Print(string(data))
	return nil
}

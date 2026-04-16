# Termbox

Your terminal environment, organised.

Termbox manages your tools, scripts, configs, and powerups.
It is a guest in your shell — it never touches your dotfiles uninvited.

---

## Setup

**1. Build termbox:**
```sh
cd $TERMBOX_HOME
go get github.com/leraniode/wondertone@v0.2.0
go mod tidy
go build -o bin/termbox ./cmd/termbox
```

**2. Run the wizard:**
```sh
./bin/termbox setup --wizard
```

**3. Add two lines to the top of `~/.zshrc`:**
```sh
# ---- termbox env ----
source $HOME/Developer/termbox/config/termbox.env
# ---- termbox zsh ----
source $TERMBOX_HOME/config/shell/zshrc
```

**4. Reload:**
```sh
source ~/.zshrc
```

---

## Commands

| Command | Description |
|---|---|
| `termbox setup --wizard` | First-time setup wizard |
| `termbox setup` | Print the two source lines |
| `termbox list` | List all registered components |
| `termbox run <script>` | Run a registered script |
| `termbox use <app>` | List configs for an app |
| `termbox use nvim.default` | Apply a config template (copies, never symlinks) |
| `termbox note <name>` | Open/create a note |
| `termbox note list` | List all notes |
| `termbox banner show` | Show the startup banner |
| `termbox banner set <file>` | Change the banner file |
| `termbox sheme list` | List installed themes |
| `termbox sheme palettes` | List Wondertone built-in palettes |
| `termbox sheme from <palette>` | Generate a .theme from a Wondertone palette |
| `termbox sheme apply <name>` | Set the active theme |
| `termbox powerup list` | List powerups |
| `termbox powerup activate <name>` | Activate a powerup |
| `termbox doctor` | Health check |
| `termbox status` | Show runtime env vars |

---

## Config kinds

| Kind | How it works |
|---|---|
| `addition` | Sourced from your shell. `termbox setup` prints the lines to add. |
| `main` | Copied to target on `termbox use app`. |
| `template` | Named variant. `termbox use app.name` copies it in, backs up current. |

---

## Powerups

Powerups are purpose-specific packs of scripts. Activate one:
```sh
termbox powerup activate container
```

Built-in: `core`, `container`, `git`.

---

## Themes

Themes are Wondertone-powered `.theme` files in `assets/themes/`.
They export `THEME_*` shell variables and emit OSC sequences to colour your terminal.

```sh
termbox sheme from "Catppuccin Mocha"   # generate from Wondertone palette
termbox sheme apply catppuccin_mocha    # activate it
```

Shipped themes: `catppuccin_mocha`, `cybergreen`, `leraniode`.

#!/usr/bin/env bash
# Termbox quick diagnostics — mirrors ShellDoctor but as a standalone script

PASS="  ✓"
FAIL="  ✗"
WARN="  ⚠"

echo ""
echo "  ┌──────────────────────────────────┐"
echo "  │  Termbox Diagnostics             │"
echo "  └──────────────────────────────────┘"
echo ""

echo "  System"
echo "    User   : $USER"
echo "    Shell  : $SHELL"
echo "    OS     : $(uname -s) $(uname -m)"
echo ""

echo "  Termbox"
if [ -n "$TERMBOX_ROOT" ]; then
    echo "$PASS TERMBOX_ROOT = $TERMBOX_ROOT"
else
    echo "$FAIL TERMBOX_ROOT is not set — run: termbox setup --wizard"
fi

registry="$TERMBOX_ROOT/config/registry.yaml"
if [ -f "$registry" ]; then
    echo "$PASS registry.yaml found"
else
    echo "$FAIL registry.yaml missing at $registry"
fi

theme_file="$TERMBOX_ROOT/config/active-theme"
if [ -f "$theme_file" ]; then
    echo "$PASS Theme: $(cat "$theme_file")"
else
    echo "$WARN No active theme set — run: termbox theme set <n>"
fi
echo ""

echo "  Core Tools"
for tool in tmux nvim starship fzf zoxide bat eza fd rg; do
    if command -v "$tool" >/dev/null 2>&1; then
        echo "$PASS $tool"
    else
        echo "$FAIL $tool (not found)"
    fi
done
echo ""

echo "  Optional Tools"
for tool in lolcat tokei docker kubectl htop procs dust delta; do
    if command -v "$tool" >/dev/null 2>&1; then
        echo "$PASS $tool"
    else
        echo "$WARN $tool (not installed — optional)"
    fi
done
echo ""

echo "  Shell Config"
zshrc="$HOME/.zshrc"
if [ -L "$zshrc" ]; then
    target=$(readlink "$zshrc")
    echo "$PASS ~/.zshrc is a symlink → $target"
elif [ -f "$zshrc" ]; then
    echo "$WARN ~/.zshrc is a regular file (not symlinked from termbox)"
else
    echo "$FAIL ~/.zshrc not found"
fi
echo ""

# Startup time (zsh only)
if command -v zsh >/dev/null 2>&1; then
    startup_ms=$(( $(zsh -i -c 'exit' 2>/dev/null; echo $((SECONDS * 1000))) ))
    # Use a different approach
    start=$(date +%s%3N 2>/dev/null || echo 0)
    zsh -i -c 'exit' 2>/dev/null
    end=$(date +%s%3N 2>/dev/null || echo 0)
    if [ "$start" != "0" ]; then
        elapsed=$((end - start))
        if [ "$elapsed" -gt 500 ]; then
            echo "$WARN Shell startup: ~${elapsed}ms (slow — check plugins)"
        else
            echo "$PASS Shell startup: ~${elapsed}ms"
        fi
    fi
fi
echo ""

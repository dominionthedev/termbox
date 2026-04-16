#!/usr/bin/env bash
# setup.sh — Termbox pre-setup script
#
# Run this BEFORE `termbox setup --wizard` to detect where termbox is installed
# and bootstrap config/termbox.env with the correct TERMBOX_HOME.
#
# Usage:
#   bash setup.sh            auto-detect and write termbox.env
#   bash setup.sh --check    just print detected path, don't write anything
#   bash setup.sh --help     show this message

set -euo pipefail

# ── Helpers ───────────────────────────────────────────────────────────────────
RED='\033[0;31m'; GREEN='\033[0;32m'; YELLOW='\033[1;33m'
CYAN='\033[0;36m'; BOLD='\033[1m'; RESET='\033[0m'

info()    { printf "  ${CYAN}→${RESET}  %s\n" "$*"; }
ok()      { printf "  ${GREEN}✓${RESET}  %s\n" "$*"; }
warn()    { printf "  ${YELLOW}⚠${RESET}  %s\n" "$*"; }
fail()    { printf "  ${RED}✗${RESET}  %s\n" "$*"; exit 1; }

# ── Detect termbox home ───────────────────────────────────────────────────────
detect_home() {
    # 1. Already set in environment
    if [[ -n "${TERMBOX_HOME:-}" ]]; then
        echo "$TERMBOX_HOME"
        return
    fi

    # 2. Resolve from this script's own location
    local script_dir
    script_dir="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
    if [[ -f "$script_dir/config/registry.yaml" ]]; then
        echo "$script_dir"
        return
    fi

    # 3. Common install locations
    local candidates=(
        "$HOME/Developer/termbox"
        "$HOME/.termbox"
        "$HOME/.config/termbox"
        "/usr/local/share/termbox"
    )
    for candidate in "${candidates[@]}"; do
        if [[ -f "$candidate/config/registry.yaml" ]]; then
            echo "$candidate"
            return
        fi
    done

    # 4. Give up
    echo ""
}

# ── Main ─────────────────────────────────────────────────────────────────────
main() {
    local check_only=false

    for arg in "$@"; do
        case "$arg" in
            --check)  check_only=true ;;
            --help|-h)
                echo "Usage: bash setup.sh [--check] [--help]"
                echo ""
                echo "  --check   show detected path without writing anything"
                echo "  --help    show this message"
                exit 0
                ;;
        esac
    done

    printf "\n${BOLD}  Termbox Pre-Setup${RESET}\n\n"

    local tb_home
    tb_home="$(detect_home)"

    if [[ -z "$tb_home" ]]; then
        fail "Could not detect termbox home directory."
        printf "\n  Set TERMBOX_HOME manually and re-run:\n"
        printf "    export TERMBOX_HOME=/path/to/termbox\n    bash setup.sh\n\n"
        exit 1
    fi

    ok "Detected TERMBOX_HOME = $tb_home"

    if $check_only; then
        printf "\n  (--check mode: nothing written)\n\n"
        exit 0
    fi

    local env_file="$tb_home/config/termbox.env"

    # If termbox.env already has TERMBOX_HOME set to the right value, skip
    if [[ -f "$env_file" ]] && grep -q "TERMBOX_HOME=\"$tb_home\"" "$env_file" 2>/dev/null; then
        ok "termbox.env already configured correctly"
    else
        # Patch or create termbox.env
        if [[ -f "$env_file" ]]; then
            info "Patching TERMBOX_HOME in existing termbox.env..."
            local tmp="$env_file.tmp"
            # Replace any existing TERMBOX_HOME line
            sed "s|^export TERMBOX_HOME=.*|export TERMBOX_HOME=\"$tb_home\"|" "$env_file" > "$tmp"
            # If no line was replaced, add it after the header comment block
            if ! grep -q "TERMBOX_HOME=" "$tmp"; then
                echo "export TERMBOX_HOME=\"$tb_home\"" >> "$tmp"
            fi
            mv "$tmp" "$env_file"
        else
            info "Creating minimal termbox.env..."
            mkdir -p "$(dirname "$env_file")"
            cat > "$env_file" << ENVEOF
# termbox.env — created by setup.sh
# Run 'termbox setup --wizard' to fill in the full configuration.

export TERMBOX_HOME="$tb_home"
export TERMBOX_DEV_FOLDER="\$HOME/Developer"
export NOTE_FOLDER="\$HOME/Developer/notes"
export BANNER="\$TERMBOX_HOME/assets/dominiondev.banner"
export TERMBOX_SHOW_BANNER="true"
export SHELL_THEME="catppuccin_mocha"
export TERMBOX_DEFAULT_NVIM="default"
export TERMBOX_POWERUP="core"
ENVEOF
        fi
        ok "Written: $env_file"
    fi

    # Build the termbox binary if not already present
    local binary="$tb_home/bin/termbox"
    if [[ -f "$binary" ]]; then
        ok "Binary already exists: $binary"
    else
        if command -v go &>/dev/null; then
            info "Building termbox binary..."
            if (cd "$tb_home" && go build -o "$binary" ./cmd/termbox 2>/dev/null); then
                ok "Built: $binary"
            else
                warn "Build failed — run manually: cd $tb_home && go build -o bin/termbox ./cmd/termbox"
            fi
        else
            warn "Go not found — build manually: cd $tb_home && go build -o bin/termbox ./cmd/termbox"
        fi
    fi

    printf "\n${BOLD}  Next steps:${RESET}\n\n"
    printf "  1. Add to the TOP of ~/.zshrc:\n\n"
    printf "       ${CYAN}# ---- termbox env ----${RESET}\n"
    printf "       ${CYAN}source %s${RESET}\n" "$env_file"
    printf "       ${CYAN}# ---- termbox zsh ----${RESET}\n"
    printf "       ${CYAN}source \$TERMBOX_HOME/config/shell/zshrc${RESET}\n\n"
    printf "  2. Run the full wizard for complete setup:\n\n"
    printf "       ${CYAN}termbox setup --wizard${RESET}\n\n"
}

main "$@"

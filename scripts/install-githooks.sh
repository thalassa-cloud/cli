#!/usr/bin/env bash
set -euo pipefail

repo_root="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"

git -C "${repo_root}" config core.hooksPath .githooks
chmod +x "${repo_root}/.githooks/pre-commit"

echo "Installed git hooks from ${repo_root}/.githooks"

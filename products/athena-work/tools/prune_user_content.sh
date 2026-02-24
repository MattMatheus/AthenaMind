#!/usr/bin/env bash
set -euo pipefail

script_dir="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
root_dir="$(git -C "$script_dir" rev-parse --show-toplevel 2>/dev/null || (cd "$script_dir/.." && pwd))"
system_json_default="$root_dir/operating-system/system.json"
system_json="$system_json_default"
export_zip=""
dry_run=false
auto_confirm=false

usage() {
  cat <<'USAGE'
usage: tools/prune_user_content.sh [--yes] [--dry-run] [--export-zip <zip_path>] [--system-json <path>]

Destructive operation:
- Removes user-owned research/idea content based on operating-system/system.json
- Recreates sanitized placeholder files

Options:
  --yes                 Skip interactive confirmation prompt
  --dry-run             Show actions only, do not modify files
  --export-zip <path>   Export prune-target files to a zip before deletion
  --system-json <path>  Override system layout json path
  -h, --help            Show this help
USAGE
}

require_cmd() {
  if ! command -v "$1" >/dev/null 2>&1; then
    echo "error: required command not found: $1" >&2
    exit 1
  fi
}

while [[ $# -gt 0 ]]; do
  case "$1" in
    --yes)
      auto_confirm=true
      shift
      ;;
    --dry-run)
      dry_run=true
      shift
      ;;
    --export-zip)
      export_zip="${2:-}"
      if [[ -z "$export_zip" ]]; then
        echo "error: --export-zip requires a path" >&2
        exit 1
      fi
      shift 2
      ;;
    --system-json)
      system_json="${2:-}"
      if [[ -z "$system_json" ]]; then
        echo "error: --system-json requires a path" >&2
        exit 1
      fi
      shift 2
      ;;
    -h|--help)
      usage
      exit 0
      ;;
    *)
      echo "error: unknown argument '$1'" >&2
      usage
      exit 1
      ;;
  esac
done

require_cmd git
require_cmd python3

if [[ ! -f "$system_json" ]]; then
  echo "error: system json not found: $system_json" >&2
  exit 1
fi

if ! git -C "$root_dir" rev-parse --is-inside-work-tree >/dev/null 2>&1; then
  echo "error: not inside a git repository: $root_dir" >&2
  exit 1
fi

delete_paths=()
while IFS= read -r line; do
  [[ -n "$line" ]] && delete_paths+=("$line")
done < <(python3 - "$system_json" <<'PY'
import json, sys
with open(sys.argv[1], "r", encoding="utf-8") as f:
    data = json.load(f)
for p in data["agent_section"]["user_content_prune"]["delete_paths"]:
    print(p)
PY
)

delete_globs=()
while IFS= read -r line; do
  [[ -n "$line" ]] && delete_globs+=("$line")
done < <(python3 - "$system_json" <<'PY'
import json, sys
with open(sys.argv[1], "r", encoding="utf-8") as f:
    data = json.load(f)
for p in data["agent_section"]["user_content_prune"]["delete_globs"]:
    print(p)
PY
)

confirm_token="$(python3 - "$system_json" <<'PY'
import json, sys
with open(sys.argv[1], "r", encoding="utf-8") as f:
    data = json.load(f)
print(data["agent_section"]["user_content_prune"]["confirmation_token"])
PY
)"

targets=()
while IFS= read -r line; do
  [[ -n "$line" ]] && targets+=("$line")
done < <(
  (
    for p in "${delete_paths[@]}"; do
      echo "$p"
    done
    for g in "${delete_globs[@]}"; do
      git -C "$root_dir" ls-files "$g"
    done
  ) | sed '/^$/d' | awk '!seen[$0]++' | sort
)

if [[ "${#targets[@]}" -eq 0 ]]; then
  echo "No matching prune targets found."
  exit 0
fi

echo "Prune targets (${#targets[@]}):"
for t in "${targets[@]}"; do
  echo " - $t"
done

if [[ "$auto_confirm" != true ]]; then
  echo
  echo "This operation is destructive."
  echo "Type '$confirm_token' to continue:"
  read -r typed
  if [[ "$typed" != "$confirm_token" ]]; then
    echo "abort: confirmation token mismatch"
    exit 1
  fi
fi

if [[ -n "$export_zip" ]]; then
  if [[ "$export_zip" != /* ]]; then
    export_zip="$root_dir/$export_zip"
  fi
  export_dir="$(dirname "$export_zip")"
  if [[ "$dry_run" == true ]]; then
    echo "DRY-RUN: would export zip to $export_zip"
  else
    mkdir -p "$export_dir"
    if [[ ! -w "$export_dir" ]]; then
      echo "error: export directory is not writable: $export_dir" >&2
      echo "hint: re-run with a writable path or with elevated permissions." >&2
      exit 1
    fi
    tmp_manifest="$(mktemp)"
    trap 'rm -f "$tmp_manifest"' EXIT
    for t in "${targets[@]}"; do
      printf '%s\n' "$t" >> "$tmp_manifest"
    done
    python3 - "$root_dir" "$export_zip" "$tmp_manifest" <<'PY'
import os, sys, zipfile
root, out_zip, manifest = sys.argv[1:4]
with zipfile.ZipFile(out_zip, "w", compression=zipfile.ZIP_DEFLATED) as zf:
    with open(manifest, "r", encoding="utf-8") as f:
        for line in f:
            rel = line.strip()
            if not rel:
                continue
            abs_path = os.path.join(root, rel)
            if os.path.exists(abs_path):
                zf.write(abs_path, rel)
print(out_zip)
PY
    echo "exported: $export_zip"
  fi
fi

if [[ "$dry_run" == true ]]; then
  echo "DRY-RUN: would delete target files and rewrite placeholders."
  exit 0
fi

for t in "${targets[@]}"; do
  git -C "$root_dir" rm -f --ignore-unmatch "$t" >/dev/null 2>&1 || true
  if [[ -e "$root_dir/$t" ]]; then
    rm -rf "$root_dir/$t"
  fi
done

python3 - "$system_json" "$root_dir" <<'PY'
import json, os, sys
system_json, root = sys.argv[1], sys.argv[2]
with open(system_json, "r", encoding="utf-8") as f:
    data = json.load(f)
for item in data["agent_section"]["user_content_prune"]["placeholders"]:
    rel = item["path"]
    content = item["content"]
    out = os.path.join(root, rel)
    os.makedirs(os.path.dirname(out), exist_ok=True)
    with open(out, "w", encoding="utf-8", newline="\n") as f:
      f.write(content)
PY

echo "Prune complete. Review and commit changes."

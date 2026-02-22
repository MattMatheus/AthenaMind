#!/usr/bin/env bash
set -euo pipefail

root_dir="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
eng_intake="$root_dir/backlog/engineering/intake"
arch_intake="$root_dir/backlog/architecture/intake"

valid_story_status='intake|ready|active|qa|done|blocked'
valid_bug_status='intake|ready|active|qa|done|blocked'
valid_bug_priority='P0|P1|P2|P3'

failures=0

fail() {
  echo "FAIL: $1" >&2
  failures=$((failures + 1))
}

require_line() {
  local file="$1"
  local pattern="$2"
  local message="$3"
  if ! rg -q -- "$pattern" "$file"; then
    fail "$message ($file)"
  fi
}

for file in "$eng_intake"/*.md; do
  [ -e "$file" ] || continue
  base="$(basename "$file")"
  case "$base" in
    STORY_TEMPLATE.md|BUG_TEMPLATE.md)
      continue
      ;;
  esac

  if [[ "$base" == STORY-* ]]; then
    require_line "$file" '^## Metadata$' 'missing metadata header'
    require_line "$file" '^[[:space:]]*-[[:space:]]*`id`:[[:space:]]*STORY-' 'engineering story id must start with STORY-'
    require_line "$file" '^[[:space:]]*-[[:space:]]*`owner_persona`:[[:space:]]*' 'missing owner_persona'
    require_line "$file" "^[[:space:]]*-[[:space:]]*\`status\`:[[:space:]]*($valid_story_status)$" 'story status invalid'
    require_line "$file" '^[[:space:]]*-[[:space:]]*`idea_id`:[[:space:]]*' 'missing idea_id'
    require_line "$file" '^[[:space:]]*-[[:space:]]*`phase`:[[:space:]]*(v0.1|v0.2|v0.3)$' 'missing or invalid phase'
    require_line "$file" '^[[:space:]]*-[[:space:]]*`adr_refs`:[[:space:]]*' 'missing adr_refs'
    require_line "$file" '^[[:space:]]*-[[:space:]]*`success_metric`:[[:space:]]*' 'missing success_metric'
    require_line "$file" '^[[:space:]]*-[[:space:]]*`release_checkpoint`:[[:space:]]*(required|deferred)$' 'missing or invalid release_checkpoint'
  elif [[ "$base" == BUG-* ]]; then
    require_line "$file" '^## Metadata$' 'missing metadata header'
    require_line "$file" '^[[:space:]]*-[[:space:]]*`id`:[[:space:]]*BUG-' 'engineering bug id must start with BUG-'
    require_line "$file" "^[[:space:]]*-[[:space:]]*\`priority\`:[[:space:]]*($valid_bug_priority)$" 'bug priority invalid'
    require_line "$file" '^[[:space:]]*-[[:space:]]*`reported_by`:[[:space:]]*' 'missing reported_by'
    require_line "$file" '^[[:space:]]*-[[:space:]]*`source_story`:[[:space:]]*' 'missing source_story'
    require_line "$file" "^[[:space:]]*-[[:space:]]*\`status\`:[[:space:]]*($valid_bug_status)$" 'bug status invalid'
    require_line "$file" '^[[:space:]]*-[[:space:]]*`phase`:[[:space:]]*(v0.1|v0.2|v0.3)$' 'missing or invalid phase'
    require_line "$file" '^[[:space:]]*-[[:space:]]*`adr_refs`:[[:space:]]*' 'missing adr_refs'
    require_line "$file" '^[[:space:]]*-[[:space:]]*`impact_metric`:[[:space:]]*' 'missing impact_metric'
  else
    fail "unexpected file type in engineering intake: $file"
  fi

  if rg -q '^[[:space:]]*-[[:space:]]*`id`:[[:space:]]*ARCH-' "$file"; then
    fail "architecture id found in engineering intake: $file"
  fi
done

for file in "$arch_intake"/*.md; do
  [ -e "$file" ] || continue
  base="$(basename "$file")"
  case "$base" in
    ARCH_STORY_TEMPLATE.md)
      continue
      ;;
  esac

  require_line "$file" '^## Metadata$' 'missing metadata header'
  require_line "$file" '^[[:space:]]*-[[:space:]]*`id`:[[:space:]]*ARCH-' 'architecture story id must start with ARCH-'
  require_line "$file" '^[[:space:]]*-[[:space:]]*`owner_persona`:[[:space:]]*' 'missing owner_persona'
  require_line "$file" '^[[:space:]]*-[[:space:]]*`status`:[[:space:]]*(intake|ready|active|qa|done|blocked)$' 'architecture status invalid'
  require_line "$file" '^[[:space:]]*-[[:space:]]*`idea_id`:[[:space:]]*' 'missing idea_id'
  require_line "$file" '^[[:space:]]*-[[:space:]]*`phase`:[[:space:]]*(v0.1|v0.2|v0.3)$' 'missing or invalid phase'
  require_line "$file" '^[[:space:]]*-[[:space:]]*`adr_refs`:[[:space:]]*' 'missing adr_refs'
  require_line "$file" '^[[:space:]]*-[[:space:]]*`decision_owner`:[[:space:]]*' 'missing decision_owner'
  require_line "$file" '^[[:space:]]*-[[:space:]]*`success_metric`:[[:space:]]*' 'missing success_metric'

  if rg -q '^[[:space:]]*-[[:space:]]*`id`:[[:space:]]*(STORY-|BUG-)' "$file"; then
    fail "engineering id found in architecture intake: $file"
  fi
done

if [ "$failures" -ne 0 ]; then
  exit 1
fi

echo "PASS: intake metadata and lane-separation validation"

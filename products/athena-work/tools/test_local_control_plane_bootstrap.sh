#!/usr/bin/env bash
set -euo pipefail

script_dir="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
root_dir="$(git -C "$script_dir" rev-parse --show-toplevel 2>/dev/null || (cd "$script_dir/.." && pwd))"
source "$root_dir/tools/lib/doc_test_harness.sh"

compose_file="$root_dir/docker-compose.local.yml"
quickstart="$root_dir/knowledge-base/operations/LOCAL_CONTROL_PLANE_QUICKSTART.md"
reset_script="$root_dir/tools/workspace_reset.sh"

doc_test_init

doc_assert_exists "$compose_file" "Local compose file exists"
doc_assert_contains "$compose_file" "services:" "Compose defines services root"
doc_assert_contains "$compose_file" "  db:" "Compose includes db service"
doc_assert_contains "$compose_file" "  api:" "Compose includes api service"
doc_assert_contains "$compose_file" "  ui:" "Compose includes ui service"
doc_assert_contains "$compose_file" 'healthcheck:' "Compose includes health checks"
doc_assert_contains "$compose_file" 'AZURE_OPENAI_ENDPOINT' "Compose supports Azure OpenAI endpoint configuration"
doc_assert_contains "$compose_file" 'AZURE_OPENAI_API_KEY' "Compose supports Azure OpenAI key configuration"
doc_assert_contains "$compose_file" 'OPENAI_API_KEY' "Compose supports OpenAI key configuration"
doc_assert_contains "$compose_file" "volumes:" "Compose defines volumes"
doc_assert_contains "$compose_file" 'athenawork-db-data' "Compose defines persistent db volume"

doc_assert_exists "$quickstart" "Local control plane quickstart exists"
doc_assert_contains "$quickstart" 'docker compose -f docker-compose.local.yml up --build -d' "Quickstart documents startup command"
doc_assert_contains "$quickstart" 'docker compose -f docker-compose.local.yml down' "Quickstart documents teardown command"
doc_assert_contains "$quickstart" '/api/v1/model/respond' "Quickstart documents model response endpoint check"
doc_assert_contains "$quickstart" './tools/workspace_reset.sh' "Quickstart documents reset command"

doc_assert_exists "$reset_script" "Workspace reset script exists"
doc_assert_contains "$reset_script" 'docker compose -f docker-compose.local.yml down -v --remove-orphans' "Reset script removes compose volumes"

doc_test_finish

#!/usr/bin/env bash
set -euo pipefail

script_dir="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
root_dir="$(git -C "$script_dir" rev-parse --show-toplevel 2>/dev/null || (cd "$script_dir/.." && pwd))"
source "$root_dir/tools/lib/doc_test_harness.sh"

gate_script="$root_dir/products/athena-work/tools/check_security_nonce_gate.sh"
nonce_script="$root_dir/products/athena-work/tools/generate_security_nonce.sh"
contracts_file="$root_dir/products/athena-work/operating-system/contracts/SECURITY_BACKING_CONTRACT_PATHS.txt"
policy_file="$root_dir/products/athena-work/operating-system/contracts/SECURITY_CHANGE_CONTROL_POLICY_V1.md"
pipeline_file="$root_dir/azure-pipelines.yml"

doc_test_init

doc_assert_exists "$gate_script" "Security nonce gate script exists"
doc_assert_exists "$nonce_script" "Security nonce generator script exists"
doc_assert_exists "$contracts_file" "Security backing contracts path list exists"
doc_assert_exists "$policy_file" "Security change control policy exists"

doc_assert_contains "$contracts_file" "AGENTS.md" "AGENTS.md is protected by security gate"
doc_assert_contains "$gate_script" "ATHENA_SECURITY_CHANGE_NONCE" "Gate requires nonce input"
doc_assert_contains "$gate_script" "ATHENA_SECURITY_CHANGE_OTP_PRIMARY" "Gate requires primary OTP input"
doc_assert_contains "$gate_script" "ATHENA_SECURITY_CHANGE_OTP_SECURITY" "Gate requires security OTP input"
doc_assert_contains "$gate_script" "ATHENA_SECURITY_GATE_WINDOW_MINUTES" "Gate supports tunable nonce window"
doc_assert_contains "$gate_script" "LAUNCH AUTHORIZED" "Gate enforces launch authorization phrase"

doc_assert_contains "$nonce_script" "security_nonce=" "Nonce generator emits nonce output"
doc_assert_contains "$nonce_script" "ATHENA_SECURITY_CHANGE_NONCE" "Nonce generator emits export guidance"
doc_assert_contains "$policy_file" "Dual confirmation phrase is required" "Policy defines dual confirmation requirement"
doc_assert_contains "$policy_file" "CI must run" "Policy defines CI enforcement"
doc_assert_contains "$pipeline_file" "check_security_nonce_gate.sh" "Azure pipeline executes security nonce gate"

doc_test_finish

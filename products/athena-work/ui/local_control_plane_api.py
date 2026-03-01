from http.server import BaseHTTPRequestHandler, HTTPServer
import json
from datetime import datetime, timezone
from pathlib import Path
import os
from urllib.parse import urlparse, parse_qs, quote
import base64
import subprocess
import tempfile
import urllib.request
import urllib.error

DEFAULT_STATE_FILE = Path("/workspace/products/athena-work/operating-system/state/backend_read_model_v1.json")
RUNTIME_STATE_FILE = Path("/workspace/products/athena-work/operating-system/state/runtime/backend_read_model_v1.local.json")
MIGRATION_LOG_FILE = Path("/workspace/products/athena-work/operating-system/state/runtime/migration_events_v1.local.json")
ALLOWED_THEMES = {"default", "graphite", "ember", "high-contrast"}
CORE_DOCS_INDEX = [
    {
        "id": "humans",
        "title": "HUMANS Guide",
        "summary": "Operator workflow and stage transition conventions.",
        "path": "products/athena-work/HUMANS.md",
    },
    {
        "id": "development-cycle",
        "title": "Development Cycle",
        "summary": "Lifecycle and queue movement rules for delivery.",
        "path": "products/athena-work/DEVELOPMENT_CYCLE.md",
    },
    {
        "id": "stage-exit-gates",
        "title": "Stage Exit Gates",
        "summary": "Criteria for planning, PM, engineering, and QA exits.",
        "path": "knowledge-base/process/STAGE_EXIT_GATES.md",
    },
    {
        "id": "local-control-quickstart",
        "title": "Local Control Plane Quickstart",
        "summary": "Runbook for starting and validating local workspace services.",
        "path": "knowledge-base/operations/LOCAL_CONTROL_PLANE_QUICKSTART.md",
    },
]
INGEST_TOOL = "/workspace/products/athena-work/tools/ingest_artifact_bundle.sh"
LAUNCH_PACKAGE_TOOL = "/workspace/products/athena-work/tools/generate_launch_authorization_package.sh"
LAUNCH_VALIDATE_TOOL = "/workspace/products/athena-work/tools/validate_launch_authorization_package.sh"
LAUNCH_PACKAGE_DIR = Path("/workspace/products/athena-work/operating-system/handoff")


def now_utc():
    return datetime.now(timezone.utc).strftime("%Y-%m-%dT%H:%M:%SZ")


def load_state():
    runtime_path = Path(os.getenv("ATHENA_RUNTIME_STATE_FILE", str(RUNTIME_STATE_FILE)))
    if runtime_path.exists():
        return json.loads(runtime_path.read_text(encoding="utf-8"))
    if DEFAULT_STATE_FILE.exists():
        return json.loads(DEFAULT_STATE_FILE.read_text(encoding="utf-8"))
    return {}


def workspace_roots():
    roots = [Path("/workspace"), Path(__file__).resolve().parents[3]]
    seen = set()
    ordered = []
    for root in roots:
        key = str(root)
        if key in seen:
            continue
        seen.add(key)
        ordered.append(root)
    return ordered


def resolve_workspace_path(path):
    relative = str(path or "").lstrip("/")
    for root in workspace_roots():
        candidate = root / relative
        if candidate.exists():
            return candidate
    return workspace_roots()[0] / relative


def _normalize_card_ids(values):
    normalized = []
    seen = set()
    for raw in values or []:
        value = str(raw or "").strip()
        if not value or value in seen:
            continue
        seen.add(value)
        normalized.append(value)
    return normalized


def _normalize_timeline(values):
    events = []
    for item in values or []:
        if not isinstance(item, dict):
            continue
        events.append(
            {
                "event_type": str(item.get("event_type") or "").strip() or "unknown",
                "result": str(item.get("result") or "").strip() or "unknown",
                "label": str(item.get("label") or "").strip() or "event",
                "story_id": str(item.get("story_id") or "").strip() or "unknown-story",
                "cycle_id": str(item.get("cycle_id") or "").strip() or "unknown-cycle",
                "correlation_id": str(item.get("correlation_id") or "").strip(),
                "timestamp": str(item.get("timestamp") or "").strip(),
            }
        )
    return events


def build_workspace_entity_model(state):
    engineering = state.get("engineering", {}) if isinstance(state.get("engineering"), dict) else {}
    architecture = state.get("architecture", {}) if isinstance(state.get("architecture"), dict) else {}

    eng_active = _normalize_card_ids(engineering.get("active"))
    eng_qa = _normalize_card_ids(engineering.get("qa"))
    arch_active = _normalize_card_ids(architecture.get("active"))
    timeline = _normalize_timeline(state.get("timeline"))
    blockers = _normalize_card_ids(state.get("blockers") or ["none"])
    next_story = str(state.get("next_story") or "none").strip() or "none"

    cards = []
    for card_id in eng_active:
        cards.append({"id": card_id, "lane": "engineering.active", "kind": "story"})
    for card_id in eng_qa:
        cards.append({"id": card_id, "lane": "engineering.qa", "kind": "story"})
    for card_id in arch_active:
        cards.append({"id": card_id, "lane": "architecture.active", "kind": "architecture_story"})

    cycles = sorted({event.get("cycle_id", "unknown-cycle") for event in timeline})
    docs = [{"id": item["id"], "path": item["path"]} for item in CORE_DOCS_INDEX]

    card_ids = {item["id"] for item in cards}
    relationship_errors = []

    eng_active_set = set(eng_active)
    eng_qa_set = set(eng_qa)
    arch_active_set = set(arch_active)
    overlap_eng = sorted(eng_active_set.intersection(eng_qa_set))
    overlap_arch = sorted((eng_active_set.union(eng_qa_set)).intersection(arch_active_set))

    if overlap_eng:
        relationship_errors.append(
            f"card_in_multiple_engineering_lanes:{','.join(overlap_eng)}"
        )
    if overlap_arch:
        relationship_errors.append(
            f"card_crosses_engineering_architecture_lanes:{','.join(overlap_arch)}"
        )

    if next_story != "none" and next_story not in card_ids:
        relationship_errors.append(f"next_story_not_in_active_or_qa:{next_story}")

    for event in timeline:
        if not event["correlation_id"]:
            relationship_errors.append(
                f"timeline_missing_correlation_id:{event['label']}:{event['story_id']}"
            )
        if not event["timestamp"]:
            relationship_errors.append(
                f"timeline_missing_timestamp:{event['label']}:{event['story_id']}"
            )

    entity_model = {
        "schema_version": "1",
        "state_version": state.get("state_version", "unknown"),
        "entities": {
            "cards": cards,
            "cycles": cycles,
            "blockers": blockers,
            "docs": docs,
        },
        "relationships": {
            "engineering_active_cards": eng_active,
            "engineering_qa_cards": eng_qa,
            "architecture_active_cards": arch_active,
            "timeline_events": timeline,
            "next_story": next_story,
        },
        "relationship_errors": relationship_errors,
    }
    return entity_model


def _resolve_story_path(story_id):
    candidates = [
        Path(f"products/athena-work/delivery-backlog/engineering/active/{story_id}.md"),
        Path(f"products/athena-work/delivery-backlog/engineering/qa/{story_id}.md"),
        Path(f"products/athena-work/delivery-backlog/engineering/done/{story_id}.md"),
    ]
    for candidate in candidates:
        resolved = resolve_workspace_path(candidate)
        if resolved.exists():
            return str(candidate)
    return ""


def build_docs_index(state=None):
    if state is None:
        state = load_state()
    model = build_workspace_entity_model(state)

    docs = list(CORE_DOCS_INDEX)
    docs.append(
        {
            "id": "engineering-active-queue",
            "title": "Engineering Active Queue",
            "summary": "Canonical ranked execution queue for engineering stage launch.",
            "path": "products/athena-work/delivery-backlog/engineering/active/README.md",
        }
    )

    story_ids = []
    next_story = model["relationships"].get("next_story", "none")
    if next_story and next_story != "none":
        story_ids.append(next_story)
    story_ids.extend(model["relationships"].get("engineering_active_cards", []))
    story_ids.extend(model["relationships"].get("engineering_qa_cards", []))

    seen_story = set()
    for story_id in story_ids:
        if story_id in seen_story:
            continue
        seen_story.add(story_id)
        path = _resolve_story_path(story_id)
        if not path:
            continue
        docs.append(
            {
                "id": f"story-{story_id.lower()}",
                "title": f"Story Link: {story_id}",
                "summary": "Auto-linked from canonical workspace state.",
                "path": path,
            }
        )

    seen_paths = set()
    ordered = []
    for item in docs:
        path = item.get("path", "")
        if not path or path in seen_paths:
            continue
        seen_paths.add(path)
        ordered.append(item)
    return {"generated_at": now_utc(), "docs": ordered}


def build_board():
    state = load_state()
    entity_model = build_workspace_entity_model(state)
    relationship_errors = entity_model.get("relationship_errors", [])
    mode = os.getenv("ATHENA_UI_ACCESSIBILITY_MODE", "normal").strip().lower()
    if mode not in {"normal", "low-vision"}:
        mode = "normal"
    theme = os.getenv("ATHENA_UI_THEME", "default").strip().lower() or "default"
    if theme not in ALLOWED_THEMES:
        theme = "default"
    return {
        "generated_at": now_utc(),
        "schema_version": state.get("schema_version", "1"),
        "state_version": state.get("state_version", "unknown"),
        "current_stage": state.get("current_stage", "engineering"),
        "next_story": state.get("next_story", "none"),
        "blockers": state.get("blockers", ["none"]),
        "required_confirmation": state.get(
            "required_confirmation",
            {
                "required": True,
                "status": "unconfirmed",
                "reason": "Direction-changing actions require valid human confirmation.",
            },
        ),
        "engineering": {
            "active": entity_model["relationships"]["engineering_active_cards"],
            "qa": entity_model["relationships"]["engineering_qa_cards"],
            "done_count": (state.get("engineering", {}) or {}).get("done_count", 0),
        },
        "architecture": {
            "active": entity_model["relationships"]["architecture_active_cards"],
            "qa_count": (state.get("architecture", {}) or {}).get("qa_count", 0),
        },
        "drift_alerts": (state.get("drift_alerts", ["none"]) or []) + relationship_errors,
        "research_comm_exception": state.get(
            "research_comm_exception",
            {
                "active": False,
                "note": "Undocumented agent communication is blocked outside research mode.",
            },
        ),
        "entity_model": {
            "schema_version": entity_model["schema_version"],
            "state_version": entity_model["state_version"],
            "relationship_error_count": len(relationship_errors),
        },
        "ui_defaults": {
            "accessibility_mode": mode,
            "theme": theme,
            "available_themes": sorted(ALLOWED_THEMES),
        },
    }


def build_timeline():
    state = load_state()
    entity_model = build_workspace_entity_model(state)
    return {
        "generated_at": now_utc(),
        "events": entity_model["relationships"]["timeline_events"],
        "entity_model": {
            "schema_version": entity_model["schema_version"],
            "state_version": entity_model["state_version"],
            "relationship_error_count": len(entity_model.get("relationship_errors", [])),
        },
    }


def build_doc_view(doc_id, docs):
    doc = next((item for item in docs if item["id"] == doc_id), None)
    if not doc:
        return None
    doc_path = resolve_workspace_path(doc["path"])
    if not doc_path.exists():
        return {
            "id": doc["id"],
            "title": doc["title"],
            "path": doc["path"],
            "content": "Document not found in local workspace.",
        }
    content = doc_path.read_text(encoding="utf-8")
    return {
        "id": doc["id"],
        "title": doc["title"],
        "path": doc["path"],
        "content": content[:12000],
    }


def build_latest_launch_package():
    files = sorted(LAUNCH_PACKAGE_DIR.glob("LAUNCH_AUTHZ_*.json"))
    if not files:
        return {"code": "NO_PACKAGE", "reason": "no launch package generated yet"}
    latest = files[-1]
    payload = json.loads(latest.read_text(encoding="utf-8"))
    payload["path"] = str(latest).replace("/workspace/", "")
    payload["code"] = "OK"
    return payload


def migration_log_path():
    return Path(os.getenv("ATHENA_MIGRATION_LOG_FILE", str(MIGRATION_LOG_FILE)))


def load_migration_events():
    path = migration_log_path()
    if not path.exists():
        return []
    try:
        payload = json.loads(path.read_text(encoding="utf-8"))
    except Exception:  # noqa: BLE001
        return []
    if not isinstance(payload, dict):
        return []
    events = payload.get("events")
    if not isinstance(events, list):
        return []
    normalized = []
    for item in events:
        if not isinstance(item, dict):
            continue
        normalized.append(
            {
                "ingested_at": str(item.get("ingested_at") or "").strip(),
                "filename": str(item.get("filename") or "").strip(),
                "code": str(item.get("code") or "").strip(),
                "reason": str(item.get("reason") or "").strip(),
                "accepted_objects": int(item.get("accepted_objects") or 0),
                "rejected_objects": int(item.get("rejected_objects") or 0),
                "migrated_objects": int(item.get("migrated_objects") or 0),
                "source_object": str(item.get("source_object") or "").strip(),
                "output_state_file": str(item.get("output_state_file") or "").strip(),
                "migration_path": str(item.get("migration_path") or "").strip(),
            }
        )
    return normalized


def append_migration_event(event):
    path = migration_log_path()
    events = load_migration_events()
    events.insert(0, event)
    trimmed = events[:50]
    payload = {
        "schema_version": "1",
        "generated_at": now_utc(),
        "read_only": True,
        "events": trimmed,
    }
    path.parent.mkdir(parents=True, exist_ok=True)
    path.write_text(json.dumps(payload, indent=2) + "\n", encoding="utf-8")


def build_migration_view():
    events = load_migration_events()
    return {
        "code": "OK",
        "generated_at": now_utc(),
        "read_only": True,
        "event_count": len(events),
        "events": events,
    }


def resolve_model_provider():
    azure_endpoint = (os.getenv("AZURE_OPENAI_ENDPOINT") or "").strip().rstrip("/")
    azure_deployment = (os.getenv("AZURE_OPENAI_DEPLOYMENT_NAME") or "").strip()
    azure_api_version = (os.getenv("AZURE_OPENAI_API_VERSION") or "").strip() or "2024-10-21"
    azure_api_key = (os.getenv("AZURE_OPENAI_API_KEY") or "").strip()

    if azure_endpoint and azure_deployment and azure_api_key:
        return {
            "provider": "azure_openai",
            "endpoint": azure_endpoint,
            "deployment": azure_deployment,
            "api_version": azure_api_version,
            "api_key": azure_api_key,
        }

    openai_api_key = (os.getenv("OPENAI_API_KEY") or "").strip()
    openai_model = (os.getenv("OPENAI_MODEL") or "").strip() or "gpt-4o-mini"
    if openai_api_key:
        return {
            "provider": "openai",
            "model": openai_model,
            "api_key": openai_api_key,
        }
    return None


def model_chat_completion(prompt):
    provider = resolve_model_provider()
    if not provider:
        raise RuntimeError(
            "model provider not configured; set AZURE_OPENAI_ENDPOINT/AZURE_OPENAI_API_KEY/"
            "AZURE_OPENAI_DEPLOYMENT_NAME or OPENAI_API_KEY"
        )

    timeout_s = 20
    try:
        timeout_s = int(os.getenv("ATHENA_MODEL_TIMEOUT_S", "20"))
    except Exception:  # noqa: BLE001
        timeout_s = 20
    timeout_s = max(5, min(timeout_s, 60))

    payload = {
        "messages": [{"role": "user", "content": prompt}],
        "temperature": 0.2,
    }
    headers = {"content-type": "application/json"}
    descriptor = ""
    if provider["provider"] == "azure_openai":
        url = (
            f"{provider['endpoint']}/openai/deployments/{quote(provider['deployment'])}"
            f"/chat/completions?api-version={provider['api_version']}"
        )
        headers["api-key"] = provider["api_key"]
        descriptor = provider["deployment"]
    else:
        url = "https://api.openai.com/v1/chat/completions"
        payload["model"] = provider["model"]
        headers["authorization"] = f"Bearer {provider['api_key']}"
        descriptor = provider["model"]

    req = urllib.request.Request(
        url,
        data=json.dumps(payload).encode("utf-8"),
        headers=headers,
        method="POST",
    )

    try:
        with urllib.request.urlopen(req, timeout=timeout_s) as resp:
            body = resp.read().decode("utf-8")
    except urllib.error.HTTPError as exc:
        detail = exc.read().decode("utf-8", errors="replace")
        raise RuntimeError(f"model request failed ({exc.code}): {detail[:600]}") from exc
    except Exception as exc:  # noqa: BLE001
        raise RuntimeError(f"model request failed: {exc}") from exc

    parsed = json.loads(body or "{}")
    choices = parsed.get("choices") if isinstance(parsed, dict) else None
    if not choices or not isinstance(choices, list):
        raise RuntimeError("model response did not include choices")
    message = choices[0].get("message", {}) if isinstance(choices[0], dict) else {}
    content = message.get("content") if isinstance(message, dict) else None
    text = str(content or "").strip()
    if not text:
        raise RuntimeError("model response was empty")

    return {
        "provider": provider["provider"],
        "model_or_deployment": descriptor,
        "response_text": text,
    }


class Handler(BaseHTTPRequestHandler):
    def write_cors_headers(self):
        self.send_header("access-control-allow-origin", "*")
        self.send_header("access-control-allow-methods", "GET, POST, OPTIONS")
        self.send_header("access-control-allow-headers", "content-type")

    def send_json(self, payload, code=200):
        body = json.dumps(payload).encode("utf-8")
        self.send_response(code)
        self.send_header("content-type", "application/json")
        self.write_cors_headers()
        self.send_header("content-length", str(len(body)))
        self.end_headers()
        self.wfile.write(body)

    def do_OPTIONS(self):
        self.send_response(204)
        self.write_cors_headers()
        self.end_headers()

    def do_POST(self):
        parsed = urlparse(self.path)
        path = parsed.path

        if path == "/api/v1/planning/export":
            self.handle_planning_export()
            return
        if path == "/api/v1/artifacts/ingest":
            self.handle_artifact_ingest()
            return
        if path == "/api/v1/launch/package":
            self.handle_launch_package()
            return
        if path == "/api/v1/launch/validate":
            self.handle_launch_validate()
            return
        if path == "/api/v1/model/respond":
            self.handle_model_respond()
            return
        self.send_json({"code": "ERR_ROUTE_NOT_FOUND", "reason": "unsupported path"}, 404)

    def handle_planning_export(self):
        if self.path != "/api/v1/planning/export":
            self.send_json({"code": "ERR_ROUTE_NOT_FOUND", "reason": "unsupported path"}, 404)
            return

        content_length = int(self.headers.get("content-length", "0"))
        payload = json.loads(self.rfile.read(content_length) or "{}")
        direction = (payload.get("direction") or "").strip()
        constraints = (payload.get("constraints") or "").strip()
        risks = (payload.get("risks") or "").strip()
        next_stage = (payload.get("next_stage") or "").strip()
        confirmed_by = (payload.get("confirmed_by") or "").strip()

        if not direction or not constraints or not risks or not next_stage:
            self.send_json(
                {
                    "code": "ERR_API_INPUT_INVALID",
                    "reason": "direction, constraints, risks, and next_stage are required",
                },
                400,
            )
            return

        ts = datetime.now(timezone.utc).strftime("%Y%m%dT%H%M%SZ")
        out_dir = Path("/workspace/product-research/planning/sessions")
        out_dir.mkdir(parents=True, exist_ok=True)
        out_path = out_dir / f"PLAN-WORKBENCH-{ts}.md"
        out_path.write_text(
            "\n".join(
                [
                    "# Planning Workbench Export",
                    "",
                    "## Metadata",
                    f"- `created_at_utc`: {now_utc()}",
                    f"- `source`: workspace-ui-planning-workbench-v1",
                    "",
                    "## Direction",
                    direction,
                    "",
                    "## Constraints",
                    constraints,
                    "",
                    "## Risks",
                    risks,
                    "",
                    "## Next Stage",
                    next_stage,
                    "",
                    "## Confirmation",
                    f"- `confirmed_by`: {confirmed_by or 'unconfirmed'}",
                ]
            )
            + "\n",
            encoding="utf-8",
        )
        self.send_json(
            {
                "code": "OK",
                "reason": "planning artifact exported",
                "path": str(out_path).replace("/workspace/", ""),
            }
        )

    def handle_artifact_ingest(self):
        content_length = int(self.headers.get("content-length", "0"))
        payload = json.loads(self.rfile.read(content_length) or "{}")
        filename = (payload.get("filename") or "").strip()
        bundle_b64 = payload.get("bundle_base64") or ""

        if not filename.lower().endswith(".zip"):
            self.send_json(
                {"code": "ERR_API_INPUT_INVALID", "reason": "filename must end with .zip"},
                400,
            )
            return
        if not bundle_b64:
            self.send_json(
                {"code": "ERR_API_INPUT_INVALID", "reason": "bundle_base64 is required"},
                400,
            )
            return

        try:
            bundle_bytes = base64.b64decode(bundle_b64, validate=True)
        except Exception:  # noqa: BLE001
            self.send_json(
                {"code": "ERR_API_INPUT_INVALID", "reason": "bundle_base64 is not valid base64"},
                400,
            )
            return

        if len(bundle_bytes) > 10 * 1024 * 1024:
            self.send_json(
                {"code": "ERR_API_INPUT_INVALID", "reason": "bundle exceeds 10MB limit"},
                400,
            )
            return

        with tempfile.TemporaryDirectory(prefix="athena-ingest-") as tmp_dir:
            zip_path = Path(tmp_dir) / filename
            zip_path.write_bytes(bundle_bytes)
            proc = subprocess.run(
                [INGEST_TOOL, "--zip", str(zip_path), "--root", "/workspace"],
                capture_output=True,
                text=True,
                check=False,
            )

        body = (proc.stdout or "").strip()
        if not body:
            body = json.dumps(
                {
                    "code": "ERR_ARTIFACT_INGEST_FAILED",
                    "reason": (proc.stderr or "unknown ingest failure").strip(),
                }
            )
        try:
            parsed = json.loads(body)
        except Exception:  # noqa: BLE001
            parsed = {
                "code": "ERR_ARTIFACT_INGEST_FAILED",
                "reason": body[:500],
            }

        migration_event = {
            "ingested_at": now_utc(),
            "filename": filename,
            "code": str(parsed.get("code") or "ERR_ARTIFACT_INGEST_FAILED"),
            "reason": str(parsed.get("reason") or ""),
            "accepted_objects": int(parsed.get("accepted_objects") or 0),
            "rejected_objects": int(parsed.get("rejected_objects") or 0),
            "migrated_objects": int(parsed.get("migrated_objects") or 0),
            "source_object": str(parsed.get("source_object") or ""),
            "output_state_file": str(parsed.get("output_state_file") or ""),
            "migration_path": "v1->v2" if int(parsed.get("migrated_objects") or 0) > 0 else "v2->v2",
        }
        append_migration_event(migration_event)

        status = 200 if proc.returncode == 0 else 400
        self.send_json(parsed, status)

    def handle_launch_package(self):
        proc = subprocess.run(
            [LAUNCH_PACKAGE_TOOL, "--root", "/workspace"],
            capture_output=True,
            text=True,
            check=False,
        )
        body = (proc.stdout or "").strip()
        if not body:
            body = json.dumps(
                {
                    "code": "ERR_LAUNCH_PACKAGE_FAILED",
                    "reason": (proc.stderr or "unknown launch package failure").strip(),
                }
            )
        try:
            parsed = json.loads(body)
        except Exception:  # noqa: BLE001
            parsed = {"code": "ERR_LAUNCH_PACKAGE_FAILED", "reason": body[:500]}
        status = 200 if proc.returncode == 0 else 400
        self.send_json(parsed, status)

    def handle_launch_validate(self):
        content_length = int(self.headers.get("content-length", "0"))
        payload = json.loads(self.rfile.read(content_length) or "{}")
        package_path = (payload.get("package_path") or "").strip()
        if not package_path:
            self.send_json(
                {"code": "ERR_API_INPUT_INVALID", "reason": "package_path is required"},
                400,
            )
            return
        full_path = package_path
        if not full_path.startswith("/workspace/"):
            full_path = f"/workspace/{package_path.lstrip('/')}"
        proc = subprocess.run(
            [LAUNCH_VALIDATE_TOOL, "--package", full_path, "--root", "/workspace"],
            capture_output=True,
            text=True,
            check=False,
        )
        if proc.returncode == 0:
            self.send_json(
                {
                    "code": "OK",
                    "reason": "launch package valid",
                    "output": (proc.stdout or "").strip(),
                }
            )
            return
        self.send_json(
            {
                "code": "ERR_LAUNCH_PACKAGE_INVALID",
                "reason": (proc.stdout or proc.stderr or "validation failed").strip(),
            },
            400,
        )

    def handle_model_respond(self):
        content_length = int(self.headers.get("content-length", "0"))
        try:
            payload = json.loads(self.rfile.read(content_length) or "{}")
        except Exception:  # noqa: BLE001
            self.send_json(
                {"code": "ERR_API_INPUT_INVALID", "reason": "body must be valid JSON"},
                400,
            )
            return

        prompt = (payload.get("prompt") or "").strip()
        if not prompt:
            self.send_json(
                {"code": "ERR_API_INPUT_INVALID", "reason": "prompt is required"},
                400,
            )
            return
        if len(prompt) > 8000:
            self.send_json(
                {"code": "ERR_API_INPUT_INVALID", "reason": "prompt exceeds 8000 characters"},
                400,
            )
            return

        try:
            result = model_chat_completion(prompt)
        except Exception as exc:  # noqa: BLE001
            self.send_json(
                {"code": "ERR_MODEL_PROVIDER", "reason": str(exc)},
                400,
            )
            return

        self.send_json(
            {
                "code": "OK",
                "generated_at": now_utc(),
                "provider": result["provider"],
                "model_or_deployment": result["model_or_deployment"],
                "response_text": result["response_text"],
            }
        )

    def do_GET(self):
        parsed = urlparse(self.path)
        path = parsed.path
        query = parse_qs(parsed.query)

        if path == "/health":
            self.send_json({"status": "ok", "service": "api"})
            return
        if path == "/api/v1/read-model/board":
            self.send_json(build_board())
            return
        if path == "/api/v1/read-model/timeline":
            self.send_json(build_timeline())
            return
        if path == "/api/v1/docs/index":
            self.send_json(build_docs_index())
            return
        if path == "/api/v1/docs/view":
            doc_id = (query.get("id", [""])[0] or "").strip()
            docs = build_docs_index().get("docs", [])
            payload = build_doc_view(doc_id, docs)
            if payload is None:
                self.send_json({"code": "ERR_DOC_NOT_FOUND", "reason": "unknown doc id"}, 404)
                return
            self.send_json(payload)
            return
        if path == "/api/v1/artifacts/migrations":
            self.send_json(build_migration_view())
            return
        if path == "/api/v1/launch/latest":
            payload = build_latest_launch_package()
            status = 200 if payload.get("code") == "OK" else 404
            self.send_json(payload, status)
            return
        self.send_json(
            {
                "service": "api",
                "message": "AthenaWork local control-plane bootstrap",
            }
        )

    def log_message(self, fmt, *args):
        return


if __name__ == "__main__":
    HTTPServer(("0.0.0.0", 8787), Handler).serve_forever()

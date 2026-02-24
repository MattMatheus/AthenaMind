# Local Experiment Infra (Podman)

Use this stack to expand retrieval experiments locally at zero cloud cost.

## Services
- Qdrant: vector search experiments (`http://localhost:6333`)
- Neo4j: memory graph experiments (`bolt://localhost:7687`)
- Postgres+pgvector: hybrid SQL/vector experiments (`localhost:5432`)

## Start
```bash
podman compose -f docker-compose.experiment.yml up -d
```

## Stop
```bash
podman compose -f docker-compose.experiment.yml down
```

## Reset Volumes
```bash
podman compose -f docker-compose.experiment.yml down -v
```

## Notes
- This stack is optional and does not change current AthenaMind runtime behavior.
- Current retrieval remains SQLite + local embedding endpoint unless explicitly changed in code.
- For Neo4j dev login, default is `neo4j/devpassword`.

## Retrieval Backend Env
When using `--retrieval-backend qdrant`, configure:
- `ATHENA_QDRANT_URL` (default `http://localhost:6333`)
- `ATHENA_QDRANT_COLLECTION` (default `athena_memories`)
- `ATHENA_QDRANT_API_KEY` (optional)

When using `--retrieval-backend neo4j`, configure:
- `ATHENA_NEO4J_HTTP_URL` (default `http://localhost:7474`)
- `ATHENA_NEO4J_USER` (default `neo4j`)
- `ATHENA_NEO4J_PASSWORD` (default `devpassword`)
- `ATHENA_NEO4J_DATABASE` (default `neo4j`)

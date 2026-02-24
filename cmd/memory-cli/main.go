package main

import (
	"context"
	"errors"
	"fmt"
	"os"

	"athenamind/internal/telemetry"
)

func main() {
	shutdown, initErr := telemetry.InitOTel(context.Background())
	if initErr != nil {
		fmt.Fprintf(os.Stderr, "warning: otel initialization failed: %v\n", initErr)
		shutdown = func(context.Context) error { return nil }
	}
	runShutdown := func() {
		if err := shutdown(context.Background()); err != nil {
			fmt.Fprintf(os.Stderr, "warning: otel shutdown failed: %v\n", err)
		}
	}
	defer runShutdown()

	if len(os.Args) < 2 {
		exitErr(errors.New("usage: memory-cli <write|retrieve|evaluate|bootstrap|verify|episode|reembed-changed|sync-qdrant|telemetry> [flags]"))
	}

	var err error
	switch os.Args[1] {
	case "write":
		err = runWrite(os.Args[2:])
	case "retrieve":
		err = runRetrieve(os.Args[2:])
	case "snapshot":
		err = runSnapshot(os.Args[2:])
	case "serve-read-gateway":
		err = runServeReadGateway(os.Args[2:])
	case "api-retrieve":
		err = runAPIRetrieve(os.Args[2:])
	case "evaluate":
		err = runEvaluate(os.Args[2:])
	case "bootstrap":
		err = runBootstrap(os.Args[2:])
	case "reindex-all":
		err = runReindexAll(os.Args[2:])
	case "crawl":
		err = runCrawl(os.Args[2:])
	case "reembed-changed":
		err = runReembedChanged(os.Args[2:])
	case "sync-qdrant":
		err = runSyncQdrant(os.Args[2:])
	case "verify":
		err = runVerify(os.Args[2:])
	case "episode":
		err = runEpisode(os.Args[2:])
	case "telemetry":
		err = runTelemetry(os.Args[2:])
	default:
		err = fmt.Errorf("unknown command: %s", os.Args[1])
	}

	if err != nil {
		runShutdown()
		exitErr(err)
	}
}

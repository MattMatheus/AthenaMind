package main

import (
	"errors"
	"fmt"
	"os"
)

func main() {
	if len(os.Args) < 2 {
		exitErr(errors.New("usage: memory-cli <write|retrieve|evaluate> [flags]"))
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
	default:
		err = fmt.Errorf("unknown command: %s", os.Args[1])
	}

	if err != nil {
		exitErr(err)
	}
}

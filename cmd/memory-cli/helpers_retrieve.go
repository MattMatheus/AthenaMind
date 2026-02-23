package main

import (
	"net/http"

	"athenamind/internal/gateway"
	"athenamind/internal/retrieval"
)

func retrieve(root, query, domain string) (retrieveResult, error) {
	return retrieval.Retrieve(root, query, domain)
}

func readGatewayHandler(root string) http.Handler {
	return gateway.ReadGatewayHandler(root)
}

func apiRetrieveWithFallback(root, gatewayURL string, req apiRetrieveRequest, traceID string, client *http.Client) (apiRetrieveResponse, error) {
	return gateway.APIRetrieveWithFallback(root, gatewayURL, req, traceID, client)
}

func gatewayRetrieve(gatewayURL string, req apiRetrieveRequest, client *http.Client) (apiRetrieveResponse, error) {
	return gateway.GatewayRetrieve(gatewayURL, req, client)
}

func isSemanticConfident(top, second float64) bool {
	return retrieval.IsSemanticConfident(top, second)
}

func indexEntryEmbedding(root, entryID, endpoint, sessionID string) (string, error) {
	return retrieval.IndexEntryEmbedding(root, entryID, endpoint, sessionID)
}

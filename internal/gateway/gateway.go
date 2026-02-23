package gateway

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"athenamind/internal/retrieval"
	"athenamind/internal/types"
)

func ReadGatewayHandler(root string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}

		var req types.APIRetrieveRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "ERR_API_INPUT_INVALID: invalid JSON payload", http.StatusBadRequest)
			return
		}
		if strings.TrimSpace(req.Query) == "" || strings.TrimSpace(req.SessionID) == "" {
			http.Error(w, "ERR_API_INPUT_INVALID: query and session_id are required", http.StatusBadRequest)
			return
		}

		result, err := retrieval.Retrieve(root, req.Query, req.Domain)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		traceID := fmt.Sprintf("%s:%d", strings.TrimSpace(req.SessionID), time.Now().UTC().UnixNano())
		resp := types.APIRetrieveResponse{
			SelectedID:    result.SelectedID,
			SelectionMode: result.SelectionMode,
			SourcePath:    result.SourcePath,
			Confidence:    result.Confidence,
			Reason:        result.Reason,
			TraceID:       traceID,
		}
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(resp)
	})
}

func APIRetrieveWithFallback(root, gatewayURL string, req types.APIRetrieveRequest, traceID string, client *http.Client) (types.APIRetrieveResponse, error) {
	if strings.TrimSpace(gatewayURL) != "" {
		resp, err := GatewayRetrieve(gatewayURL, req, client)
		if err == nil {
			local, localErr := retrieval.Retrieve(root, req.Query, req.Domain)
			if localErr != nil {
				return types.APIRetrieveResponse{}, localErr
			}
			if resp.SelectedID != local.SelectedID || resp.SelectionMode != local.SelectionMode || resp.SourcePath != local.SourcePath {
				return types.APIRetrieveResponse{}, errors.New("ERR_API_CLI_PARITY_MISMATCH: gateway response diverged from CLI retrieve contract")
			}
			if strings.TrimSpace(resp.TraceID) == "" {
				resp.TraceID = traceID
			}
			resp.GatewayEndpoint = gatewayURL
			return resp, nil
		}
		local, localErr := retrieval.Retrieve(root, req.Query, req.Domain)
		if localErr != nil {
			return types.APIRetrieveResponse{}, localErr
		}
		return types.APIRetrieveResponse{
			SelectedID:      local.SelectedID,
			SelectionMode:   local.SelectionMode,
			SourcePath:      local.SourcePath,
			Confidence:      local.Confidence,
			Reason:          local.Reason,
			TraceID:         traceID,
			FallbackUsed:    true,
			FallbackCode:    "ERR_API_WRAPPER_UNAVAILABLE",
			FallbackReason:  err.Error(),
			GatewayEndpoint: gatewayURL,
		}, nil
	}

	local, err := retrieval.Retrieve(root, req.Query, req.Domain)
	if err != nil {
		return types.APIRetrieveResponse{}, err
	}
	return types.APIRetrieveResponse{
		SelectedID:    local.SelectedID,
		SelectionMode: local.SelectionMode,
		SourcePath:    local.SourcePath,
		Confidence:    local.Confidence,
		Reason:        local.Reason,
		TraceID:       traceID,
	}, nil
}

func GatewayRetrieve(gatewayURL string, req types.APIRetrieveRequest, client *http.Client) (types.APIRetrieveResponse, error) {
	body, err := json.Marshal(req)
	if err != nil {
		return types.APIRetrieveResponse{}, err
	}
	endpoint := strings.TrimRight(gatewayURL, "/") + "/memory/retrieve"
	httpReq, err := http.NewRequest(http.MethodPost, endpoint, bytes.NewReader(body))
	if err != nil {
		return types.APIRetrieveResponse{}, err
	}
	httpReq.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(httpReq)
	if err != nil {
		return types.APIRetrieveResponse{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		data, _ := io.ReadAll(resp.Body)
		return types.APIRetrieveResponse{}, fmt.Errorf("gateway status %d: %s", resp.StatusCode, strings.TrimSpace(string(data)))
	}

	var out types.APIRetrieveResponse
	if err := json.NewDecoder(resp.Body).Decode(&out); err != nil {
		return types.APIRetrieveResponse{}, err
	}
	if out.SelectedID == "" || out.SelectionMode == "" || out.SourcePath == "" {
		return types.APIRetrieveResponse{}, errors.New("ERR_API_INPUT_INVALID: gateway response missing required retrieval fields")
	}
	return out, nil
}

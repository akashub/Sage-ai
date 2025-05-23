// // // backend/go/internal/llm/bridge.go
// package llm

// import (
// 	"bytes"
// 	"context"
// 	"encoding/json"
// 	"fmt"
// 	"net/http"
// 	"sage-ai-v2/pkg/logger"
// 	"time"
// )

// type Bridge struct {
//     baseURL    string
//     httpClient *http.Client
//     sessionID  string
// }

// type AnalyzeRequest struct {
//     Question string                 `json:"question"`
//     Schema   map[string]interface{} `json:"schema"`
// }

// type AnalyzeResponse struct {
//     Analysis map[string]interface{} `json:"analysis"`
//     Error    string                `json:"error,omitempty"`
// }

// type GenerateQueryRequest struct {
//     Analysis map[string]interface{} `json:"analysis"`
//     Schema   map[string]interface{} `json:"schema"`
// }

// type GenerateQueryResponse struct {
//     Query string `json:"query"`
//     Error string `json:"error,omitempty"`
// }

// func CreateBridge(baseURL string) *Bridge {
//     return &Bridge{
//         baseURL: baseURL,
//         httpClient: &http.Client{
//             Timeout: 120 * time.Second,  // Increase timeout to 120 seconds
//             Transport: &http.Transport{
//                 TLSHandshakeTimeout:   10 * time.Second,
//                 ResponseHeaderTimeout: 110 * time.Second,
//                 ExpectContinueTimeout: 1 * time.Second,
//                 MaxIdleConns:          100,
//                 MaxConnsPerHost:       100,
//                 IdleConnTimeout:       90 * time.Second,
//             },
//         },
//     }
// }

// type HealingRequest struct {
//     ValidationResult map[string]interface{} `json:"validation_result"`
//     OriginalQuery   string                 `json:"original_query"`
//     Analysis        map[string]interface{} `json:"analysis"`
//     Schema          map[string]interface{} `json:"schema"`
// }

// type HealingResult struct {
//     HealdQuery         string     `json:"healed_query"`
//     RequiresReanalysis bool       `json:"requires_reanalysis"`
//     Changes           []Change    `json:"changes_made"`
//     Confidence        float64     `json:"confidence"`
//     RequiresHumanReview bool      `json:"requires_human_review"`
//     Notes             string      `json:"notes"`
// }

// type Change struct {
//     Issue     string `json:"issue"`
//     Fix       string `json:"fix"`
//     Reasoning string `json:"reasoning"`
// }

// func (b *Bridge) SetSession(sessionID string) {
//     b.sessionID = sessionID
// }

// func (b *Bridge) Analyze(ctx context.Context, question string, schema map[string]interface{}) (map[string]interface{}, error) {
//     logger.InfoLogger.Printf("Bridge: Making analyze request")
//     logger.DebugLogger.Printf("Question: %s", question)
//     logger.DebugLogger.Printf("Schema: %+v", schema)

//     req := AnalyzeRequest{
//         Question: question,
//         Schema:   schema,
//     }

//     resp, err := b.MakeRequest(ctx, "/analyze", req, b.sessionID)
//     if err != nil {
//         logger.ErrorLogger.Printf("Analyze request failed: %v", err)
//         return nil, fmt.Errorf("error making analyze request: %w", err)
//     }

//     var analyzeResp AnalyzeResponse
//     if err := json.Unmarshal(resp, &analyzeResp); err != nil {
//         logger.ErrorLogger.Printf("Failed to unmarshal response: %v", err)
//         return nil, fmt.Errorf("error unmarshaling analyze response: %w", err)
//     }

//     if analyzeResp.Error != "" {
//         logger.ErrorLogger.Printf("LLM error: %s", analyzeResp.Error)
//         return nil, fmt.Errorf("LLM error: %s", analyzeResp.Error)
//     }

//     logger.InfoLogger.Printf("Analysis completed successfully")
//     logger.DebugLogger.Printf("Analysis result: %+v", analyzeResp.Analysis)
//     return analyzeResp.Analysis, nil
// }

// func (b *Bridge) GenerateQuery(ctx context.Context, analysis map[string]interface{}, schema map[string]interface{}) (string, error) {
//     logger.InfoLogger.Printf("Bridge: Making generate request")
//     logger.DebugLogger.Printf("Analysis: %+v", analysis)
//     logger.DebugLogger.Printf("Schema: %+v", schema)

//     requestData := map[string]interface{}{
//         "analysis": analysis,
//         "schema":  schema,
//     }

//     resp, err := b.MakeRequest(ctx, "/generate", requestData, b.sessionID)
//     if err != nil {
//         logger.ErrorLogger.Printf("Generate request failed: %v", err)
//         return "", fmt.Errorf("error making generate request: %w", err)
//     }

//     var genResponse struct {
//         Query string `json:"query"`
//     }
//     if err := json.Unmarshal(resp, &genResponse); err != nil {
//         logger.ErrorLogger.Printf("Failed to unmarshal response: %v", err)
//         return "", fmt.Errorf("error unmarshaling generate response: %w", err)
//     }

//     return genResponse.Query, nil
// }

// func (b *Bridge) ValidateQuery(ctx context.Context, query string, schema map[string]interface{}) (map[string]interface{}, error) {
//     requestData := map[string]interface{}{
//         "query":  query,
//         "schema": schema,
//     }

//     resp, err := b.MakeRequest(ctx, "/validate", requestData, b.sessionID)
//     if err != nil {
//         logger.ErrorLogger.Printf("Validate request failed: %v", err)
//         return nil, fmt.Errorf("error making validate request: %w", err)
//     }

//     var validation map[string]interface{}
//     if err := json.Unmarshal(resp, &validation); err != nil {
//         logger.ErrorLogger.Printf("Failed to unmarshal validation response: %v", err)
//         return nil, fmt.Errorf("error unmarshaling validation response: %w", err)
//     }

//     return validation, nil
// }

// // func (b *Bridge) MakeRequest(ctx context.Context, endpoint string, payload interface{}, sessionID string) ([]byte, error) {
// //     // Add session ID to payload
// //     requestData := map[string]interface{}{
// //         "session_id": sessionID,
// //         "data":      payload,
// //     }

// //     jsonData, err := json.Marshal(requestData)
// //     if err != nil {
// //         return nil, fmt.Errorf("error marshaling request: %w", err)
// //     }

// //     req, err := http.NewRequestWithContext(
// //         ctx,
// //         "POST",
// //         b.baseURL+endpoint,
// //         bytes.NewBuffer(jsonData),
// //     )
// //     if err != nil {
// //         return nil, fmt.Errorf("error creating request: %w", err)
// //     }

// //     req.Header.Set("Content-Type", "application/json")

// //     resp, err := b.httpClient.Do(req)
// //     if err != nil {
// //         return nil, fmt.Errorf("error making request: %w", err)
// //     }
// //     defer resp.Body.Close()

// //     var respBody bytes.Buffer
// //     _, err = respBody.ReadFrom(resp.Body)
// //     if err != nil {
// //         return nil, fmt.Errorf("error reading response: %w", err)
// //     }

// //     if resp.StatusCode != http.StatusOK {
// //         return nil, fmt.Errorf("error response from LLM service: %s", respBody.String())
// //     }

// //     return respBody.Bytes(), nil
// // }
// func (b *Bridge) MakeRequest(ctx context.Context, endpoint string, requestData interface{}, sessionID string) ([]byte, error) {
//     logger.InfoLogger.Printf("Making request to endpoint: %s", endpoint)
//     logger.DebugLogger.Printf("Request payload: %s", requestData)

//     ctx, cancel := context.WithTimeout(ctx, 110*time.Second)
//     defer cancel()

//     wrappedRequest := map[string]interface{}{
//         "session_id": sessionID,
//         "data":      requestData,
//     }

//     jsonData, err := json.Marshal(wrappedRequest)
//     if err != nil {
//         logger.ErrorLogger.Printf("Failed to marshal request: %v", err)
//         return nil, fmt.Errorf("error marshaling request: %w", err)
//     }

//     req, err := http.NewRequestWithContext(
//         ctx,
//         "POST",
//         b.baseURL+endpoint,
//         bytes.NewBuffer(jsonData),
//     )
//     if err != nil {
//         logger.ErrorLogger.Printf("Failed to create request: %v", err)
//         return nil, fmt.Errorf("error creating request: %w", err)
//     }

//     req.Header.Set("Content-Type", "application/json")

//     logger.InfoLogger.Printf("Sending request to %s", b.baseURL+endpoint)
//     resp, err := b.httpClient.Do(req)
//     if err != nil {
//         logger.ErrorLogger.Printf("Request failed: %v", err)
//         return nil, fmt.Errorf("error making request: %w", err)
//     }
//     defer resp.Body.Close()

//     var respBody bytes.Buffer
//     _, err = respBody.ReadFrom(resp.Body)
//     if err != nil {
//         logger.ErrorLogger.Printf("Failed to read response: %v", err)
//         return nil, fmt.Errorf("error reading response: %w", err)
//     }

//     if resp.StatusCode != http.StatusOK {
//         logger.ErrorLogger.Printf("Received non-200 status code: %d, body: %s", resp.StatusCode, respBody.String())
//         return nil, fmt.Errorf("error response from LLM service: %s", respBody.String())
//     }

//     logger.DebugLogger.Printf("Response received: %s", respBody.String())
//     return respBody.Bytes(), nil
// }

// func (b *Bridge) HealQuery(
//     ctx context.Context,
//     validationResult map[string]interface{},
//     originalQuery string,
//     analysis map[string]interface{},
//     schema map[string]interface{},
// ) (*HealingResult, error) {
//     req := HealingRequest{
//         ValidationResult: validationResult,
//         OriginalQuery:   originalQuery,
//         Analysis:        analysis,
//         Schema:          schema,
//     }

//     resp, err := b.MakeRequest(ctx, "/heal", req, b.sessionID)
//     if err != nil {
//         return nil, fmt.Errorf("error making healing request: %w", err)
//     }

//     var healingResult HealingResult
//     if err := json.Unmarshal(resp, &healingResult); err != nil {
//         return nil, fmt.Errorf("error unmarshaling healing response: %w", err)
//     }

//     return &healingResult, nil
// }

// backend/go/internal/llm/bridge.go
package llm

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"sage-ai-v2/pkg/logger"
	"time"
)

type Bridge struct {
	baseURL    string
	httpClient *http.Client
	sessionID  string
	llmConfig  *LLMConfig
}

type AnalyzeResponse struct {
	Analysis map[string]interface{} `json:"analysis"`
	Error    string                 `json:"error,omitempty"`
}

type GenerateQueryResponse struct {
	Query string `json:"query"`
	Error string `json:"error,omitempty"`
}

func (b *Bridge) GetSessionID() string {
    return b.sessionID
}

func (b *Bridge) SetLLMConfig(config *LLMConfig) {
    b.llmConfig = config
}

func CreateBridge(baseURL string) *Bridge {
	return &Bridge{
		baseURL: baseURL,
		httpClient: &http.Client{
			Timeout: 120 * time.Second,
			Transport: &http.Transport{
				TLSHandshakeTimeout:   10 * time.Second,
				ResponseHeaderTimeout: 110 * time.Second,
				ExpectContinueTimeout: 1 * time.Second,
				MaxIdleConns:          100,
				MaxConnsPerHost:       100,
				IdleConnTimeout:       90 * time.Second,
			},
		},
	}
}

type HealingRequest struct {
	ValidationResult map[string]interface{} `json:"validation_result"`
	OriginalQuery    string                 `json:"original_query"`
	Analysis         map[string]interface{} `json:"analysis"`
	Schema           map[string]interface{} `json:"schema"`
}

type HealingResult struct {
	HealdQuery          string   `json:"healed_query"`
	RequiresReanalysis  bool     `json:"requires_reanalysis"`
	Changes             []Change `json:"changes_made"`
	Confidence          float64  `json:"confidence"`
	RequiresHumanReview bool     `json:"requires_human_review"`
	Notes               string   `json:"notes"`
}

type Change struct {
	Issue     string `json:"issue"`
	Fix       string `json:"fix"`
	Reasoning string `json:"reasoning"`
}

func (b *Bridge) SetSession(sessionID string) {
	b.sessionID = sessionID
}

func (b *Bridge) Analyze(ctx context.Context, question string, schema map[string]interface{}) (map[string]interface{}, error) {
	logger.InfoLogger.Printf("Bridge: Making analyze request")
	logger.DebugLogger.Printf("Question: %s", question)
	logger.DebugLogger.Printf("Schema: %+v", schema)

	req := map[string]interface{}{
		"question": question,
		"schema":   schema,
	}

	resp, err := b.MakeRequest(ctx, "/analyze", req, b.sessionID)
	if err != nil {
		logger.ErrorLogger.Printf("Analyze request failed: %v", err)
		return nil, fmt.Errorf("error making analyze request: %w", err)
	}

	var analyzeResp AnalyzeResponse
	if err := json.Unmarshal(resp, &analyzeResp); err != nil {
		logger.ErrorLogger.Printf("Failed to unmarshal response: %v", err)
		return nil, fmt.Errorf("error unmarshaling analyze response: %w", err)
	}

	if analyzeResp.Error != "" {
		logger.ErrorLogger.Printf("LLM error: %s", analyzeResp.Error)
		return nil, fmt.Errorf("LLM error: %s", analyzeResp.Error)
	}

	logger.InfoLogger.Printf("Analysis completed successfully")
	logger.DebugLogger.Printf("Analysis result: %+v", analyzeResp.Analysis)
	return analyzeResp.Analysis, nil
}

// AnalyzeWithKnowledge is an enhanced version of Analyze that includes knowledge context
func (b *Bridge) AnalyzeWithKnowledge(ctx context.Context, analysisRequest map[string]interface{}) (map[string]interface{}, error) {
	logger.InfoLogger.Printf("Bridge: Making analyze with knowledge request")
	logger.DebugLogger.Printf("Analysis request: %+v", analysisRequest)

	resp, err := b.MakeRequest(ctx, "/analyze_with_knowledge", analysisRequest, b.sessionID)
	if err != nil {
		logger.ErrorLogger.Printf("Analyze with knowledge request failed: %v", err)
		return nil, fmt.Errorf("error making analyze request: %w", err)
	}

	var analyzeResp AnalyzeResponse
	if err := json.Unmarshal(resp, &analyzeResp); err != nil {
		logger.ErrorLogger.Printf("Failed to unmarshal response: %v", err)
		return nil, fmt.Errorf("error unmarshaling analyze response: %w", err)
	}

	if analyzeResp.Error != "" {
		logger.ErrorLogger.Printf("LLM error: %s", analyzeResp.Error)
		return nil, fmt.Errorf("LLM error: %s", analyzeResp.Error)
	}

	logger.InfoLogger.Printf("Analysis with knowledge completed successfully")
	logger.DebugLogger.Printf("Analysis result: %+v", analyzeResp.Analysis)
	return analyzeResp.Analysis, nil
}

func (b *Bridge) GenerateQuery(ctx context.Context, analysis map[string]interface{}, schema map[string]interface{}) (string, error) {
	logger.InfoLogger.Printf("Bridge: Making generate request")
	logger.DebugLogger.Printf("Analysis: %+v", analysis)
	logger.DebugLogger.Printf("Schema: %+v", schema)

	requestData := map[string]interface{}{
		"analysis": analysis,
		"schema":   schema,
	}

	resp, err := b.MakeRequest(ctx, "/generate", requestData, b.sessionID)
	if err != nil {
		logger.ErrorLogger.Printf("Generate request failed: %v", err)
		return "", fmt.Errorf("error making generate request: %w", err)
	}

	var genResponse struct {
		Query string `json:"query"`
	}
	if err := json.Unmarshal(resp, &genResponse); err != nil {
		logger.ErrorLogger.Printf("Failed to unmarshal response: %v", err)
		return "", fmt.Errorf("error unmarshaling generate response: %w", err)
	}

	return genResponse.Query, nil
}

// GenerateQueryWithKnowledge is an enhanced version of GenerateQuery that includes knowledge context
func (b *Bridge) GenerateQueryWithKnowledge(ctx context.Context, request map[string]interface{}) (string, error) {
	logger.InfoLogger.Printf("Bridge: Making generate with knowledge request")
	logger.DebugLogger.Printf("Generate request: %+v", request)

	if kc, ok := request["knowledge_context"]; ok {
		kcStr, _ := json.MarshalIndent(kc, "", "  ")
		logger.InfoLogger.Printf("Knowledge context: %s", string(kcStr))
	}

	resp, err := b.MakeRequest(ctx, "/generate_with_knowledge", request, b.sessionID)
	if err != nil {
		logger.ErrorLogger.Printf("Generate with knowledge request failed: %v", err)
		return "", fmt.Errorf("error making generate request: %w", err)
	}

	var genResponse GenerateQueryResponse
	if err := json.Unmarshal(resp, &genResponse); err != nil {
		logger.ErrorLogger.Printf("Failed to unmarshal response: %v", err)
		return "", fmt.Errorf("error unmarshaling generate response: %w", err)
	}

	if genResponse.Error != "" {
		logger.ErrorLogger.Printf("LLM error: %s", genResponse.Error)
		return "", fmt.Errorf("LLM error: %s", genResponse.Error)
	}

	logger.InfoLogger.Printf("Query generation with knowledge completed successfully")
	return genResponse.Query, nil
}

// MakeRequest makes an HTTP request to the Python service
func (b *Bridge) MakeRequest(ctx context.Context, endpoint string, requestData interface{}, sessionID string) ([]byte, error) {
	logger.InfoLogger.Printf("Making request to endpoint: %s", endpoint)
	logger.DebugLogger.Printf("Request payload: %v", requestData)

	ctx, cancel := context.WithTimeout(ctx, 110*time.Second)
	defer cancel()

	wrappedRequest := map[string]interface{}{
		"session_id": sessionID,
		"data":       requestData,
	}

	// Include LLM config if available
    if b.llmConfig != nil {
        wrappedRequest["llm_config"] = b.llmConfig
    }

	jsonData, err := json.Marshal(wrappedRequest)
	if err != nil {
		logger.ErrorLogger.Printf("Failed to marshal request: %v", err)
		return nil, fmt.Errorf("error marshaling request: %w", err)
	}

	req, err := http.NewRequestWithContext(
		ctx,
		"POST",
		b.baseURL+endpoint,
		bytes.NewBuffer(jsonData),
	)
	if err != nil {
		logger.ErrorLogger.Printf("Failed to create request: %v", err)
		return nil, fmt.Errorf("error creating request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")

	logger.InfoLogger.Printf("Sending request to %s", b.baseURL+endpoint)
	resp, err := b.httpClient.Do(req)
	if err != nil {
		logger.ErrorLogger.Printf("Request failed: %v", err)
		return nil, fmt.Errorf("error making request: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		logger.ErrorLogger.Printf("Failed to read response: %v", err)
		return nil, fmt.Errorf("error reading response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		logger.ErrorLogger.Printf("Received non-200 status code: %d, body: %s", resp.StatusCode, respBody)
		return nil, fmt.Errorf("error response from LLM service: %s", respBody)
	}

	logger.DebugLogger.Printf("Response received: %s", respBody)
	return respBody, nil
}

func (b *Bridge) ValidateQuery(ctx context.Context, query string, schema map[string]interface{}) (map[string]interface{}, error) {
	requestData := map[string]interface{}{
		"query":  query,
		"schema": schema,
	}

	resp, err := b.MakeRequest(ctx, "/validate", requestData, b.sessionID)
	if err != nil {
		logger.ErrorLogger.Printf("Validate request failed: %v", err)
		return nil, fmt.Errorf("error making validate request: %w", err)
	}

	var validation map[string]interface{}
	if err := json.Unmarshal(resp, &validation); err != nil {
		logger.ErrorLogger.Printf("Failed to unmarshal validation response: %v", err)
		return nil, fmt.Errorf("error unmarshaling validation response: %w", err)
	}

	return validation, nil
}

func (b *Bridge) HealQuery(
	ctx context.Context,
	validationResult map[string]interface{},
	originalQuery string,
	analysis map[string]interface{},
	schema map[string]interface{},
) (*HealingResult, error) {
	req := HealingRequest{
		ValidationResult: validationResult,
		OriginalQuery:    originalQuery,
		Analysis:         analysis,
		Schema:           schema,
	}

	resp, err := b.MakeRequest(ctx, "/heal", req, b.sessionID)
	if err != nil {
		return nil, fmt.Errorf("error making healing request: %w", err)
	}

	var healingResult HealingResult
	if err := json.Unmarshal(resp, &healingResult); err != nil {
		return nil, fmt.Errorf("error unmarshaling healing response: %w", err)
	}

	return &healingResult, nil
}

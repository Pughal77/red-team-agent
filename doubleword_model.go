package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"iter"
	"net/http"

	"google.golang.org/adk/v2/model"
	"google.golang.org/genai"
)

// DoublewordModel implements the model.LLM interface from Google ADK.
type DoublewordModel struct {
	modelName string
	apiKey    string
}

// NewDoublewordModel creates a new instance of DoublewordModel.
func NewDoublewordModel(modelName, apiKey string) *DoublewordModel {
	return &DoublewordModel{
		modelName: modelName,
		apiKey:    apiKey,
	}
}

// Name returns the name of the model.
func (m *DoublewordModel) Name() string {
	return m.modelName
}

// GenerateContent handles calling the Doubleword completions API.
func (m *DoublewordModel) GenerateContent(ctx context.Context, req *model.LLMRequest, stream bool) iter.Seq2[*model.LLMResponse, error] {
	return func(yield func(*model.LLMResponse, error) bool) {
		// Translate Contents to OpenAI messages format
		var messages []map[string]any
		for _, c := range req.Contents {
			role := c.Role
			if role == "model" {
				role = "assistant"
			} else if role == "" {
				role = "user"
			}

			var content string
			var toolCalls []map[string]any
			for _, p := range c.Parts {
				if p.Text != "" {
					content += p.Text
				}
				if p.FunctionCall != nil {
					toolCalls = append(toolCalls, map[string]any{
						"id":   "call_" + p.FunctionCall.Name,
						"type": "function",
						"function": map[string]any{
							"name":      p.FunctionCall.Name,
							"arguments": serializeMap(p.FunctionCall.Args),
						},
					})
				}
				if p.FunctionResponse != nil {
					// Format function execution result back to LLM context
					toolMsg := map[string]any{
						"role":         "tool",
						"tool_call_id": "call_" + p.FunctionResponse.Name,
						"content":      serializeResponse(p.FunctionResponse.Response),
					}
					messages = append(messages, toolMsg)
				}
			}

			// Avoid sending duplicates of tool response message structures
			var hasToolResponse bool
			for _, p := range c.Parts {
				if p.FunctionResponse != nil {
					hasToolResponse = true
					break
				}
			}
			if !hasToolResponse {
				msg := map[string]any{
					"role":    role,
					"content": content,
				}
				if len(toolCalls) > 0 {
					msg["tool_calls"] = toolCalls
				}
				messages = append(messages, msg)
			}
		}

		// Translate Tools to OpenAI tool declarations
		var openAITools []map[string]any
		if req.Config != nil && len(req.Config.Tools) > 0 {
			for _, t := range req.Config.Tools {
				if t.FunctionDeclarations != nil {
					for _, fd := range t.FunctionDeclarations {
						openAITools = append(openAITools, map[string]any{
							"type": "function",
							"function": map[string]any{
								"name":        fd.Name,
								"description": fd.Description,
								"parameters":  fd.Parameters,
							},
						})
					}
				}
			}
		}

		// Build chat completions request body
		reqBodyObj := map[string]any{
			"model":    m.modelName,
			"messages": messages,
		}
		if len(openAITools) > 0 {
			reqBodyObj["tools"] = openAITools
		}

		reqBodyBytes, err := json.Marshal(reqBodyObj)
		if err != nil {
			yield(nil, fmt.Errorf("failed to marshal request payload: %w", err))
			return
		}

		// Execute HTTP POST to Doubleword real-time API endpoint
		httpReq, err := http.NewRequestWithContext(ctx, "POST", "https://api.doubleword.ai/v1/chat/completions", bytes.NewReader(reqBodyBytes))
		if err != nil {
			yield(nil, fmt.Errorf("failed to create http request: %w", err))
			return
		}
		httpReq.Header.Set("Content-Type", "application/json")
		httpReq.Header.Set("Authorization", "Bearer "+m.apiKey)

		resp, err := http.DefaultClient.Do(httpReq)
		if err != nil {
			yield(nil, fmt.Errorf("http request failed: %w", err))
			return
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			body, _ := io.ReadAll(resp.Body)
			yield(nil, fmt.Errorf("api error (status %d): %s", resp.StatusCode, string(body)))
			return
		}

		// Parse the JSON response
		var completionResponse struct {
			Choices []struct {
				Message struct {
					Content   string `json:"content"`
					ToolCalls []struct {
						ID       string `json:"id"`
						Type     string `json:"type"`
						Function struct {
							Name      string `json:"name"`
							Arguments string `json:"arguments"`
						} `json:"function"`
					} `json:"tool_calls"`
				} `json:"message"`
				FinishReason string `json:"finish_reason"`
			} `json:"choices"`
		}

		if err := json.NewDecoder(resp.Body).Decode(&completionResponse); err != nil {
			yield(nil, fmt.Errorf("failed to decode response: %w", err))
			return
		}

		if len(completionResponse.Choices) == 0 {
			yield(nil, fmt.Errorf("empty choices in response"))
			return
		}

		choice := completionResponse.Choices[0]

		// Build genai.Part slices to represent response content
		var parts []*genai.Part
		if choice.Message.Content != "" {
			parts = append(parts, &genai.Part{Text: choice.Message.Content})
		}
		for _, tc := range choice.Message.ToolCalls {
			var args map[string]any
			_ = json.Unmarshal([]byte(tc.Function.Arguments), &args)
			parts = append(parts, &genai.Part{
				FunctionCall: &genai.FunctionCall{
					Name: tc.Function.Name,
					Args: args,
				},
			})
		}

		llmResponse := &model.LLMResponse{
			Content: &genai.Content{
				Role:  "model",
				Parts: parts,
			},
		}

		yield(llmResponse, nil)
	}
}

func serializeMap(m map[string]any) string {
	b, _ := json.Marshal(m)
	return string(b)
}

func serializeResponse(val any) string {
	switch v := val.(type) {
	case string:
		return v
	case []byte:
		return string(v)
	default:
		b, err := json.Marshal(v)
		if err != nil {
			return fmt.Sprintf("%v", v)
		}
		return string(b)
	}
}

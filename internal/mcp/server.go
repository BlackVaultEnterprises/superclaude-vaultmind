package mcp

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/gorilla/websocket"
	"github.com/opencode-ai/opencode/internal/logging"
	"github.com/opencode-ai/opencode/internal/superclaude"
)

// MCPServer implements the Model Context Protocol server
type MCPServer struct {
	upgrader websocket.Upgrader
	handler  *superclaude.SuperClaudeHandler
	sessions sync.Map
	mu       sync.RWMutex
}

// NewMCPServer creates a new MCP server
func NewMCPServer(handler *superclaude.SuperClaudeHandler) *MCPServer {
	return &MCPServer{
		upgrader: websocket.Upgrader{
			CheckOrigin: func(r *http.Request) bool {
				// Allow connections from any origin for development
				// TODO: Restrict in production
				return true
			},
		},
		handler: handler,
	}
}

// MCPRequest represents an incoming MCP request
type MCPRequest struct {
	ID      string          `json:"id"`
	Method  string          `json:"method"`
	Params  json.RawMessage `json:"params"`
	Context MCPContext      `json:"context"`
}

// MCPContext provides execution context
type MCPContext struct {
	SessionID   string            `json:"session_id"`
	WorkingDir  string            `json:"working_dir"`
	Environment map[string]string `json:"environment"`
	Capabilities []string         `json:"capabilities"`
}

// MCPResponse represents an MCP response
type MCPResponse struct {
	ID     string      `json:"id"`
	Result interface{} `json:"result,omitempty"`
	Error  *MCPError   `json:"error,omitempty"`
}

// MCPError represents an error response
type MCPError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

// ServeHTTP handles WebSocket connections
func (s *MCPServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	conn, err := s.upgrader.Upgrade(w, r, nil)
	if err != nil {
		logging.Error("Failed to upgrade connection", "error", err)
		return
	}
	defer conn.Close()

	sessionID := r.Header.Get("X-Session-ID")
	if sessionID == "" {
		sessionID = generateSessionID()
	}

	logging.Info("New MCP connection", "session_id", sessionID)
	
	// Handle messages
	for {
		var req MCPRequest
		err := conn.ReadJSON(&req)
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				logging.Error("WebSocket error", "error", err)
			}
			break
		}

		// Process request
		resp := s.handleRequest(req)
		
		// Send response
		if err := conn.WriteJSON(resp); err != nil {
			logging.Error("Failed to write response", "error", err)
			break
		}
	}
}

// handleRequest processes an MCP request
func (s *MCPServer) handleRequest(req MCPRequest) MCPResponse {
	switch req.Method {
	case "initialize":
		return s.handleInitialize(req)
	case "execute":
		return s.handleExecute(req)
	case "complete":
		return s.handleComplete(req)
	case "analyze":
		return s.handleAnalyze(req)
	case "capabilities":
		return s.handleCapabilities(req)
	default:
		return MCPResponse{
			ID: req.ID,
			Error: &MCPError{
				Code:    -32601,
				Message: fmt.Sprintf("Method not found: %s", req.Method),
			},
		}
	}
}

// handleInitialize initializes a new MCP session
func (s *MCPServer) handleInitialize(req MCPRequest) MCPResponse {
	var params struct {
		Name    string `json:"name"`
		Version string `json:"version"`
	}
	
	if err := json.Unmarshal(req.Params, &params); err != nil {
		return errorResponse(req.ID, -32602, "Invalid params")
	}

	// Store session info
	s.sessions.Store(req.Context.SessionID, &MCPSession{
		ID:          req.Context.SessionID,
		WorkingDir:  req.Context.WorkingDir,
		Environment: req.Context.Environment,
	})

	return MCPResponse{
		ID: req.ID,
		Result: map[string]interface{}{
			"capabilities": []string{
				"superclaude.commands",
				"superclaude.personas",
				"file.read",
				"file.write",
				"bash.execute",
				"git.operations",
			},
			"version": "1.0.0",
		},
	}
}

// handleExecute executes a SuperClaude command
func (s *MCPServer) handleExecute(req MCPRequest) MCPResponse {
	var params struct {
		Command string `json:"command"`
		Input   string `json:"input"`
	}
	
	if err := json.Unmarshal(req.Params, &params); err != nil {
		return errorResponse(req.ID, -32602, "Invalid params")
	}

	// Execute SuperClaude command
	ctx := context.Background()
	handled, err := s.handler.HandleCommand(ctx, req.Context.SessionID, params.Command)
	
	if err != nil {
		return errorResponse(req.ID, -32603, err.Error())
	}

	if !handled {
		return errorResponse(req.ID, -32604, "Not a SuperClaude command")
	}

	return MCPResponse{
		ID: req.ID,
		Result: map[string]interface{}{
			"status": "success",
			"command": params.Command,
		},
	}
}

// handleComplete provides command completion
func (s *MCPServer) handleComplete(req MCPRequest) MCPResponse {
	var params struct {
		Input  string `json:"input"`
		Cursor int    `json:"cursor"`
	}
	
	if err := json.Unmarshal(req.Params, &params); err != nil {
		return errorResponse(req.ID, -32602, "Invalid params")
	}

	// Generate completions
	completions := generateCompletions(params.Input, params.Cursor)

	return MCPResponse{
		ID: req.ID,
		Result: map[string]interface{}{
			"completions": completions,
		},
	}
}

// handleAnalyze analyzes code or project
func (s *MCPServer) handleAnalyze(req MCPRequest) MCPResponse {
	var params struct {
		Path  string   `json:"path"`
		Types []string `json:"types"`
	}
	
	if err := json.Unmarshal(req.Params, &params); err != nil {
		return errorResponse(req.ID, -32602, "Invalid params")
	}

	// Run analysis using SuperClaude
	ctx := context.Background()
	command := fmt.Sprintf("/user:analyze %s", params.Path)
	
	handled, err := s.handler.HandleCommand(ctx, req.Context.SessionID, command)
	if err != nil {
		return errorResponse(req.ID, -32603, err.Error())
	}

	if !handled {
		return errorResponse(req.ID, -32604, "Analysis failed")
	}

	return MCPResponse{
		ID: req.ID,
		Result: map[string]interface{}{
			"status": "analyzing",
			"path":   params.Path,
		},
	}
}

// handleCapabilities returns server capabilities
func (s *MCPServer) handleCapabilities(req MCPRequest) MCPResponse {
	return MCPResponse{
		ID: req.ID,
		Result: map[string]interface{}{
			"commands": superclaude.GetAvailableCommands(),
			"personas": superclaude.GetAvailablePersonas(),
			"flags":    superclaude.GetAvailableFlags(),
		},
	}
}

// MCPSession represents an active MCP session
type MCPSession struct {
	ID          string
	WorkingDir  string
	Environment map[string]string
}

// Helper functions

func errorResponse(id string, code int, message string) MCPResponse {
	return MCPResponse{
		ID: id,
		Error: &MCPError{
			Code:    code,
			Message: message,
		},
	}
}

func generateSessionID() string {
	// Simple session ID generation
	return fmt.Sprintf("mcp-%d", time.Now().UnixNano())
}

func generateCompletions(input string, cursor int) []string {
	// Simple completion logic
	if strings.HasPrefix(input, "/user:") {
		return []string{
			"/user:analyze",
			"/user:build",
			"/user:test",
			"/user:improve",
			"/user:design",
			"/user:review",
			"/user:deploy",
			"/user:document",
		}
	}
	
	if strings.HasPrefix(input, "/persona:") {
		return []string{
			"/persona:architect",
			"/persona:frontend",
			"/persona:backend",
			"/persona:security",
			"/persona:qa",
			"/persona:performance",
		}
	}
	
	return []string{}
}
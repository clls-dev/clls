package lsp

type InitializeResult struct {
	Capabilities *ServerCapabilities         `json:"capabilities,omitempty"`
	ServerInfo   *InitializeResultServerInfo `json:"serverInfo,omitempty"`
}

type InitializeResultServerInfo struct {
	Name    string `json:"name"`
	Version string `json:"version,omitempty"`
}

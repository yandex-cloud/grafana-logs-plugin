package plugin

type loggingRequest struct {
	GroupID          string             `json:"groupId"`
	Limit            int                `json:"limit"`
	QueryText        string             `json:"queryText"`
	Levels           []Level            `json:"levels"`
	Stream           string             `json:"stream"`
	ResourceType     string             `json:"resourceType"`
	ResourceIDs      []string           `json:"resourceIds"`
	AddPayloadFields []string           `json:"addPayloadFields"`
	DerivedFields    []derivedFieldRule `json:"derivedFields"`
}

const apiKeyJsonInSettings = "apiKeyJson"

type loggingConfig struct {
	APIEndpoint  string              `json:"apiEndpoint"`
	FolderID     string              `json:"folderId"`
	DerivedLinks []derivedLinkConfig `json:"derivedLinks"`
}

type derivedLinkConfig struct {
	Field       string `json:"field"`
	Title       string `json:"title"`
	URL         string `json:"url"`
	TargetBlank bool   `json:"targetBlank"`
}

type suggestQueryRequest struct {
	GroupID      string `json:"groupId"`
	ResourceType string `json:"resourceType"`
}

type suggestQueryResponse struct {
	Groups        []string `json:"groups"`
	ResourceTypes []string `json:"resourceTypes"`
	ResourceIDs   []string `json:"resourceIds"`
}

type resourceError struct {
	Error string `json:"error"`
}

type derivedFieldRule struct {
	Name     string `json:"name"`
	Template string `json:"template"`
}

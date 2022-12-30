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
	APIEndpoint string `json:"apiEndpoint"`
	FolderID    string `json:"folderId"`
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

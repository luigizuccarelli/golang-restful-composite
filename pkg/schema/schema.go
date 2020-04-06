package schema

// SchemaInterface - acts as an interface wrapper
// All the go microservices will using this schema
type SchemaInterface struct {
	LastUpdate    int64       `json:"lastupdate,omitempty"`
	MetaInfo      string      `json:"metainfo,omitempty"`
	Strategy      string      `json:"strategy"`
	Requests      []Composite `json:"requests"`
	MergedContent []string    `json:"mergedcontent"`
}

// RequestSchema
type RequestSchema struct {
	MetaInfo string      `json:"metainfo,omitempty"`
	Strategy string      `json:"strategy"`
	Requests []Composite `json:"requests"`
}

// Composite schema

// Composite schema
// The strategy is as follows
// failonce - fail all if one of the results fail (default)
// failnone - don't fail even if there is a failure
// passonce - pass if at least one of the results pass
type Composite struct {
	Method     string `json:"method"`
	Url        string `json:"url"`
	Payload    string `json:"payload,omitemty"`
	Headers    string `json:"headers,omitemty"`
	Message    string `json:"message,omitempty"`
	StatusCode string `json:"statuscode,omitempty"`
	Status     string `json:"status,omitempty"`
}

// Response schema
type Response struct {
	Message    string          `json:"message,omitempty"`
	StatusCode string          `json:"statuscode,omitempty"`
	Status     string          `json:"status,omitempty"`
	Results    SchemaInterface `json:"results"`
}

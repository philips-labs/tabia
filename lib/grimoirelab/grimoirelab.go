package grimoirelab

// Projects holds all projects to be loaded in Grimoirelab
type Projects map[string]*Project

// Project holds the project resources and metadata
type Project struct {
	Metadata   Metadata `json:"meta,omitempty"`
	Git        []string `json:"git,omitempty"`
	Github     []string `json:"github,omitempty"`
	GithubRepo []string `json:"github:repo,omitempty"`
}

// Metadata hold metadata for a given project
type Metadata map[string]string

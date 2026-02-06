package tutorial

// Tutorial is the unified data model that accepts both tutorial patterns.
type Tutorial struct {
	// Common fields
	Name        string `yaml:"name"`
	Title       string `yaml:"title"`
	Description string `yaml:"description"`

	// Pattern A fields (narrative, instructions block with === delimiters)
	Instructions string `yaml:"instructions"`

	// Pattern B fields (structured modules with checkpoints)
	Duration      int      `yaml:"duration"`
	Difficulty    string   `yaml:"difficulty"`
	Prerequisites []string `yaml:"prerequisites"`
	Modules       []Module `yaml:"modules"`
	Completion    *CompletionMessage `yaml:"completion"`

	// Parsed pages (populated by loader)
	Pages []Page `yaml:"-"`
}

// Module represents a structured tutorial module (Pattern B).
type Module struct {
	ID          string       `yaml:"id"`
	Title       string       `yaml:"title"`
	Objectives  []string     `yaml:"objectives"`
	Checkpoints []Checkpoint `yaml:"checkpoints"`
}

// Checkpoint represents a validation checkpoint.
type Checkpoint struct {
	Name        string     `yaml:"name"`
	Description string     `yaml:"description"`
	Validation  Validation `yaml:"validation"`
}

// Validation defines how to validate a checkpoint.
type Validation struct {
	Type         string   `yaml:"type"` // pod_ready, deployment_ready, manual
	Selector     string   `yaml:"selector"`
	Namespace    string   `yaml:"namespace"`
	Deployment   string   `yaml:"deployment"`
	Deployments  []string `yaml:"deployments"`
	Instructions string   `yaml:"instructions"`
}

// CompletionMessage is shown when the tutorial is done.
type CompletionMessage struct {
	Message string `yaml:"message"`
}

// Page is a single rendered page in the TUI.
type Page struct {
	Title       string
	Content     string
	Checkpoints []Checkpoint
}

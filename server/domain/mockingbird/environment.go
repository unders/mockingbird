package mockingbird

// Env defines the different environment of mockingbird
type Env string

// Specifies the environments
const (
	DEV  Env = "dev"
	STAG     = "stag"
	PROD     = "prod"
)

package workflow

type Job struct {
	ID      string
	Type    string
	Payload map[string]any
}

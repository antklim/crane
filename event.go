package crane

// Event describes an event sent to crane to trigger an action.
// Crane uses event category to find relevant event handler.
type Event struct {
	Category string `json:"category"`
	Commit   string `json:"commit"`
}

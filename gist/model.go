package gist

// Response is the response coming back from gist API.
type Response struct {
	ID          string `json:"id"`
	URL         string `json:"url"`
	Description string `json:"description"`
	CreatedAt   string `json:"created_at"`
	UpdatedAt   string `json:"updated_at"`
}

// Gist represents one gist.
type Gist struct {
	Files map[string]File `json:"files"`
}

// File is one file in a Gist.
type File struct {
	Content string `json:"content"`
}

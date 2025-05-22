package utils

import gonanoid "github.com/matoous/go-nanoid/v2"

// GenerateID generates a unique 21-character NanoID
func GenerateID() string {
	id, err := gonanoid.New()
	if err != nil {
		return "" // In practice this should never happen
	}
	return id
}

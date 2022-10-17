package lib

import (
	"log"
	"testing"
)

func TestGenerationPath(t *testing.T) {
	path := GenerationPath("kyf", "txt")
	log.Print(path)
}

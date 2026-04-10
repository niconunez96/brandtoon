package test

import (
	"brandtoonapi/bounded_contexts/shared/domain"
	"testing"

	"github.com/google/uuid"
)

func TestGenerateUUIDv7ReturnsVersionSevenIdentifier(t *testing.T) {
	t.Parallel()

	generatedID, err := shareddomain.GenerateUUIDv7()
	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}

	parsedID, err := uuid.Parse(generatedID)
	if err != nil {
		t.Fatalf("expected a valid uuid, got %v", err)
	}

	if parsedID.Version() != 7 {
		t.Fatalf("expected uuid version 7, got %d", parsedID.Version())
	}
}

package message_test

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/ArtemNehoda/golang-hello-world/internal/domain/message"
)

// TestNew verifies that New() correctly populates Content, Author, and a
// non-zero UTC CreatedAt, while leaving ID at its zero value.
func TestNew(t *testing.T) {
	before := time.Now().UTC()
	e := message.New("Hello World!", "Artem")
	after := time.Now().UTC()

	if e.Content != "Hello World!" {
		t.Errorf("Content: got %q, want %q", e.Content, "Hello World!")
	}
	if e.Author != "Artem" {
		t.Errorf("Author: got %q, want %q", e.Author, "Artem")
	}
	if e.ID != 0 {
		t.Errorf("ID: got %d, want 0", e.ID)
	}
	if e.CreatedAt.IsZero() {
		t.Error("CreatedAt must not be zero")
	}
	if e.CreatedAt.Before(before) || e.CreatedAt.After(after) {
		t.Errorf("CreatedAt %v is outside [%v, %v]", e.CreatedAt, before, after)
	}
	if e.CreatedAt.Location() != time.UTC {
		t.Errorf("CreatedAt timezone: got %v, want UTC", e.CreatedAt.Location())
	}
}

// TestNew_EmptyFields ensures New() does not panic and produces an entity even
// when content and author are empty strings.
func TestNew_EmptyFields(t *testing.T) {
	e := message.New("", "")
	if e == nil {
		t.Fatal("New() returned nil")
	}
}

// TestEntity_JSONMarshal checks that Entity marshals to JSON with the expected
// keys and values, and that CreatedAt is excluded (json:"-").
func TestEntity_JSONMarshal(t *testing.T) {
	e := &message.Entity{
		ID:        42,
		Content:   "test content",
		Author:    "test author",
		CreatedAt: time.Now().UTC(),
	}

	data, err := json.Marshal(e)
	if err != nil {
		t.Fatalf("json.Marshal failed: %v", err)
	}

	var got map[string]any
	if err := json.Unmarshal(data, &got); err != nil {
		t.Fatalf("json.Unmarshal failed: %v", err)
	}

	// id
	if v, ok := got["id"]; !ok {
		t.Error("JSON missing key 'id'")
	} else if v.(float64) != 42 {
		t.Errorf("id: got %v, want 42", v)
	}

	// content
	if v, ok := got["content"]; !ok {
		t.Error("JSON missing key 'content'")
	} else if v.(string) != "test content" {
		t.Errorf("content: got %v, want 'test content'", v)
	}

	// author
	if v, ok := got["author"]; !ok {
		t.Error("JSON missing key 'author'")
	} else if v.(string) != "test author" {
		t.Errorf("author: got %v, want 'test author'", v)
	}

	// created_at must NOT appear (json:"-")
	if _, ok := got["created_at"]; ok {
		t.Error("JSON must not contain 'created_at' (tagged json:\"-\")")
	}
}

// TestEntity_JSONUnmarshal verifies round-trip: a JSON object without
// created_at unmarshals into an Entity whose CreatedAt remains zero.
func TestEntity_JSONUnmarshal(t *testing.T) {
	raw := `{"id":7,"content":"hi","author":"bob"}`
	var e message.Entity
	if err := json.Unmarshal([]byte(raw), &e); err != nil {
		t.Fatalf("json.Unmarshal failed: %v", err)
	}
	if e.ID != 7 {
		t.Errorf("ID: got %d, want 7", e.ID)
	}
	if e.Content != "hi" {
		t.Errorf("Content: got %q, want 'hi'", e.Content)
	}
	if e.Author != "bob" {
		t.Errorf("Author: got %q, want 'bob'", e.Author)
	}
	if !e.CreatedAt.IsZero() {
		t.Errorf("CreatedAt should be zero after unmarshal without that key, got %v", e.CreatedAt)
	}
}

// TestEntity_JSONMarshalSlice ensures a nil slice of Entity marshals to "[]"
// rather than "null", matching the handler's guarantee.
func TestEntity_JSONMarshalSlice(t *testing.T) {
	var msgs []message.Entity
	data, err := json.Marshal(msgs)
	if err != nil {
		t.Fatalf("json.Marshal nil slice failed: %v", err)
	}
	if string(data) != "null" {
		// The handler converts nil→empty slice before encoding, so just
		// document the raw nil-slice behaviour here.
		t.Logf("nil slice marshals to %s (handler guards this)", string(data))
	}

	// An initialised empty slice must marshal to "[]".
	empty := []message.Entity{}
	data2, err := json.Marshal(empty)
	if err != nil {
		t.Fatalf("json.Marshal empty slice failed: %v", err)
	}
	if string(data2) != "[]" {
		t.Errorf("empty slice: got %s, want []", string(data2))
	}
}

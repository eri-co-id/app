package file

import (
	"strings"
	"testing"
)

func TestUploadKeySanitizesPathTraversal(t *testing.T) {
	cases := []string{"../../etc/passwd", "../../../x", "a/b/c.jpg", ".."}
	for _, name := range cases {
		key, err := uploadKey("riding", 5, name)
		if err != nil {
			t.Fatalf("unexpected error for %q: %v", name, err)
		}
		if strings.Contains(key, "/..") || strings.Contains(key, "../") {
			t.Fatalf("key %q derived from %q still contains path traversal", key, name)
		}
	}
}

func TestUploadKeyIncludesObjectAndRecord(t *testing.T) {
	key, err := uploadKey("riding", 5, "photo.jpg")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.HasPrefix(key, "riding/5/") {
		t.Fatalf("expected key to start with riding/5/, got %q", key)
	}
	if !strings.HasSuffix(key, "-photo.jpg") {
		t.Fatalf("expected key to end with the original filename, got %q", key)
	}
}

package output

import (
	"fmt"
	"testing"
)

func TestBrandedIDString(t *testing.T) {
	t.Parallel()

	id := NewBrandedID[GraphNodeIDBrand](testGraphNodeID)
	if got := id.String(); got != testGraphNodeID {
		t.Errorf("String() = %q, want %q", got, testGraphNodeID)
	}

	empty := NewBrandedID[GraphNodeIDBrand]("")
	if got := empty.String(); got != "" {
		t.Errorf("String() = %q, want empty", got)
	}
}

func TestBrandedIDMarshalText(t *testing.T) {
	t.Parallel()

	id := NewBrandedID[GraphNodeIDBrand]("hello")

	data, err := id.MarshalText()
	if err != nil {
		t.Fatalf("MarshalText() error = %v", err)
	}

	if string(data) != "hello" {
		t.Errorf("MarshalText() = %q, want %q", string(data), "hello")
	}
}

func TestBrandedIDUnmarshalText(t *testing.T) {
	t.Parallel()

	var id BrandedID[GraphNodeIDBrand]
	if err := id.UnmarshalText([]byte("world")); err != nil {
		t.Fatalf("UnmarshalText() error = %v", err)
	}

	if got := id.Get(); got != "world" {
		t.Errorf("after UnmarshalText, Get() = %q, want %q", got, "world")
	}
}

func TestBrandedIDFormat(t *testing.T) {
	t.Parallel()

	id := NewBrandedID[GraphNodeIDBrand](testGraphNodeID)

	t.Run("%s", func(t *testing.T) {
		t.Parallel()

		got := fmt.Sprintf("%s", id)
		if got != testGraphNodeID {
			t.Errorf("%%s = %q, want %q", got, testGraphNodeID)
		}
	})

	t.Run("%#v", func(t *testing.T) {
		t.Parallel()

		got := fmt.Sprintf("%#v", id)
		want := `BrandedID{"test-id"}`

		if got != want {
			t.Errorf("%%#v = %q, want %q", got, want)
		}
	})

	t.Run("%v", func(t *testing.T) {
		t.Parallel()

		got := fmt.Sprintf("%v", id)
		if got != testGraphNodeID {
			t.Errorf("%%v = %q, want %q", got, testGraphNodeID)
		}
	})
}

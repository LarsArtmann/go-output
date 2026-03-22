package output

import (
	"strings"
	"testing"
)

func TestD2Diagram(t *testing.T) {
	t.Run("basic diagram", func(t *testing.T) {
		d := NewD2Diagram()
		d.AddTable("users", []D2Column{
			{Name: "id", Type: "int"},
			{Name: "name", Type: "string"},
		})

		got := d.Render()

		if !strings.Contains(got, "users:") {
			t.Error("Render() should contain table name")
		}
		if !strings.Contains(got, "id: int") {
			t.Error("Render() should contain column definitions")
		}
	})

	t.Run("chaining", func(t *testing.T) {
		d := NewD2Diagram().
			AddTable("users", []D2Column{
				{Name: "id", Type: "int"},
			})

		if d == nil {
			t.Error("Method chaining should return non-nil")
		}
	})

	t.Run("multiple tables", func(t *testing.T) {
		d := NewD2Diagram()
		d.AddTable("users", []D2Column{{Name: "id", Type: "int"}})
		d.AddTable("posts", []D2Column{{Name: "id", Type: "int"}})

		got := d.Render()

		if !strings.Contains(got, "users:") {
			t.Error("Render() should contain users table")
		}
		if !strings.Contains(got, "posts:") {
			t.Error("Render() should contain posts table")
		}
	})
}

func TestNewD2Diagram(t *testing.T) {
	d := NewD2Diagram()
	// Verify diagram is initialized properly
	_ = d.tables // Just ensure field is accessible
}

package testhelpers

import (
	"testing"
)

func TestErrorWriter(t *testing.T) {
	t.Parallel()

	w := &ErrorWriter{}
	n, err := w.Write([]byte("data"))

	if n != 0 {
		t.Errorf("ErrorWriter.Write() n = %d, want 0", n)
	}

	if err != ErrWrite {
		t.Errorf("ErrorWriter.Write() err = %v, want ErrWrite", err)
	}
}

func TestWriteNThenFailWriter(t *testing.T) {
	t.Parallel()

	t.Run("writes successfully while remaining", func(t *testing.T) {
		t.Parallel()

		w := &WriteNThenFailWriter{Remaining: 2}

		n, err := w.Write([]byte("data"))
		if err != nil {
			t.Errorf("Write() err = %v, want nil", err)
		}

		if n != 4 {
			t.Errorf("Write() n = %d, want 4", n)
		}

		if w.Remaining != 1 {
			t.Errorf("Remaining = %d, want 1", w.Remaining)
		}
	})

	t.Run("fails when exhausted", func(t *testing.T) {
		t.Parallel()

		w := &WriteNThenFailWriter{Remaining: 0}

		n, err := w.Write([]byte("data"))
		if n != 0 {
			t.Errorf("Write() n = %d, want 0", n)
		}

		if err != ErrWrite {
			t.Errorf("Write() err = %v, want ErrWrite", err)
		}
	})

	t.Run("decrements to zero then fails", func(t *testing.T) {
		t.Parallel()

		w := &WriteNThenFailWriter{Remaining: 1}

		n, err := w.Write([]byte("first"))
		if err != nil {
			t.Errorf("first Write() err = %v, want nil", err)
		}

		if n != 5 {
			t.Errorf("first Write() n = %d, want 5", n)
		}

		n, err = w.Write([]byte("second"))
		if n != 0 {
			t.Errorf("second Write() n = %d, want 0", n)
		}

		if err != ErrWrite {
			t.Errorf("second Write() err = %v, want ErrWrite", err)
		}
	})
}

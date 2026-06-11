package testhelpers

import (
	"errors"
	"testing"
)

func assertErrWrite(t *testing.T, n int, err error) {
	t.Helper()

	if n != 0 {
		t.Errorf("Write() n = %d, want 0", n)
	}

	if !errors.Is(err, ErrWrite) {
		t.Errorf("Write() err = %v, want ErrWrite", err)
	}
}

func TestErrorWriter(t *testing.T) {
	t.Parallel()

	w := &ErrorWriter{}
	n, err := w.Write([]byte("data"))
	assertErrWrite(t, n, err)
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
		assertErrWrite(t, n, err)
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
		assertErrWrite(t, n, err)
	})
}

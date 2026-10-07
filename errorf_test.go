package werr

import (
	"errors"
	"strings"
	"testing"
)

func TestErrorf(t *testing.T) {
	t.Run("без оборачивания", func(t *testing.T) {
		err := Errorf("invalid config: %s", "app.yaml")
		got := err.Error()

		if !strings.Contains(got, "invalid config: app.yaml") {
			t.Errorf("Errorf() = %q, want contains 'invalid config: app.yaml'", got)
		}
		// Кадр должен быть: адрес, пробел, сообщение.
		if !strings.Contains(got, ") invalid config: app.yaml") {
			t.Errorf("Errorf() = %q, want address separated by space", got)
		}
	})

	t.Run("с оборачиванием через %w", func(t *testing.T) {
		base := errors.New("connection refused")
		err := Errorf("failed to create client: %w", base)
		got := err.Error()

		if !strings.Contains(got, "failed to create client") {
			t.Errorf("Errorf() = %q, want contains 'failed to create client'", got)
		}
		if !strings.Contains(got, "connection refused") {
			t.Errorf("Errorf() = %q, want contains 'connection refused'", got)
		}
		if !strings.Contains(got, "failed to create client: connection refused") {
			t.Errorf("Errorf() = %q, want message joined by ': '", got)
		}
	})

	t.Run("errors.Is через %w", func(t *testing.T) {
		base := errors.New("base")
		err := Errorf("context: %w", base)

		if !errors.Is(err, base) {
			t.Errorf("errors.Is(err, base) = false, want true")
		}
	})

	t.Run("Cause через %w", func(t *testing.T) {
		base := errors.New("base")
		err := Errorf("context: %w", base)

		if Cause(err) != base {
			t.Errorf("Cause() = %v, want base", Cause(err))
		}
	})

	t.Run("вложенный %w с werr", func(t *testing.T) {
		e := New("root")
		e = Wrap(e, "first")
		e = Errorf("second: %w", e)
		e = Wrap(e, "third")

		got := e.Error()
		iThird := strings.Index(got, "third")
		iSecond := strings.Index(got, "second")
		iFirst := strings.Index(got, "first")
		iRoot := strings.Index(got, "root")

		if !(iThird < iSecond && iSecond < iFirst && iFirst < iRoot) {
			t.Errorf("wrong order: %q", got)
		}
	})
}

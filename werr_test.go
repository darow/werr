package werr

import (
	"errors"
	"strings"
	"testing"
)

func TestWrapNil(t *testing.T) {
	if err := Wrap(nil); err != nil {
		t.Fatalf("Wrap(nil) = %v, want nil", err)
	}
	if err := Wrap(nil, "msg"); err != nil {
		t.Fatalf("Wrap(nil, msg) = %v, want nil", err)
	}
}

func TestWrapForms(t *testing.T) {
	base := errors.New("base")

	// Форма 1: без сообщения.
	e1 := Wrap(base)
	if !strings.Contains(e1.Error(), "base") {
		t.Errorf("Wrap(err) = %q, want contains base", e1.Error())
	}

	// Форма 2: со строкой.
	e2 := Wrap(base, "context")
	if !strings.Contains(e2.Error(), "context") || !strings.Contains(e2.Error(), "base") {
		t.Errorf("Wrap(err, msg) = %q, want contains context and base", e2.Error())
	}

	// Форма 3: с форматом и аргументами.
	e3 := Wrap(base, "id %d", 42)
	if !strings.Contains(e3.Error(), "id 42") {
		t.Errorf("Wrap(err, fmt, args) = %q, want contains id 42", e3.Error())
	}
}

func TestOrder(t *testing.T) {
	// Создаём цепочку и проверяем порядок кадров:
	// самый свежий Wrap должен идти первым.
	e := New("root")
	e = Wrap(e, "first")
	e = Wrap(e, "second")
	e = Wrap(e, "third")

	got := e.Error()
	iThird := strings.Index(got, "third")
	iSecond := strings.Index(got, "second")
	iFirst := strings.Index(got, "first")
	iRoot := strings.Index(got, "root")

	if !(iThird < iSecond && iSecond < iFirst && iFirst < iRoot) {
		t.Errorf("wrong order: %q", got)
	}
}

func TestCause(t *testing.T) {
	base := errors.New("base")
	e := Wrap(Wrap(Wrap(base, "a"), "b"), "c")
	if Cause(e) != base {
		t.Errorf("Cause() = %v, want base", Cause(e))
	}
}

func TestErrorsIs(t *testing.T) {
	base := errors.New("base")
	e := Wrap(base, "context")
	if !errors.Is(e, base) {
		t.Errorf("errors.Is(e, base) = false, want true")
	}
}

func TestFormat(t *testing.T) {
	base := errors.New("base")
	e := Wrap(base, "context")
	got := e.Error()

	// Адрес отделён от сообщения пробелом, уровни — ": ".
	if !strings.Contains(got, ") context") {
		t.Errorf("address not separated by space: %q", got)
	}
	if !strings.Contains(got, "context: base") {
		t.Errorf("levels not separated by ': ': %q", got)
	}
}

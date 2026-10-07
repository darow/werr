// Package werr — минималистичный пакет для оборачивания ошибок
// с сохранением стека вызовов и путей относительно корня проекта.
package werr

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"sync"
)

// Frame описывает одну точку в стеке вызовов.
type Frame struct {
	File     string
	Line     int
	Function string
}

// level — один уровень оборачивания: сообщение и место вызова.
type level struct {
	msg   string
	frame Frame
}

// wrappedError хранит все уровни оборачивания и исходную причину.
type wrappedError struct {
	levels []level // от самого глубокого (создание) к самому свежему (последний Wrap)
	cause  error   // исходная ошибка, если она не наша
}

func (e *wrappedError) Error() string {
	// Собираем цепочку в формате: frame1: msg1: frame2: msg2: ...: cause
	var s string
	for i := len(e.levels) - 1; i >= 0; i-- {
		l := e.levels[i]
		if l.frame.File != "" {
			s += fmt.Sprintf("%s:%d (%s) ", l.frame.File, l.frame.Line, l.frame.Function)
		}
		if l.msg != "" {
			s += l.msg
		}
		if i > 0 || e.cause != nil {
			s += ": "
		}
	}
	if e.cause != nil {
		s += e.cause.Error()
	}
	return s
}

// Cause возвращает исходную (самую вложенную) ошибку.
func (e *wrappedError) Cause() error {
	if e.cause == nil {
		return e
	}
	if c, ok := e.cause.(interface{ Cause() error }); ok {
		return c.Cause()
	}
	return e.cause
}

// Unwrap нужен для совместимости с errors.Is/errors.As из stdlib.
func (e *wrappedError) Unwrap() error { return e.cause }

// --- Поиск корня проекта ---

var (
	projectRoot     string
	projectRootOnce sync.Once
)

func projectRootDir() string {
	projectRootOnce.Do(func() {
		wd, err := os.Getwd()
		if err != nil {
			projectRoot = ""
			return
		}
		dir := wd
		for {
			if _, err := os.Stat(filepath.Join(dir, "go.mod")); err == nil {
				projectRoot = dir
				return
			}
			parent := filepath.Dir(dir)
			if parent == dir {
				projectRoot = wd
				return
			}
			dir = parent
		}
	})
	return projectRoot
}

func relativePath(abs string) string {
	root := projectRootDir()
	if root == "" {
		return abs
	}
	rel, err := filepath.Rel(root, abs)
	if err != nil {
		return abs
	}
	return filepath.ToSlash(rel)
}

func callerFrame(skip int) Frame {
	pc, file, line, ok := runtime.Caller(skip)
	if !ok {
		return Frame{}
	}
	fn := runtime.FuncForPC(pc)
	fnName := ""
	if fn != nil {
		fnName = fn.Name()
	}
	return Frame{
		File:     relativePath(file),
		Line:     line,
		Function: fnName,
	}
}

// --- Публичный API ---

// New создаёт новую ошибку с захватом места вызова.
func New(message string) error {
	return &wrappedError{
		levels: []level{{msg: message, frame: callerFrame(2)}},
	}
}

// Cause возвращает исходную ошибку (совместимо с pkg/errors.Cause).
func Cause(err error) error {
	for err != nil {
		if c, ok := err.(interface{ Cause() error }); ok {
			err = c.Cause()
			continue
		}
		break
	}
	return err
}

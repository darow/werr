package werr

// Errorf создаёт новую ошибку с форматированием и захватом места вызова.
// Совместима с fmt.Errorf, включая глагол %w для оборачивания.
//
// Поддерживает две формы:
//
//	Errorf("static message")                  // без оборачивания
//	Errorf("failed to process id %d: %w", id, err) // с оборачиванием
//
// Если в строке есть %w, соответствующая ошибка становится cause.

func Errorf(format string, args ...interface{}) error {
	// fmt.Errorf корректно обрабатывает %w, включая вложение.
	flat := fmt.Errorf(format, args...)

	// Ищем cause: если flat — это *fmt.wrapError, у него есть Unwrap.
	cause := errors.Unwrap(flat)

	msg := flat.Error()
	if cause != nil {
		// Убираем суффикс ": <cause>" из сообщения, он добавится при печати.
		msg = strings.TrimSuffix(msg, ": "+cause.Error())
	}

	if cause == nil {
		return &wrappedError{
			levels: []level{{msg: msg, frame: callerFrame(2)}},
		}
	}

	if we, ok := cause.(*wrappedError); ok {
		levels := make([]level, len(we.levels), len(we.levels)+1)
		copy(levels, we.levels)
		levels = append(levels, level{msg: msg, frame: callerFrame(2)})
		return &wrappedError{levels: levels, cause: we.cause}
	}

	return &wrappedError{
		levels: []level{{msg: msg, frame: callerFrame(2)}},
		cause:  cause,
	}
}

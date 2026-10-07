package werr

// Wrap оборачивает ошибку, добавляя контекст и место вызова.
//
//	Wrap(err)                                 // только ошибка, без сообщения
//	Wrap(err, "message")                      // ошибка + строка
//	Wrap(err, "failed to process id %d", id)  // ошибка + формат + аргументы
//
// Если err == nil, возвращает nil.
func Wrap(err error, args ...interface{}) error {
	if err == nil {
		return nil
	}

	var msg string
	if len(args) > 0 {
		if format, ok := args[0].(string); ok {
			if len(args) == 1 {
				msg = format
			} else {
				msg = fmt.Sprintf(format, args[1:]...)
			}
		}
	}

	newLevel := level{msg: msg, frame: callerFrame(2)}

	// Если причина — наша ошибка, добавляем уровень поверх её уровней.
	if we, ok := err.(*wrappedError); ok {
		levels := make([]level, len(we.levels), len(we.levels)+1)
		copy(levels, we.levels)
		levels = append(levels, newLevel)
		return &wrappedError{levels: levels, cause: we.cause}
	}

	// Иначе причина — чужая ошибка, кладём её в cause.
	return &wrappedError{
		levels: []level{newLevel},
		cause:  err,
	}
}

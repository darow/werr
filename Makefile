.PHONY: test test-verbose test-race test-cover lint vet fmt clean

PKG ?= ./

# Обычный прогон тестов.
test:
	go test $(PKG)

# Подробный вывод — видно каждый тест.
test-verbose:
	go test -v $(PKG)

# Прогон с детектором гонок. Полезно, если пакет начнут использовать из горутин.
test-race:
	go test -race $(PKG)

# Прогон с покрытием и HTML-отчётом.
test-cover:
	go test -coverprofile=coverage.out $(PKG)
	go tool cover -html=coverage.out -o coverage.html
	@echo "Отчёт: coverage.html"

# Стандартные проверки.
vet:
	go vet $(PKG)

# Форматирование кода.
fmt:
	gofmt -w .

# Всё сразу — то, что имеет смысл гонять перед пушем.
lint: fmt vet test

# Удаление артефактов.
clean:
	rm -f coverage.out coverage.html
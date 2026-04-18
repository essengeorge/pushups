APP_NAME=punkpushups
BINARY_PATH=./$(APP_NAME)
SERVICE_NAME=pushups.service

.PHONY: build update restart status logs

build:
	@echo "--- Сборка проекта... ---"
	go build -o $(BINARY_PATH) ./cmd/main.go
	@echo "Готово: $(BINARY_PATH)"

restart:
	@echo "--- Перезапуск сервиса $(SERVICE_NAME)... ---"
	sudo systemctl restart $(SERVICE_NAME)
	@systemctl status $(SERVICE_NAME) --no-pager

status:
	systemctl status $(SERVICE_NAME)

logs:
	journalctl -u $(SERVICE_NAME) -f

update:
	@echo "--- Обновление кода из Git... ---"
	git pull origin main
	$(MAKE) build
	$(MAKE) restart

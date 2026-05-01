APP_NAME=punkpushups
BINARY_PATH=./$(APP_NAME)
SERVICE_NAME=pushups.service

.PHONY: build update restart status logs

build:
	@echo "--- Building... ---"
	go build -o $(BINARY_PATH) .
	@echo "Готово: $(BINARY_PATH)"

restart:
	@echo "--- Restarting service $(SERVICE_NAME)... ---"
	sudo systemctl restart $(SERVICE_NAME)
	@systemctl status $(SERVICE_NAME) --no-pager

status:
	systemctl status $(SERVICE_NAME)

logs:
	journalctl -u $(SERVICE_NAME) -f

update:
	@echo "--- Pulling code from Git... ---"
	git pull origin main
	$(MAKE) build
	$(MAKE) restart

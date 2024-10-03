# ── Database ────────────────────────────────────────────────────────────────────

.PHONY: db_up
db_up:
	docker-compose up postgres

.PHONY: db_up_d
db_up_d:
	docker-compose up postgres -d

.PHONY: db_down
db_down:
	docker-compose down postgres

# ── API ─────────────────────────────────────────────────────────────────────────

.PHONY: up
run_app:
	docker-compose up

.PHONY: run_app
up:
	docker-compose up --build

.PHONY: clean
clean: 
	docker-compose down --rmi all --volumes --remove-orphans

.PHONY: restart
restart: 
	docker-compose down
	docker-compose up --build

.PHONY: test
test:
	@echo "Running all tests..."
	if go test ./... -v -count=1; then \
		echo "All tests passed successfully!"; \
	else \
		echo "Some tests failed!"; \
		exit 1; \
	fi

.PHONY: swag
swag:
	swag init -g cmd/api/main.go -o internal/swagger/docs
	@echo "Swagger documentation has been regenerated and saved in internal/swagger/docs"
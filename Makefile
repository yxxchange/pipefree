## gen sql
gen:
	@echo "Generating SQL files..."
	go run ./cmd/gorm_gen/main.go


migrate:
	@echo "Running migrations..."
	go run ./cmd/migrate/main.go
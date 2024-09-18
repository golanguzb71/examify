SETUP_SCRIPT_PATH=./scripts/setup_databases.sh
DROP_SCRIPT_PATH=./scripts/drop_databases.sh

.PHONY: all clean

migrate_up:
	@echo "Running database setup script..."
	./$(SETUP_SCRIPT_PATH)
migrate_down:
	@echo "Running database drop script..."
	./$(DROP_SCRIPT_PATH)

## help: print this help message
help:
	@echo "Usage:"
	@sed -n 's/^##//p' ${MAKEFILE_LIST} | column -t -s ":" | sed -e 's/^/  /'

## run-compose: runs the project using docker compose
.PHONY: run-compose
run-compose:
	docker-compose up --build
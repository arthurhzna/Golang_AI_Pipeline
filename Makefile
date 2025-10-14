## Live reload:
watch-prepare: ## Install the tools required for the watch command
	curl -sSfL https://raw.githubusercontent.com/cosmtrek/air/master/install.sh | sh

watch: ## Run the service with hot reload
	bin/air

## Build:
build: ## Build the service
	go build -o task-queue

## Docker:
docker-compose: ## Start the service in docker
	docker-compose up -d --build --force-recreate


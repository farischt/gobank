NAME=gobank
VERSION=0.0.1

## deps: Install the dependencies.
deps:
	@go mod download

## build: Compile the binary file.
build: 
	@go build -o bin/$(NAME)

## run: Run the web server in development mode.
run: build
	@./bin/$(NAME) -e dev

## dev: Run the web server in development mode with a watcher.
## make sure to have https://github.com/githubnemo/CompileDaemon in your device
dev:
	@CompileDaemon -command="./$(NAME) -e dev" -exclude-dir=.git -include=Makefile -include="*.json"

## start: Run the web server in production mode.
start: 
	@./bin/$(NAME) -e prod

## clean: Remove previous build.
clean:
	@rm -rf bin

## migrate-create: Create a new migration.
migrate-create:
	@migrate create -ext sql -dir database/migrations -seq $(name)

## migrate-up: Run the migrations.
migrate-up: 
	@migrate -path database/migrations -database "postgres://gobank:root@localhost:5432/gobank-db?sslmode=disable" -verbose up

## migrate-down: Rollback the migrations.
migrate-down:
	@migrate -path database/migrations -database "postgres://gobank:root@localhost:5432/gobank-db?sslmode=disable" -verbose down

## migrate-goto: Go to a specific migration version. Use the -v flag to specify the version. (e.g. make migrate-goto v=1)
migrate-goto:
	@migrate -path database/migrations -database "postgres://gobank:root@localhost:5432/gobank-db?sslmode=disable" -verbose goto $(v)

## migrate-fix: Force the migrations. Use the -v flag to specify the version.
migrate-fix: 
	@migrate -path database/migrations -database "postgres://gobank:root@localhost:5432/gobank-db?sslmode=disable" -verbose force $(v)

## test: Run the tests and generate the coverage report.
test: 
	@go test -v ./... -cover
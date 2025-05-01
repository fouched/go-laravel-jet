BINARY_NAME=celeritasApp.exe
SHELL=cmd.exe

## build: builds all binaries
build:
# 	can use vendor below to store and update local copy of dependencies
#	@go mod vendor
	@go build -o tmp/${BINARY_NAME} .
	@echo Celeritas built!

run: build
	@echo Staring Celeritas...
	@.\tmp\${BINARY_NAME} &
	@echo Celeritas started!

clean:
	@echo Cleaning...
	@go clean
	@del .\tmp\${BINARY_NAME}
	@echo Cleaned!

start-compose:
	docker-compose up -d

stop-compose:
	docker-compose down

test:
	@echo Testing...
	@go test ./...
	@echo Done!

start: run

stop:
	@echo "Starting the front end..."
	@taskkill /IM ${BINARY_NAME} /F
	@echo Stopped Celeritas

restart: stop start
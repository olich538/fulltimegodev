build:
	@go build -o bin/api 

run: build
	@./bin/api

seed:
	@go run scripts/seed.go

test:
	@go test -v ./...

docker:
	echo "building an image"
	@docker build -t api .
	echo "Running API inside Docker container"
	@docker run -p 3000:3000 api
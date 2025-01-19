dev:
	air
build:
	go build -o ./build/main cmd/app/main.go
serve:
	sudo docker-compose up
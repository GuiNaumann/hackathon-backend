build:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o server main.go

docker:
	docker build -t hackathon:latest .

docker-compose:
	docker-compose down && sleep 5 && docker-compose up -d

all: build docker docker-compose

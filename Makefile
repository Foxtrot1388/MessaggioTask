.SILENT:

run:
	docker-compose up -d

initapi:
	swag init -g ./cmd/app/main.go -o ./api
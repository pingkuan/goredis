run:
	docker run --name redis -d -p 6379:6379 redis:7.0
	go run main.go
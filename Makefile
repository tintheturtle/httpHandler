build:
	go build -o main main.go

run:
	go run main.go -url=http://sweet-worker.tintheturtle.workers.dev:80/links -profile=7

clean:
	go clean
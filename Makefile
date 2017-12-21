build:
	GOOS=linux go build -v
	docker build -t dl .
	
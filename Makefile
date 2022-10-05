docker-build:
	docker build --no-cache -f build/Dockerfile -t todo .

run:
	docker run -d -p8081:8081 todo

test:
	go test ./tests/
# Build the Docker image
build:
	docker build -t redis-implementation:latest .

# Run the container
run:
	docker run -it --rm \
		--name redis-implementation \
		-p 9999:9999 \
		redis-implementation:latest

build-and-run:build run
# Build the Docker image
build:
	docker build -t redis-implementation:latest .

# Run the container
run:
	docker run -d \
		--name redis-implementation \
		-p 9999:9999 \
		redis-implementation:latest
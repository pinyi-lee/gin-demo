docker:
	@rm -f build/Dockerfile
	@cp build/Dockerfile.temp build/Dockerfile
	@docker build -t gin_service -f build/Dockerfile .

start:
	@docker-compose -f build/docker-compose.yml up -d

stop:
	@docker-compose -f build/docker-compose.yml stop

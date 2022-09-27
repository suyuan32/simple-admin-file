version := 0.0.8
docker:
	sudo docker build -f Dockerfile -t fileapi:$(version) .

run-docker:
	sudo docker run -d --name fileapi-$(version) --network docker-compose_simple-admin --network-alias fileapi -p 9101:9101 fileapi:$(version)
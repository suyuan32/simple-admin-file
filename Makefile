RED  =  "\e[31;1m"
GREEN = "\e[32;1m"
YELLOW = "\e[33;1m"
BLUE  = "\e[34;1m"
PURPLE = "\e[35;1m"
CYAN  = "\e[36;1m"

docker:
	sudo docker build -f Dockerfile -t ${DOCKER_USERNAME}/fileapi:${VERSION} .
	@printf $(GREEN)"[SUCCESS] build docker successfully"

publish-docker:
	echo "${DOCKER_PASSWORD}" | docker login --username ${DOCKER_USERNAME} --password-stdin http://${REPO}
	docker push ${REPO}/${DOCKER_USERNAME}/fileapi:${VERSION}
	@printf $(GREEN)"[SUCCESS] publish docker successfully"

gen-api:
	goctls api go --api ./api/desc/file.api --dir ./api --transErr=true
	swagger generate spec --output=./file.yml --scan-models
	@printf $(GREEN)"[SUCCESS] generate API successfully"

gen-ent:
	go run -mod=mod entgo.io/ent/cmd/ent generate --template glob="./pkg/ent/template/*.tmpl" ./pkg/ent/schema
	@printf $(GREEN)"[SUCCESS] generate ent successfully"

serve-swagger:
	lsof -i:36666 | awk 'NR!=1 {print $2}' | xargs killall -9 || true
	@printf $(GREEN)"[SUCCESS] serve swagger-ui successfully"
	swagger serve -F=swagger --port 36667 file.yml
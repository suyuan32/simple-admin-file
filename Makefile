docker:
	sudo docker build -f Dockerfile -t ${DOCKER_USERNAME}/fileapi:${VERSION} .

publish-docker:
	echo "${DOCKER_PASSWORD}" | docker login --username ${DOCKER_USERNAME} --password-stdin http://${REPO}
	docker push ${REPO}/${DOCKER_USERNAME}/fileapi:${VERSION}
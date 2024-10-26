.PHONY: start
start: 
	docker run -d -p 9000:80 --name iamfoo traefik/whoami

.PHONY: stop
stop:
	docker stop iamfoo && docker rm iamfoo
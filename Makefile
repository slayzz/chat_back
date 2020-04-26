APP=chat-service
NGINX=nginx

start:
	docker-compose up -d

stop:
	docker-compose down


rebuild:
	docker-compose up -d --no-deps --build --force-recreate $(APP)
	docker-compose logs -f $(APP)

rnginx:
	docker-compose up -d --no-deps --build --force-recreate $(NGINX)

logs:
	docker-compose logs -f $(APP)

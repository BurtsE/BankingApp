run:
	docker compose -f .\deploy\compose.yml -p "banking_app" up --build  app
down:
	docker compose -f .\deploy\compose.yml -p "banking_app" down
bootstrap:
	bash deploy/bootstrap.sh
up:
	docker compose -f deploy/docker-compose.yml up -d
down:
	docker compose -f deploy/docker-compose.yml down
migrate:
	bash scripts/migrate_all.sh
seed:
	bash scripts/seed_data.sh
test:
	bash scripts/integration_test.sh
logs:
	docker compose -f deploy/docker-compose.yml logs -f --tail=100
health:
	bash scripts/health_check.sh
reset:
	bash scripts/dev_reset.sh
clean:
	docker compose -f deploy/docker-compose.yml down -v --remove-orphans

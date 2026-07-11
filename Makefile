ENV ?= .env
include $(ENV)
export

# Define db url
DB_URL=mysql://$(DB_USER):$(DB_PASS)@tcp($(DB_HOST):$(DB_PORT))/$(DB_NAME)
ifeq ($(DB_TLS),true)
DB_URL:=$(DB_URL)?tls=true
endif

# Desc: To create a new migration
# Usage: make migration name=create_users
migration:
	migrate create -ext sql -dir ./scripts/migrations -seq $(name)

# Desc: To create insert seed data migration
# Usage: make seed name=insert_initial_data
seed:
	migrate create -ext sql -dir ./scripts/seeds -seq $(name)

# Usage: make fakes [ARGS="--students-csv=..."]
fakes:
	go run ./scripts/faker $(ARGS)

backfill:
	go run ./scripts/faker --backfill-iir

# Desc: Fetch PSGC data from API and save to JSON (run once, commit the file)
# Usage: make psgc-fetch
psgc-fetch:
	go run cmd/psgc/psgc.go

# Desc: Seed locations from psgc_data.json into the database
# Usage: make locations
locations:
	go run cmd/locations/locations.go

# Desc: To refresh database with cli
# Usage: make migrate-up
migrate-up:
	migrate -path scripts/migrations -database "$(DB_URL)" up

# Desc: To undo database migrations
# Usage: make migrate-down
migrate-down:
	migrate -path scripts/migrations -database "$(DB_URL)" down

# Desc: To undo all database migrations
# Usage: make migrate-reset
migrate-reset:
	migrate -path scripts/migrations -database "$(DB_URL)" drop -f

# Desc: To apply seed data to database
# Usage: make seed-up
seed-up:
	migrate -path scripts/seeds -database \
	"$(DB_URL)$(if $(filter true,$(DB_TLS)),&,?)x-migrations-table=seed_migrations" up

refresh: migrate-reset
	migrate -path scripts/migrations -database "$(DB_URL)" up 30
	$(MAKE) seed-up
	migrate -path scripts/migrations -database "$(DB_URL)" up
	$(MAKE) locations

# Desc: To generate swagger docs
# Usage: make swagger-internal
swagger-internal:
	swag init -g main.go \
	--parseDependency --parseInternal \
	--dir ./cmd/api,\
	./internal/features/auth,\
	./internal/features/users,\
	./internal/features/appointments,\
	./internal/features/slips,\
	./internal/features/students \
	--output ./docs/internal_docs \
    --instanceName internal

swagger-integrations:
	swag init -g main.go \
	--parseDependency --parseInternal \
	--dir ./cmd/api,\
	./internal/features/students/integrations,\
	./internal/features/m2mclients \
	--output ./docs/integrations \
	--instanceName integrations

compose-up:
	docker compose --env-file $(ENV) up --build

QA_HEAVY ?= false

compose-staging:
	docker compose -f docker-compose.staging.yml \
		--env-file $(ENV) up --build -d

compose-prod:
	docker compose -f docker-compose.prod.yml \
		--env-file $(ENV) up --build -d

# Desc: Build and push production API image
push-api-prod:
	docker build --network=host --target prod \
		-t ghcr.io/olazo-johnalbert/capstone-api:latest .
	docker push ghcr.io/olazo-johnalbert/capstone-api:latest

# Desc: Build and push staging API image
push-api-staging:
	docker build --network=host --target prod \
		-t ghcr.io/olazo-johnalbert/capstone-staging-api:latest .
	docker push ghcr.io/olazo-johnalbert/capstone-staging-api:latest

# Desc: Build and push production AI image
push-ai-prod:
	docker build --network=host \
		-t ghcr.io/olazo-johnalbert/capstone-ai:latest \
		-f ../guisis-ai/Dockerfile ../guisis-ai
	docker push ghcr.io/olazo-johnalbert/capstone-ai:latest

# Desc: Build and push staging AI image
push-ai-staging:
	docker build --network=host \
		-t ghcr.io/olazo-johnalbert/capstone-staging-ai:latest \
		-f ../guisis-ai/Dockerfile ../guisis-ai
	docker push ghcr.io/olazo-johnalbert/capstone-staging-ai:latest

# Desc: Push both production images
push-all-prod: push-api-prod push-ai-prod

# Desc: Push both staging images
push-all-staging: push-api-staging
ifeq ($(QA_HEAVY),true)
push-all-staging: push-ai-staging
endif

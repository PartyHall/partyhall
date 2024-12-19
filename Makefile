include env
export

VERSION = 0.8
COMMIT = $(shell git rev-parse --short HEAD)

init:
	@echo "Ensuring frontend folders exists"
	@mkdir -p backend/frontend/appliance backend/frontend/app
	@echo "Removing containers..."
	@docker compose down --remove-orphans
	@echo "Dropping the database..."
	@sudo rm -rf ./0_DATA/data/database.sqlite
	@sudo rm -rf ./0_DATA/data/events
	@echo "Starting the containers..."
	@docker compose up -d
	# @$(MAKE) fixtures

run-app:
	@cd backend && go run .

hwhandler:
	@cd backend && go run . hwhandler

compile-sdk:
	@sudo rm -rf sdk/dist
	@docker run --rm -v $(PWD)/sdk:/sdk -w /sdk node:lts npm install
	@docker run --rm -v $(PWD)/sdk:/sdk -w /sdk node:lts npx tsc

build-release:
	@docker buildx build -f docker/prod/Dockerfile --build-arg PARTYHALL_VERSION=$(VERSION) --build-arg PARTYHALL_COMMIT=$(COMMIT) -t partyhall:latest --load .
	@docker run --rm -v $(PWD)/build:/binaries partyhall:latest /bin/sh -c 'cp /partyhall-*-linux-* /binaries/'

fixtures:
	@$(MAKE) create-users
	@$(MAKE) create-events


create-users:
	@echo "Creating the default users"
	@cd backend && go run . user create-admin --username admin --password password --name Administrator
	@cd backend && go run . user create user password "Some user"

create-events:
	@echo "Creating an event"
	@curl -s -o /dev/null -L -X POST 'http://localhost:8080/api/webapp/events' -H "Authorization: Bearer $(VITE_PARTYHALL_APPLIANCE_JWT)" -H 'Content-Type: application/json' --data-raw '{"name":"New event","author":"Some author","date":"2024-01-10T11:58:00Z","location":"Some place"}'
	@echo "Creating a second event"
	@curl -s -o /dev/null -L -X POST 'http://localhost:8080/api/webapp/events' -H "Authorization: Bearer $(VITE_PARTYHALL_APPLIANCE_JWT)" -H 'Content-Type: application/json' --data-raw '{"name":"Second event","author":"Another author","date":"2024-01-11T21:12:00Z","location":"At the beach"}'

	@echo "Setting mode"
	@curl -s -o /dev/null -L -X POST 'http://localhost:8080/api/webapp/settings/mode/photobooth' -H "Authorization: Bearer $(VITE_PARTYHALL_APPLIANCE_JWT)"

gen-jwt:
	@docker compose exec app go run . dev jwt

show-debug:
	@echo "Displaying debug through MQTT"
	@docker compose run --rm mosquitto mosquitto_pub -h mosquitto -t partyhall/display_debug -m ""

take-picture:
	@echo "Taking picture through MQTT"
	@docker compose run --rm mosquitto mosquitto_pub -h mosquitto -t partyhall/take_picture -m ""

shutdown:
	@echo "Shutdown through MQTT"
	@docker compose run --rm mosquitto mosquitto_pub -h mosquitto -t partyhall/shutdown -m ""

################
# Experimental #
################

lint-ts:
	@cd sdk && npx prettier . --write
	@docker compose exec admin npx prettier . --write
	@docker compose exec frontend npx prettier . --write
	@docker compose exec admin npm run lint -- --fix
	@docker compose exec frontend npm run lint -- --fix

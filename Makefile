.PHONY: mocks
mocks:
	sleep 1 && rm -rfd mocks && mockery

swag:
	swag init --parseDependency --parseInternal --parseDepth 1 -g ./cmd/main.go

run:  swag
	go run cmd/main.go

build:  swag
	go build -o tmp/main cmd/main.go

build-win:  swag
	go build -o tmp/main.exe cmd/main.go
		
docker-run:
	docker run -d -p 8080:8080 \
	--name project-sprint-shopifyx-api \
	project-sprint-shopifyx-api:latest

air:
	air -c .air.toml

air-win:
	air -c .air.win.toml

# make startProm
.PHONY: startProm
startProm:
	docker run \
	-p 9090:9090 \
	--name=prometheus \
	-v $(shell pwd)/prometheus.yml:/etc/prometheus/prometheus.yml \
	-d \
	prom/prometheus

# make startGrafana
# for first timers, the username & password is both `admin`
.PHONY: startGrafana
startGrafana:
	docker volume create grafana-storage
	docker volume inspect grafana-storage
	docker run -d -p 3000:3000 --name=grafana grafana/grafana-oss || docker start grafana
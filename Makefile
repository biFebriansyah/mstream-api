APP_NAME=mstream-api
BUILD_DIR="./build/$(APP_NAME)"
DB_DRIVER=postgres
DB_SOURCE="postgresql://postgres:postgres@localhost/project?sslmode=disable&search_path=stream"
MIGRATIONS_DIR=./migrations

install:
	go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest && go get -u ./... && go mod tidy

run:
	nodemon --exec "go run" .

build:
	mkdir -p ./build && CGO_ENABLED=0 GOOS=linux go build -o ${BUILD_DIR}

test:
	go test -cover -v ./...

seed:
	go run ./seeder/*.go

migrate-init:
# make migrate-init name=mstream-genres
	migrate create -dir ${MIGRATIONS_DIR} -ext sql $(name)

migrate-up:
	migrate -path ${MIGRATIONS_DIR} -database ${DB_SOURCE} -verbose up

migrate-down:
	migrate -path ${MIGRATIONS_DIR} -database ${DB_SOURCE} -verbose down

migrate-reset:
	migrate -path ${MIGRATIONS_DIR} -database ${DB_SOURCE} drop

migrate-fix:
	migrate -path ${MIGRATIONS_DIR} -database ${DB_SOURCE} force 0

hello:
	echo "$DB_PASS"
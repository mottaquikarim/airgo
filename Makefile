# Definitions
ROOT                    := $(PWD)#$(shell dirname $(realpath $(firstword $(MAKEFILE_LIST))))
GO_HTML_COV             := ./coverage.html
GO_TEST_OUTFILE         := ./c.out

#   Format according to gofmt: https://github.com/cytopia/docker-gofmt
#   Usage:
#       make fmt
#       make fmt path=src/elastic/index_setup.go
fmt:
ifdef path
	docker run --rm -v ${ROOT}:/data cytopia/gofmt -w ${path}
else
	docker run --rm -v ${ROOT}:/data cytopia/gofmt -w .
endif

#   Usage:
#       make lint
lint:
	docker run --rm -v ${ROOT}:/app -w /app golangci/golangci-lint:v1.27.0 golangci-lint run -v

#   Start services 
compose: 
	docker-compose rm -f
	docker-compose build --no-cache 
	docker-compose up -d

# 	Spin down project
#	Usage:
#		make stop
stop:
	docker-compose down

#   Download data from airtable:
#	Usage:
#		make refresh arguments="-api-key=API_KEY -base-id=BASE_ID"
refresh: 
	docker-compose exec -T airgo_develop go build -o build/refresh_data cmd/refresh_data/main.go || true
	docker-compose exec -T airgo_develop ./build/refresh_data ${arguments}

#   Usage:
#       make test-dev
#       make test-dev package=util
test-dev:
ifdef package 
	docker-compose exec -T airgo_develop go test -v ./${NAME}/${package}/... -coverprofile=${GO_TEST_OUTFILE}
else
	docker-compose exec -T airgo_develop go test ./${NAME}/... -coverprofile=${GO_TEST_OUTFILE}
endif
	docker-compose exec -T airgo_develop go tool cover -html=${GO_TEST_OUTFILE} -o ${GO_HTML_COV}

#   Usage:
#       make logs
logs:
	docker-compose logs -f airgo_develop

develop: compose logs


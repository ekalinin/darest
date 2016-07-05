.PHONY: env build

NAME=darest
EXEC=${NAME}
GOVER=1.6.2
ENVNAME=${NAME}${GOVER}
GHBASE=github.com/ekalinin
GHNAME=${GHBASE}/${NAME}

#
# For virtual environment create with
# https://github.com/ekalinin/envirius
#
env-create:
	@bash -c ". ~/.envirius/nv && nv mk ${ENVNAME} --go-prebuilt=${GOVER}"

env-fix:
	@bash -c ". ~/.envirius/nv && nv do ${ENVNAME} 'make fix-paths'"

env-deps:
	@bash -c ". ~/.envirius/nv && nv do ${ENVNAME} 'make deps'"

env-tools:
	@bash -c ". ~/.envirius/nv && nv do ${ENVNAME} 'make tools'"


env-build:
	@bash -c ". ~/.envirius/nv && nv do ${ENVNAME} 'make build'"

env-init: env-create env-fix env-deps env-tools

env:
	@bash -c ". ~/.envirius/nv && nv use ${ENVNAME}"

#
# Other targets
#

deps:
	#@go get gopkg.in/alecthomas/kingpin.v2
	go get -u github.com/lib/pq
	go get -u github.com/labstack/echo/...

fix-paths:
	@if [ -d "${GOPATH}/src/${GHNAME}" ]; then \
		echo "Already fixed. No actions need."; \
	else \
		mkdir -p ${GOPATH}/src/${GHBASE}; \
		ln -s `pwd` ${GOPATH}/src/${GHBASE}; \
	fi

tools:
	@go get -v github.com/nsf/gocode
	@go get -v github.com/rogpeppe/godef
	@go get -v github.com/golang/lint/golint
	#@go get -u -v github.com/lukehoban/go-find-references
	@go get -v github.com/lukehoban/go-outline
	@go get -v sourcegraph.com/sqs/goreturns
	@go get -v golang.org/x/tools/cmd/gorename
	@go get -v github.com/tpng/gopkgs
	@go get -v github.com/newhook/go-symbols

build:
	@go build -a -tags netgo \
		--ldflags '-s -extldflags "-lm -lstdc++ -static"' \
		-o ${EXEC} ./main.go

#
# Utils
#
start-pg:
	# psql -h 172.17.0.1 -d postgres -U postgres -W
	# sudo docker run -p 5432:5432 --name postgres-db -d postgres:latest
	docker start postgres-db

start:
	# curl -s http://localhost:7788/festival | python -mjson.tool
	@./${NAME} -db-dbname postgres -db-host 172.17.0.1 -port 7788 \
		-db-pass postgres -db-port 5432 -db-user postgres

import-db:
	psql -h 172.17.0.1 -d postgres -U postgres -W \
		-c "\copy festival from './examples/films/data/festival.csv' DELIMITER ',' CSV HEADER;"
	psql -h 172.17.0.1 -d postgres -U postgres -W \ 
		-c "\copy competition from './examples/films/data/competition.csv' DELIMITER ',' CSV HEADER;"
	psql -h 172.17.0.1 -d postgres -U postgres -W \
		-c "\copy director from './examples/films/data/director.csv' DELIMITER ',' CSV HEADER;"
	psql -h 172.17.0.1 -d postgres -U postgres -W \
		-c "\copy film from './examples/films/data/film.csv' DELIMITER ',' CSV HEADER;"
	psql -h 172.17.0.1 -d postgres -U postgres -W \
		-c "\copy film_nomination(competition,film,won) from './examples/films/data/film_nomination.csv' DELIMITER ',' CSV HEADER;"

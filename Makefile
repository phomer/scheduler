.PHONY: all start stop clean fulltest build install test

all: build
	@echo "Built"

start: build
	(cd ./tests; ./scheduled > scheduled.log 2>&1 &)
	(cd ./tests; tail -f *.log)

stop:
	-pkill scheduled

clean: stop
	rm -rf ./tests/data
	rm -f ./tests/*.log
	rm -f ./tests/*.key

fulltest: build
	(cd ./tests; ./simple_test.sh)
	tail ./tests/*.log

build:
	go build -o tests/register cmd/register.go
	go build -o tests/scheduled cmd/scheduled.go
	go build -o tests/schedule cmd/schedule.go 

install:
	@echo "Install into root"

test:
	go test -v github.com/phomer/scheduler/datastore

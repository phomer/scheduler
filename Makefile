test: build
	(cd ./tests; ./simple_test.sh)
	tail ./tests/*.log

all: build
	@echo "Built"

build:
	go build -o tests/register cmd/register.go
	go build -o tests/scheduled cmd/scheduled.go
	go build -o tests/schedule cmd/schedule.go 

install:
	@echo "Install into root"


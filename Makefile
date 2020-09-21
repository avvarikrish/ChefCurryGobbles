all: generate build

generate:
	@for PROTODIR in proto/*; \
	do \
		./generate.sh $$(basename $$PROTODIR); \
	done

build:
	@go build -o bin/ccgobbles ccgobbles_client/ccgobbles/ccgobbles_client.go 
	@go build -o bin/metrics ccgobbles_client/metrics/metrics_client.go

submit:
	@./bc.sh

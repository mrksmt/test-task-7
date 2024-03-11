.PHONY: run-client run-server run-all-in-one docker-build-client docker-build-service prof-async

# lint ...

# test ...

run-client:
	go run cmd/client/main.go

run-servre:
	go run cmd/server/main.go

# test run client and server
run-all-in-one:
	go run cmd/all_in_one/main.go

docker-build-client:
	docker build \
	--no-cache \
	--tag test-task-7/client:latest . \
	--file dockerfile.client	

docker-build-server:
	docker build \
	--no-cache \
	--tag test-task-7/service:latest . \
	--file dockerfile.server	

prof-async:
	cd pkg/pow && \
	go test -bench=. -benchmem -memprofile memprofile.out -cpuprofile profile.out -run Benchmark_GetProofAsync && \
	go tool pprof profile.out       

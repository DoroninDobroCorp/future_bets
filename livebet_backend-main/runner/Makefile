include .env
export $(shell sed 's/=.*//' .env)

NAME=runner

build_r: ### run runner in local environment 
	@echo "Building ${NAME} binary..."
	mkdir -p ./logs
	cd ../${NAME} && env GOOS=darwin CGO_ENABLED=0 go build -o ${NAME} ./cmd
	@echo "Done!"
	./${NAME} >> ./logs/RunnerLogs 2>&1
.PHONY: ./${NAME} 

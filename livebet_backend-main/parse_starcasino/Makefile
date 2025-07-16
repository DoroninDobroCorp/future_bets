include .env
export $(shell sed 's/=.*//' .env)

NAME=parse_starcasino

build_pl:  ### run parse_starcasino in local environment
	@echo "Building ${NAME} binary..."
	cd ../${NAME} && env GOOS=darwin CGO_ENABLED=0 go build -o ${NAME} ./cmd
	@echo "Done!"
	./${NAME}
.PHONY: ./${NAME}

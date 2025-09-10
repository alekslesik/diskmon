# ==================================================================================== #
# HELPERS
# ==================================================================================== #
## help: print this help message
.PHONY: help
help:
	@echo 'Usage:'
	@sed -n 's/^##//p' ${MAKEFILE_LIST} | column -t -s ':' |  sed -e 's/^/ /'
	
## run: go  run
.PHONY: run
run:
	@cp configs/config.yaml build
	@go build -o build/diskmon cmd/diskmon/main.go
	@cd build && CONF_PATH=config.yaml ./diskmon
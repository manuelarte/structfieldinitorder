default: fmt

fmt: format-code
format-code:  ## Format go code and run the fixer, alias: fmt
	golangci-lint custom -v
	./custom-gcl run --fix ./...
.PHONY: fmt format-code


.PHONY: lambda package

build_dir := ./build
lambda_output := $(build_dir)/bootstrap


lambda:
	@echo "Building lambda binary..."
	CGO_ENABLED=0 GOOS=linux GOARCH=arm64 go build -ldflags "-s -w" -o $(lambda_output) ./main.go

package: lambda
	@echo "Packaging lambda binary..."
	zip -j $(build_dir)/function.zip $(lambda_output)

clean:
	@rm -rf $(build_dir)
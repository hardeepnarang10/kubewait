.PHONY: build
build: clean
	@CGO_ENABLED=0 go build -C ./cmd/job -ldflags "-s -w" -v -buildvcs=false -trimpath -buildmode=exe \
	-o ./../../release/service-darwin-arm64 > ./release/service-darwin-arm64.build.log 2>&1

.PHONY: run
run:
	./release/service-darwin-arm64

.PHONY: test
test:
	go test -v -cpu 2 ./...

.PHONY: clean
clean:
	@rm -f ./release/service-darwin-arm64

.PHONY: buildx
buildx:
	@docker buildx \
	build \
	--push \
	--platform linux/amd64,linux/arm64 \
	-f build/job/Dockerfile \
	-t hardeepnarang10/kubewait:0.2-alpha \
	.

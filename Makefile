VERSION := "0.0.2"

build:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -mod=readonly -ldflags="-w -s" -v -o bin/exporter cmd/exporter/main.go
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -mod=readonly -ldflags="-w -s" -v -o bin/notifier cmd/notifier/main.go

tag:
	git tag -a "v$(VERSION)" -m "Release $(VERSION)"
	git push --tags

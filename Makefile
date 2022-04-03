
local: fmt vet
	GOOS=linux GOARCH=amd64 go build -o=./bin/sfc-scheduler pkg/main/main.go


build: local
	docker build --no-cache . -t 326891007/sfc-scheduler:v1.19.9

push: build
	docker push 326891007/sfc-scheduler:v1.19.9

# Run go fmt against code
fmt:
	go fmt ./...

 # Run go vet against code
vet:
	go vet ./...
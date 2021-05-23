init:
	npx create-react-app ui
	cd ui && PUBLIC_URL=http://127.0.0.1:8080 npm run build
build:
	go build -o bin/go-create-react-app cmd/main.go

run: build
	bin/go-create-react-app
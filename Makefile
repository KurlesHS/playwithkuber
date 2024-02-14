build:
	cd hellosayer && CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o ./bin/hellosayer ./cmd/app/...
	cd telebot && CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o ./bin/telebot ./cmd/app/...

dockerize:
	cd hellosayer && docker build . -t hellosayer:v1.0
	cd telebot && docker build . -t telebot:v1.0

build-and-dockerize:

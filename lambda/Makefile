setup:
	export GOPROXY=https://gocenter.io
	export GO111MODULE=on
	gvm linkthis github.com/netroy/serverless-dns

clean:
	-@rm -f app function.zip

build: clean
	go build -o app
	zip app.zip app

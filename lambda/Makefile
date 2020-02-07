setup:
	export GO111MODULE=on
	gvm linkthis github.com/netroy/serverless-dns/lambda

clean:
	-@rm -f app function.zip

build: clean
	go build -o app
	zip app.zip app

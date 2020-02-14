setup:
	export GO111MODULE=on

clean:
	-@rm -f app function.zip

build: clean
	go build -o app
	zip app.zip app

exports:
	export GOPROXY=https://gocenter.io
	export GO111MODULE=on

clean:
	-@rm -f app function.zip

build: clean lint
	go build -o app

deploy: build
	@zip function.zip app
	aws lambda update-function-code --function-name DOHFunction --zip-file fileb://function.zip
	@make clean

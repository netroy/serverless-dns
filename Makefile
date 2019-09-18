install:
	go get -u golang.org/x/lint/golint
	go get -u github.com/golang/dep/cmd/dep
	dep ensure

clean:
	-@rm -f app function.zip

lint:
	golint

build: clean lint
	GOOS=linux go build -o app main.go

deploy: build
	@zip function.zip app
	aws lambda update-function-code --function-name DOHFunction --zip-file fileb://function.zip
	@make clean

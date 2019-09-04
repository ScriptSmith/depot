run: bin/depot
	@PATH="$(PWD)/bin:$(PATH)" heroku local

bin/depot: main.go
	go build -o bin/depot .

clean:
	rm -rf bin
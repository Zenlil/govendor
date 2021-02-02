build:
	go build

init2:
	-rm example/vendor2.json
	./govendor -f example/vendor2.json init

list1:
	./govendor -f example/vendor1.json list

copy1:
	./govendor -f example/vendor1.json copy

add1:
	./govendor -f example/vendor1.json add github.com/marstr/guid foo

delete1:
	./govendor -f example/vendor1.json delete GoShared/types foo --dryrun
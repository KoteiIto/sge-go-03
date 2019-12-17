bench:
	go test ./... -bench . -benchmem -test.run=none -benchtime 2s

.PHONY: test
test:
	go test ./... -v -race -count 1
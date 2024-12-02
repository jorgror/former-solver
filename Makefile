.PHONY: make-play make-solve play solve test

make-play:
	go build -o bin/play cli/play/main.go

make-solve:
	go build -o bin/solve cli/solve/main.go

play: make-play
	./bin/play 20241120

solve: make-solve
	./bin/solve 20241201 infinite 20 3

test:
	go test ./...
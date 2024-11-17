.PHONY: make-play make-solve play solve test

make-play:
	go build -o bin/play cli/play/main.go

make-solve:
	go build -o bin/solve cli/solve/main.go

play: make-play
	./bin/play 20241117

solve: make-solve
	./bin/solve 20241117 beam 10 5 3 50

test:
	go test ./...
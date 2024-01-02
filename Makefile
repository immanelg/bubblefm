all: build


run:
	go run .


build:
	go build -o bubblefm .


clean:
	rm -f ./bubblefm


# TODO
install: build
	mkdir -p ~/.local/bin
	cp -f ./bubblefm ~/.local/bin/
	chmod 755 ~/.local/bin/bubblefm

uninstall:
	rm -f ~/.local/bin/bubblefm


.PHONY: all build clean install uninstall run

ifdef ComSpec
    RM=del /F /Q
else
    RM=rm -f
endif

test:
	go test -v ./...

build:
	go build -C cmd -o ../bin/acb

clean:
	RM bin


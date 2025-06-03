SRC = main.go
DST := ../../bin
NAME = fill_template

.DEFAULT_GOAL: all

.PHONY: all
all: $(DST)/$(NAME).lnx64.exe $(DST)/$(NAME).mac64.exe $(DST)/$(NAME).win64.exe $(DST)/$(NAME).armm.exe

$(DST)/$(NAME).mac64.exe: $(SRC)
	GO111MODULE=on GOOS=darwin GOARCH=amd64 go build -o $@

$(DST)/$(NAME).armm.exe: $(SRC)
	GO111MODULE=on GOOS=darwin GOARCH=arm64 go build -o $@

$(DST)/$(NAME).lnx64.exe: $(SRC)
	GO111MODULE=on GOOS=linux GOARCH=amd64 go build -o $@

$(DST)/$(NAME).win64.exe: $(SRC)
	GO111MODULE=on GOOS=windows GOARCH=amd64 go build -o $@

.PHONY: clean
clean:
	rm -f $(DST)/$(NAME)*

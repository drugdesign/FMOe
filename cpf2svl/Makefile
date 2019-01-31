SRC = main.go write_cpf.go cpf/cpf.go svlwriter/svlwriter.go
DST := ../bin

.DEFAULT_GOAL: all

.PHONY: all
all: $(DST)/cpf2svl.lnx64.exe $(DST)/cpf2svl.mac64.exe $(DST)/cpf2svl.win64.exe

$(DST)/cpf2svl.mac64.exe: $(SRC)
	GO111MODULE=on GOOS=darwin GOARCH=amd64 go build -o $@

$(DST)/cpf2svl.lnx64.exe: $(SRC)
	GO111MODULE=on GOOS=linux GOARCH=amd64 go build -o $@

$(DST)/cpf2svl.win64.exe: $(SRC)
	GO111MODULE=on GOOS=windows GOARCH=amd64 go build -o $@

.PHONY: clean
clean:
	rm -f $(DST)/*

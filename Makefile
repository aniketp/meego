GOCMD=go
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GENERATOR=grammer.bnf
GENERATE=../../../../../bin/gocc


all: test run

test: run
	$(GOTEST) -v 

clean:
	rm -rf src/util src/token src/lexer src/parser src/errors
run:
	cd src; \
	$(GENERATE) $(GENERATOR) # create lexer and parser
	cd ..;

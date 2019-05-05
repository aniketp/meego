GOCMD=go
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GENERATOR=grammer.bnf
GENERATE=../../../../bin/gocc


all: test run

test:
	$(GENERATE) $(GENERATOR)
	$(GOTEST) -v 

clean:
	$(GOCLEAN)
	rm -rf util token lexer parser errors
run:
	$(GENERATE) $(GENERATOR) # create lexer and parser


TARGET = ./cw

# macro for test
ifdef DBG
omit =
else
omit = >/dev/null 2>&1
endif
it_produces  = 2>&1 | grep -q
it_exit_with = test $$? =

###### targets #######

default: build
	$(TARGET)

build: $(TARGET)

$(TARGET): *.go
	go build

test: test-cli test-unit

test-unit:
	go test

test-cli: build
	# smoke test
	$(TARGET) $(omit)
	# messages
	$(TARGET) $(it_produces) "Simple command line tool"
	$(TARGET) -h $(it_produces) "Simple command line tool"
	$(TARGET) -version $(it_produces) "chatwork-cli/cw.*ver\."
	# argument error
	$(TARGET) -invalidflag $(omit); $(it_exit_with) 2
	# work with token
	HOME= CW_API_TOKEN=token $(TARGET) get me $(omit)
	unset CW_API_TOKEN ; $(TARGET) get me $(it_produces) "environment variable not set"
	# work with config
	$(TARGET) -f example.toml get me $(omit)
	$(TARGET) -f example.toml -p personal get me $(omit)
	# headers
	$(TARGET) -v -f example.toml get me $(it_produces) "^\(>\|<\)"


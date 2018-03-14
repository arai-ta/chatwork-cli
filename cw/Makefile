
TARGET = ./cw

# macro for test
ifdef DBG
omit =
else
omit = >/dev/null 2>&1
endif
it_produces  = 2>&1 | grep -q
it_exit_with = test $$? =
without_token = unset CW_API_TOKEN

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
	# messages
	$(TARGET) $(it_produces) "Simple command line tool"
	$(TARGET) -h $(it_produces) "Simple command line tool"
	$(TARGET) -version $(it_produces) "chatwork-cli/cw.*ver\."
	# argument error
	$(TARGET) -invalidflag $(omit); $(it_exit_with) 2
	# work with token
	HOME= CW_API_TOKEN=token $(TARGET) get me $(omit)
	$(without_token) ; $(TARGET) get me $(it_produces) "environment variable not set"
	# work with config
	$(TARGET) -f example.toml get me $(omit)
	$(TARGET) -f example.toml -p personal get me $(omit)
	$(without_token) ; $(TARGET) -f notexist.toml get me $(omit); $(it_exit_with) 1
	# headers
	$(TARGET) -v -f example.toml get me $(it_produces) "^\(>\|<\)"
	# values replacing
	$(TARGET) -v -f example.toml get rooms {room_id} messages $(it_produces) "GET.*/rooms/95297208/messages"


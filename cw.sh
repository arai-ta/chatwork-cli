#!/bin/sh

set -eu

CW_URL_BASE=https://api.chatwork.com/v2
CW_RC_FILE=~/.chatworkrc

if [ -f $CW_RC_FILE ]
then
    . $CW_RC_FILE
fi

if [ ! "${CW_TOKEN:-}" ]
then
    echo "Error: Token not defined" >&2
    exit 1
fi

get() {
    endpont=${1:?}
    curl --silent -H "X-ChatWorkToken: $CW_TOKEN" $CW_URL_BASE/$endpont
}

post() {
    endpont=${1:?}
    body=$2
    curl --verbose -X POST -H "X-ChatWorkToken: $CW_TOKEN" --data-urlencode "body=$body" $CW_URL_BASE/$endpont
}

get_my_chat() {
    get rooms | jq '.[] | select(.type == "my")'
}

read_line() {
    while read LINE
    do
        echo $LINE
    done
}

code() {
    echo "[code]"
    echo "$@"
    echo "[/code]"
}

my_room=`get_my_chat | jq .room_id`

post rooms/$my_room/messages "$(code $(read_line))"



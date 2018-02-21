chatwork-cli
============

`chatwork-cli` is simple command line tools for
[chatwork API](http://developer.chatwork.com).

## Install

    $ go get github.com/arai-ta/chatwork-cli/cw

## Usage

    $ cw
    # ==> Show usage and exit

    $ cw GET /me
    # ==> HTTP GET http://api.chatwork.com/v2/me

    $ cw GET /my/tasks
    # ==> HTTP GET http://api.chatwork.com/v2/my/tasks

    $ cw get my tasks   # alternative
    # ==> HTTP GET http://api.chatwork.com/v2/my/tasks

    $ cw POST rooms "name=New room for topic X"
    # ==> HTTP POST http://api.chatwork.com/v2/rooms

## Configuration

It works with chatwork API token.
(OAuth2 will be implemented in the near future)

### Environment Variable

    $ export CW_API_TOKEN=hereisyourapitoken

### Configuration File

    $ cp example.toml ~/.chatwork.toml
    $ vi ~/.chatwork.toml
    # edit it, like this: `token = hereisyourapitoken`

## Lisence

This software is released under the MIT License.


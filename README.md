chatwork-cli
============

[![CircleCI](https://circleci.com/gh/arai-ta/chatwork-cli.svg?style=shield)](https://circleci.com/gh/arai-ta/chatwork-cli)

`chatwork-cli` is a simple command line client for [chatwork API][1].

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

## Features

### Parameter Substitution

Edit `~/.chatwork.toml` file as following:

    [values]
    mychat = "17708368"

then you can do like this:

    $ cw post rooms {mychat} messages "body=I'm hungry:("
    # ==> HTTP POST https://api.chatwork.com/v2/rooms/17708368/messages

### Listing Available Endpoints

[chatwork API][1] is providing a [RAML definition][2].
`-endpoint` option will read definition and show list of available endpoints.

### Multiple Profiles

You can use API with multiple accounts by using configuration file.
See below.

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

[1]: http://developer.chatwork.com
[2]: https://github.com/chatwork/api

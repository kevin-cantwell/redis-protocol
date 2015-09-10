## About

The `redis-protocol` tool takes a stream of redis commands and reformats them according to the [redis protocol specification](http://redis.io/topics/protocol)

## Installation

`go get -u github.com/kevin-cantwell/redis-protocol`

## Usage

`$ echo "LPUSH mylist bar" | redis-protocol`

Output is encoded as:

```
*3\r\n$5\r\nLPUSH\r\n$6\r\nmylist\r\n$4\r\nbar\r\n
```

## Unix pipes

You can pipe the output into the redis command-line tool to take advantage of [redis pipelining](http://redis.io/topics/pipelining), which is a very efficient way to import large amounts of data into redis:

`$ echo "LPUSH mylist bar" | redis-protocol | redis-cli --pipe`
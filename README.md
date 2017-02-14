## About

This project is a Go package and cli tool called `resp` for encoding or decoding the [redis serialization protocol](http://redis.io/topics/protocol).

For example, you can pipe stdout to the redis-cli tool to take advantage of [redis pipelining](http://redis.io/topics/pipelining), which is a very efficient way to execute large amounts of commands against redis:


```
$ for i in {1..1000}; do echo -n "LPUSH counts $i" | resp; done | redis-cli --pipe
All data transferred. Waiting for the last reply...
Last reply received from server.
errors: 0, replies: 1000
```

## Installation

`go get -u github.com/kevin-cantwell/resp/...`

## CLI Usage

#### Encoding

The default usage is to encode stdin (_note that quoted strings are parsed as a single field_).

command:

```
$ resp
```

stdin:

```
SET foo "biz\nbaz buz"
```

stdout (escaped for readability):

```
*3\r\n$3\r\nSET\r\n$3\r\nfoo\r\n$11\r\nbiz\nbaz buz\r\n
```

#### Decoding

When decoding, the default is to use human-friendly output. The `--raw` flag may be passed in to get just the decoded text value and nothing more.

command:

```
$ resp --decode
```

stdin:

```
*4
+OK
-ERR foo
:123
$3
bar
```

stdout:

```
1) OK
2) (error) ERR foo
3) 123
4) "bar"
```
stdout with `--raw` option:

```
OK
ERR foo
123
bar
```

## Go Package

The Go Package provides reader and writer types. See Godocs: https://godoc.org/github.com/kevin-cantwell/resp
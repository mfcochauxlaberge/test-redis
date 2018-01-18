# test-redis

This is a simple test I made to see how much memory certain amounts of members take in a set. Also, how much time it takes to create and store them all.

## Usage

You first need to run Redis and install some Go dependencies.

```
go get github.com/dchest/uniuri
go get github.com/dustin/go-humanize
go get github.com/mediocregopher/radix.v3
```

The Redis instance should be listening on `127.0.0.1:6379` by default. Edit the URL in the source code if you have to.

Then you simply clone this package and run it. You may use the command `go run redis.go`. You may also use `go install` then simply `test-redis` if your `$GOPATH/bin` directory is included in your `PATH`.

## Example

The output should look like the following:

```
About to generate save 40,000,000 members into a redis set.
Redis instance flushed in 8.456815ms.
40 batches of 1,000,000 members built in 1m0.968611461s.
40,000,000 members generated in 1m2.463567922s.
40,000,000 members saved in 1m3.528827676s.
TOTAL: 40,000,000 members generated and saved in 1m3.540111529s
```

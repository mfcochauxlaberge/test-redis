package main

import (
	"fmt"
	"time"

	"github.com/dchest/uniuri"
	"github.com/dustin/go-humanize"
	"github.com/mediocregopher/radix.v3"
)

func main() {
	const NUMBER_OF_BATCHES = 40
	const BATCH_SIZE = 1000000
	const NUMBER_OF_MEMBERS = NUMBER_OF_BATCHES * BATCH_SIZE
	const MEMBER_SIZE = 8

	done := make(chan int, 10)
	members := make(chan string, NUMBER_OF_BATCHES*BATCH_SIZE)
	batchesEmpty := make(chan []string, 10)
	batchesFull := make(chan []string, 10)

	fmt.Printf("About to generate save %s members into a redis set.\n", humanize.Comma(NUMBER_OF_MEMBERS))

	start := time.Now()

	// Connection
	pool, err := radix.NewPool("tcp", "127.0.0.1:6379", 10, nil)
	if err != nil {
		panic(err)
	}

	// Generate members
	go func() {
		start := time.Now()

		for i := 0; i < NUMBER_OF_MEMBERS; i++ {
			members <- uniuri.NewLen(8)
		}

		fmt.Printf("%s members generated in %s.\n", humanize.Comma(NUMBER_OF_MEMBERS), time.Now().Sub(start))
	}()

	// Flush everything
	go func() {
		start := time.Now()

		pool.Do(radix.Cmd(nil, "FLUSHALL"))

		fmt.Printf("Redis instance flushed in %s.\n", time.Now().Sub(start))
		done <- 0
	}()

	// Build batches
	go func() {
		start := time.Now()

		batchesSent := 0
		var batch []string

		for i := 0; i < 2; i++ {
			batchesEmpty <- make([]string, BATCH_SIZE+1)
		}

	batching:
		for {
			select {
			case batch = <-batchesEmpty:
				batch[0] = "set"
				for i := 1; i < BATCH_SIZE+1; i++ {
					batch[i] = uniuri.NewLen(MEMBER_SIZE)
				}

				batchesFull <- batch

				batchesSent++
				if batchesSent == NUMBER_OF_BATCHES {
					break batching
				}
			}
		}

		fmt.Printf("%s batches of %s members built in %s.\n", humanize.Comma(NUMBER_OF_BATCHES), humanize.Comma(BATCH_SIZE), (time.Now().Sub(start)).String())
		close(batchesFull)
	}()

	// Wait for flushing
	<-done

	go func() {
		start := time.Now()

		var batch []string
		var more bool

	saving:
		for {
			select {
			case batch, more = <-batchesFull:
				if !more {
					break saving
				}
				err = pool.Do(radix.Cmd(nil, "SADD", batch...))
				batchesEmpty <- batch
			}
		}

		done <- 0
		fmt.Printf("%s members saved in %s.\n", humanize.Comma(NUMBER_OF_MEMBERS), (time.Now().Sub(start)).String())
	}()

	<-done

	fmt.Printf("TOTAL: %s members generated and saved in %s\n", humanize.Comma(NUMBER_OF_MEMBERS), time.Now().Sub(start))
}

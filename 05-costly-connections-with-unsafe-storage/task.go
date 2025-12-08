package main

import (
	"fmt"
	"log"
	"sync"
)

type Connection interface {
	// Need call Connect before Send
	// Take time to connect
	Connect()

	// Every connection should be disconnected after use
	// Take time to disconnect
	Disconnect()

	Send(req string) (string, error)
}

type ConnectionCreator interface {
	// Create new connection
	// Will return error if there is more than maxConn
	NewConnection() (Connection, error)
}

type Saver interface {
	// Saves data to unsafe storage
	// WILL CORRUPT DATA on concurrent save
	Save(data string)
}

// SendAndSave should send all requests concurrently using at most `maxConn` simultaneous connections.
// Responses must be saved using Saver.Save.
// Be careful: Saver.Save is not safe for concurrent use.
func SendAndSave(creator ConnectionCreator, saver Saver, requests []string, maxConn int) { //+return error?
	count := len(requests)
	requestCh := make(chan string, count)
	resultCh := make(chan string, count)
	errCh := make(chan error, count)

	//request gen
	go func() {
		for _, req := range requests {
			requestCh <- req
		}
		close(requestCh)
	}()

	//processing
	wg := &sync.WaitGroup{}
	for id := range maxConn {
		conn, err := creator.NewConnection()
		if err != nil {
			errCh <- err
			continue
		}
		wg.Add(1)
		go worker(id, wg, conn, requestCh, resultCh, errCh)
	}
	wg.Wait()
	close(resultCh)
	close(errCh)

	//save and print
	for res := range resultCh {
		saver.Save(res)
	}
	for err := range errCh {
		log.Println(err)
	}
}

func worker(id int, wg *sync.WaitGroup, conn Connection, requests <-chan string, result chan<- string, errs chan<- error) {
	defer wg.Done()
	defer func() {
		if r := recover(); r != nil {
			errs <- fmt.Errorf("worker %d: panic - %v", id, r)
		}
	}()

	conn.Connect()
	defer conn.Disconnect()

	for req := range requests {
		r, err := conn.Send(req)
		if err != nil {
			errs <- err
		} else {
			result <- r
		}
	}
}

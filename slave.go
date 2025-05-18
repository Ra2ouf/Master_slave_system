package main

import (
	"bufio"
	"distributed-db/shared"
	"distributed-db/slave"

	//"distributed-db/slave"
	"fmt"
	"os"
	"strings"
	"time"
)

func main() {
	fmt.Println("[Slave] Interactive shell started...")

	scanner := bufio.NewScanner(os.Stdin)

	fmt.Print("Enter Slave ID: ")
	scanner.Scan()
	slaveID := scanner.Text()

	fmt.Print("Enter DB Name: ")
	scanner.Scan()
	dbName := scanner.Text()

	// Start async replication in background
	go slave.StartAsyncReplication(slaveID, dbName)

	for {
		fmt.Print("SQL> ")
		if !scanner.Scan() {
			break
		}
		query := scanner.Text()

		if strings.ToLower(query) == "exit" {
			break
		}

		if strings.ToLower(query) == "run-buffer" {
			for _, req := range slave.AsyncBuffer {
				slave.ReplicateLocally(req)
			}
			slave.AsyncBuffer = nil
			continue
		}

		req := shared.QueryRequest{
			SlaveID:     slaveID,
			Query:       query,
			DBName:      dbName,
			IsReplicate: true,
		}

		go slave.SendQueryToMaster(req)
		go slave.ReplicateLocally(req)

		// Add to async buffer
		slave.AsyncBuffer = append(slave.AsyncBuffer, req)

		time.Sleep(1 * time.Second)
	}
}

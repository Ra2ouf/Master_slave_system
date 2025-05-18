package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"net"
	"os"
	"os/exec"
	"strings"

	"distributed-db/shared"
)

func handleConnection(conn net.Conn) {
	defer conn.Close()
	decoder := json.NewDecoder(conn)
	var req shared.QueryRequest

	if err := decoder.Decode(&req); err != nil {
		fmt.Println("[Master] Error decoding request:", err)
		return
	}

	fmt.Printf("[Master] Received query from %s: %s\n", req.SlaveID, req.Query)

	if strings.Contains(strings.ToUpper(req.Query), "DROP") {
		fmt.Println("[Master] DROP statements are not allowed from slaves.")
		return
	}

	executeQuery(req.Query, req.DBName)
}

func executeQuery(query, dbName string) {
	cmd := exec.Command("mysql", "-uroot", "-prootroot", "-e", fmt.Sprintf("USE %s; %s", dbName, query))
	output, err := cmd.CombinedOutput()

	fmt.Println("[Master] --- Query Result ---")
	fmt.Println(string(output))

	if err != nil {
		fmt.Printf("[Master] Error: %s\n", err.Error())
	}
}

func createDatabaseIfNotExists(dbName string) {
	cmd := exec.Command("mysql", "-uroot", "-prootroot", "-e", fmt.Sprintf("CREATE DATABASE IF NOT EXISTS %s;", dbName))
	output, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Printf("[Master] Failed to create database: %s\n", output)
	} else {
		fmt.Printf("[Master] Database '%s' is ready.\n", dbName)
	}
}
func startInteractiveShell() {
	scanner := bufio.NewScanner(os.Stdin)
	fmt.Println("[Master] Interactive shell started. Type your SQL queries:")
	fmt.Print("Enter database name: ")
	scanner.Scan()
	dbName := scanner.Text()

	// ✅ تأكد من وجود قاعدة البيانات أو أنشئها
	createDatabaseIfNotExists(dbName)

	for {
		fmt.Print("SQL> ")
		if !scanner.Scan() {
			break
		}
		query := scanner.Text()
		if strings.ToLower(query) == "exit" {
			break
		}
		go executeQuery(query, dbName)
	}
}
func main() {
	ln, err := net.Listen("tcp", ":9000")
	if err != nil {
		panic(err)
	}
	fmt.Println("[Master] Running on port 9000")

	// Start listening for slave connections
	go func() {
		for {
			conn, err := ln.Accept()
			if err != nil {
				fmt.Println("[Master] Connection error:", err)
				continue
			}
			go handleConnection(conn)
		}
	}()

	startInteractiveShell()
}

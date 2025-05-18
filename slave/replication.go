package slave

import (
	"distributed-db/shared"
	"encoding/json"
	"fmt"
	"net"
	"os/exec"
	"time"
)

var AsyncBuffer []shared.QueryRequest

func SendQueryToMaster(req shared.QueryRequest) {
	conn, err := net.Dial("tcp", "172.20.10.3:9000")
	if err != nil {
		fmt.Println("[Slave] Connection to master failed:", err)
		return
	}
	defer conn.Close()

	encoder := json.NewEncoder(conn)
	if err := encoder.Encode(req); err != nil {
		fmt.Println("[Slave] Encoding error:", err)
	} else {
		fmt.Println("[Slave] Sent query to master")
	}
}

func ReplicateLocally(req shared.QueryRequest) {
	if !req.IsReplicate {
		return
	}

	fmt.Println("[Slave] Replicating:", req.Query)
	cmd := exec.Command("mysql", "-uroot", "-prootroot", "-e", fmt.Sprintf("USE %s; %s", req.DBName, req.Query))
	output, err := cmd.CombinedOutput()
	fmt.Println("[slave] --- Query Result ---")
	fmt.Println(string(output))

	if err != nil {
		fmt.Printf("[Slave] Replication error: %s\n", output)
	} else {
		fmt.Println("[Slave] Replication successful")
	}
}

func StartAsyncReplication(slaveID, dbName string) {
	for {
		if len(AsyncBuffer) > 0 {
			fmt.Println("[Slave] Async Replicating batch...")

			for _, req := range AsyncBuffer {
				ReplicateLocally(req)
			}

			AsyncBuffer = nil
		}

		time.Sleep(10 * time.Second) // batch every 10 seconds
	}
}

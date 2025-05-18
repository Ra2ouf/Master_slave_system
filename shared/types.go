package shared

type QueryRequest struct {
	SlaveID     string `json:"slave_id"`
	Query       string `json:"query"`
	DBName      string `json:"db_name"`
	IsReplicate bool   `json:"is_replicate"`
}

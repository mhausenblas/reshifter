package handler

const (
	operationSuccess = "success"
	operationFail    = "fail"
)

// BackupRequest represents the request for a backup operation.
type BackupRequest struct {
	Endpoint   string `json:"endpoint"`
	Remote     string `json:"remote"`
	Bucket     string `json:"bucket"`
	Filter     string `json:"filter"`
	APIversion string `json:"apiversion"`
}

// BackupResult represents the results of a backup operation.
type BackupResult struct {
	Outcome  string `json:"outcome"`
	BackupID string `json:"backupid"`
}

// RestoreRequest represents the request for a restore operation.
type RestoreRequest struct {
	Endpoint string `json:"endpoint"`
	BackupID string `json:"backupid"`
	Remote   string `json:"remote"`
	Bucket   string `json:"bucket"`
}

// RestoreResult represents the results of a restore operation.
type RestoreResult struct {
	Outcome      string  `json:"outcome"`
	KeysRestored int     `json:"keysrestored"`
	ElapsedTime  float64 `json:"elapsedtimeinsec"`
}

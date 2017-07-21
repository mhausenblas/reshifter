package types

// KubernetesDistro represents a Kubernetes distribution.
type KubernetesDistro int

// BackupRequest represents the request for a backup operation.
type BackupRequest struct {
	Endpoint string `json:"endpoint"`
	Remote   string `json:"remote"`
	Bucket   string `json:"bucket"`
	Filter   string `json:"filter"`
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

// EtcdResponse represents the response of an etcd2 server at /version
// endpoint. Example:
//
//		{
//			"etcdserver": "2.3.8",
//			"etcdcluster": "2.3.0"
//		}
type EtcdResponse struct {
	EtcdServerVersion  string `json:"etcdserver"`
	EtcdClusterVersion string `json:"etcdcluster"`
}

// Reap function types take a path and a value and perform
// some action on it, for example, storing it to disk or
// writing it to stdout. The arg parameter is optional
// and can be used by the function in a context-dependent way,
// for example, it can specify a directory to write to.
type Reap func(path, value string, arg interface{}) error

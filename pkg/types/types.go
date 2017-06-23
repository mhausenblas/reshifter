package types

const (
	// EscapeColon represents the : in an etcd key
	EscapeColon = "ESC_COLON"
	// ContentFile is the name of the file an etcd value is stored
	ContentFile = "content"
)

// Endpoint represents an etcd server, available in
// a certain version at a certain URL.
type Endpoint struct {
	Version string `json:"version"`
	URL     string `json:"url"`
}

// BackupResult represents the results of a backup operation.
type BackupResult struct {
	Outcome  string `json:"outcome"`
	BackupID string `json:"backupid"`
}

// RestoreResult represents the results of a restore operation.
type RestoreResult struct {
	Outcome      string `json:"outcome"`
	KeysRestored int    `json:"keysrestored"`
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

// Reap function types take a node path and a value as parameters and performs
// some side effect, such as storing, on the node
type Reap func(string, string) error

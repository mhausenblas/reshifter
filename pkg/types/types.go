package types

const (
	// DefaultWorkDir is the default work directory for backups
	DefaultWorkDir = "/tmp"
	// EscapeColon represents the ':' character in an etcd key
	EscapeColon = "ESC_COLON"
	// ContentFile is the name of the file an etcd value is stored
	ContentFile = "content"
	// KubernetesPrefix represents the etcd key prefix for core Kubernetes
	KubernetesPrefix = "/kubernetes.io"
	// KubernetesPrefixLast represents a stop marker for core Kubernetes
	KubernetesPrefixLast = "/kubernetes.io/zzzzzzzzzz"
	// OpenShiftPrefix represents the etcd key prefix for OpenShift
	OpenShiftPrefix = "/openshift.io"
	// OpenShiftPrefixLast represents a stop marker for OpenShift
	OpenShiftPrefixLast = "/openshift.io/zzzzzzzzzz"
	// ContentTypeZip represents the content type for a ZIP file
	ContentTypeZip = "application/zip"
	// NotADistro represents the fact that no Kubernetes distro-related prefixes exit in etcd
	NotADistro KubernetesDistro = iota
	// Vanilla represents the vanilla, upstream Kubernetes distribution.
	Vanilla
	// OpenShift represents an OpenShift Kubernetes distribution.
	OpenShift
)

// KubernetesDistro represents a Kubernetes distribution.
type KubernetesDistro int

// BackupRequest represents the request for a backup operation.
type BackupRequest struct {
	Endpoint string `json:"endpoint"`
	Remote   string `json:"remote"`
	Bucket   string `json:"bucket"`
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

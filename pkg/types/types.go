package types

const (
	// DefaultWorkDir is the default work directory for backups
	DefaultWorkDir = "/tmp/reshifter"
	// UploadInMemoryBufferSize is the number of bytes in main memory used for uploading backups in the HTTP API handler (currently: 500kBytes)
	UploadInMemoryBufferSize = (1 << 10) * 500
	// EscapeColon represents the ':' character in an etcd key
	EscapeColon = "ESC_COLON"
	// ContentFile is the name of the file an etcd value is stored
	ContentFile = "content"
	// LegacyKubernetesPrefix represents the etcd legacy key prefix for core Kubernetes
	LegacyKubernetesPrefix = "/registry"
	// LegacyKubernetesPrefixLast represents a stop marker for the etcd legacy key prefix for core Kubernetes
	LegacyKubernetesPrefixLast = "/registry/zzzzzzzzzz"
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
	// ReapFunctionRaw represents the reap function that dumps the values of all keys to disk.
	ReapFunctionRaw = "raw"
	// ReapFunctionRender represents the reap function that dumps the values of all keys to stdout.
	ReapFunctionRender = "render"
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

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
	// EtcdAPIVersion2 represents the API v2 as of https://godoc.org/github.com/coreos/etcd/client
	EtcdAPIVersion2 = "v2"
	// EtcdAPIVersion3 represents the API v3 as of https://godoc.org/github.com/coreos/etcd/clientv3
	EtcdAPIVersion3 = "v3"
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
	// ReapFunctionFilter represents the reap function that dumps the values of certain keys to disk.
	ReapFunctionFilter = "filter"
	// NotADistro represents the fact that no Kubernetes distro-related prefixes exit in etcd
	NotADistro KubernetesDistro = iota
	// Vanilla represents the vanilla, upstream Kubernetes distribution.
	Vanilla
	// OpenShift represents an OpenShift Kubernetes distribution.
	OpenShift
)

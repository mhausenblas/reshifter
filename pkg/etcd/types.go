package etcd

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

// reap function types take a node path and a value as parameters and performs
// some side effect, such as storing, on the node
type reap func(string, string) error

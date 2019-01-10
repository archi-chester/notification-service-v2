package tcp

var (
	TCP_PORT int    = 10466
	TCP_IP   string = "10.46.2.197"
)

const (
	MESSAGE_TYPE_ERROR  = iota
	MESSAGE_TYPE_TEST   = iota
	MESSAGE_TYPE_INSERT = iota
)

// 	структурка для хранения мессаджа
type MessagePackage struct {
	Type       int    `json:"type"`
	SourcePort int    `json:"source_port"`
	SourceIP   string `json:"source_ip"`
	Message    []byte `json:"message"`
}

package v2

// So well-used "known" identities

func NewHostProcess(host, process string) map[string]string {
	return map[string]string{
		"kind": "host",
		"host": host,
		"process": process,
	}
}

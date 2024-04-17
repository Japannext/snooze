package opensearch

var (
  byTimestamp = map[string]string{"timestamp": "desc"}
)

func sorts(ss... map[string]string) []map[string]string {
  return ss
}

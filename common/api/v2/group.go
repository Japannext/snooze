package v2

type Group struct {
  // Hash value of the key-value
  Hash []byte `json:"hash"`
  // Human readable information about the group
  Kv KeyValue `json:"kv,omitempty"`
  // Timestamp of the last "hit"
  LastHit uint64 `json:"lastHit,omitEmpty"`
  // Value of the body of the last message found.
  LastMessage KeyValue `json:"lastMessage,omitempty"`
  // Number of hits between the oldest hit and the last updated
  Hits int `json:"hits,omitempty"`
  // Time of the first "hit". Used to give feedback to the user
  // about when the "hits" counter is counting.
  FirstHit uint64 `json:"firstHit,omitEmpty"`
}

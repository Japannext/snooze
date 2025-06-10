package hash

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"sort"

	"github.com/japannext/snooze/internal/models"
)

// Compute the hash of a map of string, in a predictable way.
func ComputeHash(m map[string]string) string {
	var keys []string
	h := md5.New()
	for k := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	for _, k := range keys {
		v := m[k]
		// Double hash, in order to avoid identical key/value concat collisions
		b := md5.Sum([]byte(k))
		h.Write(b[:])
		b = md5.Sum([]byte(v))
		h.Write(b[:])
	}

	return hex.EncodeToString(h.Sum(nil))
}

// Hash a single string
func hash(s string) string {
	h := md5.New()
	b := md5.Sum([]byte(s))
	h.Write(b[:])
	return hex.EncodeToString(h.Sum(nil))
}

func HashSelector(selector models.Selector) string {
	if len(selector.MatchLabels) != 0 {
		return fmt.Sprintf("labels:%s", ComputeHash(selector.MatchLabels))
	}
	if len(selector.MatchIdentity) != 0 {
		return fmt.Sprintf("identity:%s", ComputeHash(selector.MatchIdentity))
	}
	return ""
}

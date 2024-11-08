package opensearch

func pointer[V any](v V) *V {
	return &v
}

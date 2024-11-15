package mapping

var mappings = make(map[string]map[string]string)

func Startup(maps []*Mapping) {
	for _, m := range maps {
		mappings[m.Name] = m.Map
	}
}

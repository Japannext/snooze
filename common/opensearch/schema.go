package opensearch

var alertIndex = IndexTemplate{
	Name:         "alerts-v2",
	IndexPattern: []string{"alerts-v2-*"},
	Template: Indice{
		Settings: settings,
		Mappings: []IndexMapping{
			{
				Properties: map[string]MappingProps{
					"source.kind": {Type: "keyword"},
					"source.name": {Type: "keyword"},
					"timestamp":   {Type: "unsigned_long"},
					"groupHash":   {Type: "byte"},
					"groupLabels": {Type: "object"},
					"labels":      {Type: "object"},
					"attributes":  {Type: "object"},
					"body":        {Type: "object"},
					"mute": {
						Type: "object",
						Fields: map[string]MappingProps{
							"enabled":         {Type: "boolean"},
							"component":       {Type: "keyword"},
							"rule":            {Type: "text"},
							"skipNotificaton": {Type: "boolean"},
						},
					},
				},
			},
		},
	},
}

var groupIndex = IndexTemplate{
	Name:         "group-v2",
	IndexPattern: []string{"group-v2-*"},
	Template: Indice{
		Settings: settings,
		Mappings: []IndexMapping{
			{
				Properties: map[string]MappingProps{
					"hash":     {Type: "byte"},
					"labels":   {Type: "object"},
					"hits":     {Type: "integer"},
					"lastBody": {Type: "object"},
					"lastHit":  {Type: "integer"},
					"firstHit": {Type: "unsigned_long"},
				},
			},
		},
	},
}

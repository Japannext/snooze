package opensearch

import (
	"context"
)

/*
var alertIndex = IndexTemplate{
	Name:         "alerts-v2",
	IndexPattern: []string{"alerts-v2-*"},
	Template: Indice{
		Settings: settings,
		Mappings: IndexMapping{
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
}

var logIndex = IndexTemplate{
	Name:         "log-v2",
	IndexPattern: []string{"log-v2-*"},
	Template: Indice{
		Settings: settings,
		Mappings: IndexMapping{
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
}

var groupIndex = IndexTemplate{
	Name:         "group-v2",
	IndexPattern: []string{"group-v2-*"},
	Template: Indice{
		Settings: settings,
		Mappings: IndexMapping{
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
}
*/

func (client *OpensearchLogStore) Bootstrap() error {
	ctx := context.Background()
	log.Info("Bootstrapping opensearch...")

	numberOfShards := 1
	numberOfReplicas := 2
	settings := IndexSettings{numberOfShards, numberOfReplicas}

	logTemplate := IndexTemplate{
		IndexPattern: []string{"log-v2", "log-v2-*"},
		DataStream: DataStream{TimestampField{"timestampMillis"}},
		Template: Indice{
			Settings: settings,
			Mappings: IndexMapping{
				Properties: map[string]MappingProps{
					"timestampMillis": {Type: "date", Format: "epoch_millis"},
					"source.kind": {Type: "keyword"},
					"source.name": {Type: "keyword"},
					"identity": {Type: "object"},
					"groupHash":   {Type: "keyword"},
					"groupLabels": {Type: "object"},
					"profile": {Type: "keyword"},
					"pattern": {Type: "keyword"},
					"labels":      {Type: "object"},
					"message":        {Type: "text"},
					"mute.skipNotification": {Type: "boolean"},
					"mute.skipStorage": {Type: "boolean"},
					"mute.reason": {Type: "keyword"},
				},
			},
		},
	}
	client.ensureIndex(ctx, "log-v2", logTemplate)
	client.ensureDatastream(ctx, "log-v2")

	notificationTemplate := IndexTemplate{
		IndexPattern: []string{"notification-v2"},
		DataStream: DataStream{TimestampField{"timestampMillis"}},
		Template: Indice{
			Settings: settings,
			Mappings: IndexMapping{
				Properties: map[string]MappingProps{
					"timestampMillis": {Type: "date", Format: "epoch_millis"},
					"destination.kind": {Type: "keyword"},
					"destination.name": {Type: "keyword"},
					"alertUID": {Type: "keyword"},
					"logUID": {Type: "keyword"},
					"body": {Type: "object"},
				},
			},
		},
	}

	client.ensureIndex(ctx, "notification-v2", notificationTemplate)
	client.ensureDatastream(ctx, "notification-v2")

	log.Info("Done bootstrapping indices")
	return nil
}

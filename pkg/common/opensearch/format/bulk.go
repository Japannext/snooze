package format

type Action = string

const (
	INDEX_ACTION  Action = "index"
	UPDATE_ACTION Action = "update"
	CREATE_ACTION Action = "create"
	DELETE_ACTION Action = "delete"
)

type BulkHeader = map[Action]Metadata

type Metadata struct {
	ID string `json:"_id,omitempty"`
	Index string `json:"_index,omitempty"`
}

// UPDATE action

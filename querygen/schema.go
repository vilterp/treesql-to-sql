package querygen

var Timestamps = map[string]struct{}{"created_at": {}, "updated_at": {}}

type TableDesc struct {
	Name       string
	Cols       []ColumnDesc
	PrimaryKey string
	Timestamps map[string]struct{}
}

type ColumnDesc struct {
	Name string
	Type ColType
}

type ColType int

const (
	TInt = iota
	TString
	TTimestamp
	TInet
)

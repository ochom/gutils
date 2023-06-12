package nosql

// Model is a struct that implements this interface
type Model interface {
	// TableName returns the table name
	TableName() string
}

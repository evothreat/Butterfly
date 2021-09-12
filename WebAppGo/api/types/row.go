package types

type Row interface {
	Scan(...interface{}) error
}

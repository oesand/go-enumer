package ifaces

type RowScanner interface {
	Scan(dest ...any) error
}

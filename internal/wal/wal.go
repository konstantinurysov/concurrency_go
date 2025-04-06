package wal

type WAL interface {
	Append(data []byte) error
	Close() error
}

type wal struct {
	walStorage WAL
}

func New(walStorage WAL) WAL {
	return &wal{
		walStorage: walStorage,
	}
}

func (w *wal) Append(data []byte) error {
	data = append(data, '\n')
	return w.walStorage.Append(data)
}

func (w *wal) Close() error {
	// Implement this
	return nil
}

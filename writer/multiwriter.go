package writer

type MultiWriter struct {
	writers []Writer
}

func NewMultiWriter(writers ...Writer) *MultiWriter {
	return &MultiWriter{writers: writers}
}

func (m *MultiWriter) Write(b []byte) {
	for _, w := range m.writers {
		w.Write(b)
	}
}

func (m *MultiWriter) Close() {
	for _, w := range m.writers {
		w.Close()
	}
}

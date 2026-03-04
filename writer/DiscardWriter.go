package writer

type DiscardWriter struct{}

func (d *DiscardWriter) Write(b []byte) {}

func (d *DiscardWriter) Close() {}

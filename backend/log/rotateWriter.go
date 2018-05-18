package log

import (
	"os"
	"sync"
	"time"
)

// RotateWriter is a custom implementation of io.Writer to rotate logs after some time
type RotateWriter struct {
	lock     sync.Mutex
	filename string
	fd       *os.File
	ticker   *time.Ticker
}

func newRotateWriter(filename string) *RotateWriter {
	w := new(RotateWriter)
	w.filename = filename
	w.ticker = time.NewTicker(time.Hour * 24)
	if err := w.rotate(); err != nil {
		panic(err)
	}

	go rotateService(w)

	return w
}

func rotateService(w *RotateWriter) {
	for range w.ticker.C {
		if err := w.rotate(); err != nil {
			panic(err)
		}
	}
}

func (w *RotateWriter) Write(output []byte) (int, error) {
	w.lock.Lock()
	defer w.lock.Unlock()
	return w.fd.Write(output)
}

func (w *RotateWriter) rotate() (err error) {
	w.lock.Lock()
	defer w.lock.Unlock()

	if w.fd != nil {
		err = w.fd.Close()
		w.fd = nil
		if err != nil {
			return err
		}
	}

	if _, err = os.Stat(w.filename); err == nil || os.IsExist(err) {
		newf := w.filename + "." + time.Now().Format("2006-01-02_15-04-05")
		if err := os.Rename(w.filename, newf); err != nil {
			return err
		}
	}

	w.fd, err = os.Create(w.filename)
	if err != nil {
		return err
	}
	return nil
}

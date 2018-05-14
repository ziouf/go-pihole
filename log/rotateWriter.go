package log

import (
	"os"
	"sync"
	"time"
)

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

func (w *RotateWriter) rotate() error {
	w.lock.Lock()
	defer w.lock.Unlock()

	if w.fd != nil {
		err := w.fd.Close()
		w.fd = nil
		if err != nil {
			return err
		}
	}

	if _, err := os.Stat(w.filename); err == nil || os.IsExist(err) {
		newf := w.filename + "." + time.Now().Format(time.RFC3339)
		if err := os.Rename(w.filename, newf); err != nil {
			return err
		}
	}

	if fd, err := os.Create(w.filename); err != nil {
		return err
	} else {
		w.fd = fd
	}
	return nil
}

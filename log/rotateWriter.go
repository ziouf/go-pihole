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
	w.ticker = time.NewTicker(time.Second * 10)
	if err := w.rotate(); err != nil {
		panic(err)
	}

	go func(w *RotateWriter) {
		for range w.ticker.C {
			if err := w.rotate(); err != nil {
				println(err)
			}
		}
	}(w)

	return w
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
		println("Close FD")
		err = w.fd.Close()
		w.fd = nil
		if err != nil {
			return
		}
	}

	if _, err = os.Stat(w.filename); os.IsExist(err) {
		println("Rename file")
		err = os.Rename(w.filename, w.filename+"."+time.Now().Format(time.RFC3339))
		if err != nil {
			return
		}
	} else {
		println(w.filename, "does not exists")
	}

	w.fd, err = os.Create(w.filename)
	return
}

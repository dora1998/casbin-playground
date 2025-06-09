package main

import (
	"log"

	"github.com/casbin/casbin/v2/persist"
)

type Watcher struct {
	cb     func(string)
	update chan string
	close  chan struct{}
}

func NewWatcher(updateCh chan string) (persist.Watcher, error) {
	w := &Watcher{
		update: updateCh,
	}

	w.subscribe()

	return w, nil
}

func (w *Watcher) subscribe() {
	go func() {
		for {
			select {
			case <-w.close:
				return
			case <-w.update:
				log.Println("Received update message")
				// 引数の文字列はなんでもいいっぽい
				// https://github.com/casbin/casbin/blob/c6f6cfcd1a0b22667290b2aba93290a5ee78ca5d/enforcer.go#L274
				w.cb("update")
			}
		}
	}()
}

func (w *Watcher) Close() {
	close(w.close)
}

func (w *Watcher) SetUpdateCallback(cb func(string)) error {
	w.cb = cb
	return nil
}

func (w *Watcher) Update() error {
	return nil
}

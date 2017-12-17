
package dispatcher

import (
    "sync"
    "errors"
)

type Worker interface {
    Init(id int)
    Proc(interface{})
}

type Dispatcher struct {
    queue   chan interface{}
    quit    []chan bool
    workers []Worker
    wg      sync.WaitGroup
    stop    bool
}

type WorkerFunc func(id int) Worker

func NewDispatcher(queues, works int, wf WorkerFunc) (*Dispatcher, error) {
    if queues <= 0 || works <= 0 {
        return nil, errors.New("queue and works non zero")
    }

    d := &Dispatcher{
        queue: make(chan interface{}, queues),
    }
    for i := 0; i < works; i++ {
        w := wf(i)
        w.(Worker).Init(i)
        d.workers = append(d.workers, w)
        d.quit = append(d.quit, make(chan bool))
    }
    return d, nil
}

func (d *Dispatcher) Start() {
    for n, w := range d.workers {
        go func(n int, w interface{}) {
            for {
                select {
                case v := <- d.queue:
                    w.(Worker).Proc(v)
                    d.wg.Done()
                case <- d.quit[n]:
                    return
                }
            }
        }(n, w)
    }
}

func (d *Dispatcher) Add(v interface{}) {
    d.wg.Add(1)
    d.queue <- v
}

func (d *Dispatcher) Wait() {
    if !d.stop {
        d.wg.Wait()
    }
}

func (d *Dispatcher) Stop() {
    d.stop = true
    for i, _ := range d.quit {
        go func(ii int) {
            d.quit[ii] <- true
        }(i)
    }
}

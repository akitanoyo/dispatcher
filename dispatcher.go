
package dispatcher

import (
    "sync"
    "errors"
)

type Worker interface {
    Init(id int)
    Proc(interface{})
}

// type WorkerT struct {
//     Worker
// }

type Dispatcher struct {
    queue   chan interface{}
    workers []Worker
    wg      sync.WaitGroup
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
        // fmt.Println(w)
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
    d.wg.Wait()
}

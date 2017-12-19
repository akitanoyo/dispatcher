
package dispatcher

import (
    // "fmt"
    "sync"
    "errors"
)

type Worker interface {
    Init(id int)
    Proc(interface{}) interface{}
}

type Dispatcher struct {
    queue   chan interface{}
    quit    []chan bool
    workers []Worker
    res     chan interface{}
    wg      sync.WaitGroup
    stop    bool
    ended   bool
}

type WorkerFunc func(id int) Worker

func NewDispatcher(queues, works int, wf WorkerFunc) (*Dispatcher, error) {
    if queues <= 0 || works <= 0 {
        return nil, errors.New("queue and works non zero")
    }

    d := &Dispatcher{
        queue: make(chan interface{}, queues),
        res  : make(chan interface{}, works),
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
            f := true
            for {
                select {
                case qf := <- d.quit[n]:
                    if qf {
                        // fmt.Printf("quit end %d\n", n)
                        return 
                    } else {
                        f = false
                        // fmt.Printf("quit wait %d\n", n)
                    }
                case v := <- d.queue:
                    if f {
                        r := w.(Worker).Proc(v)
                        if r != nil {
                            // fmt.Printf("-1-[%d]\n", n)
                            d.res <- r
                            // fmt.Printf("-2-[%d]\n", n)
                        }
                    }
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
    if !d.stop {
        d.wg.Wait()
        d.ended = true
    }
}

type ResultFunc func(v interface{}) error

func (d *Dispatcher) ResultWait(rsf ResultFunc) {
    once := sync.Once{}
    e := make(chan bool)
    
    go func() {
        // d.wg.Wait()
        d.Wait()
        e <- true
    }()

    for {
        select {
        case r := <- d.res:
            err := rsf(r)
            if err != nil {
                once.Do(func() {
                    // d.sendstop(false) - block -> deadlock
                    // non block
                    d.stop = true
                    for i, _ := range d.quit {
                        go func(ii int) {
                            d.quit[ii] <- false
                        }(i)
                    }
                })
            }
        case <- e:
            return
        }
    }
}

func (d *Dispatcher) sendstop(f bool) {
    wg  := sync.WaitGroup{}
    d.stop = true
    for i, _ := range d.quit {
        wg.Add(1)
        go func(ii int) {
            // fmt.Printf("stop quit %d (force %v)\n", ii, f)
            d.quit[ii] <- f
            // fmt.Printf("stop quited %d (force %v)\n", ii, f)
            wg.Done()
        }(i)
    }
    wg.Wait()
}

func (d *Dispatcher) Close() {
    if !d.ended {
        d.sendstop(false)
        d.Wait()
    }

    d.sendstop(true)
}

package dispatcher

type Worker interface {  
	Init(id int)  
	Proc(interface{}) interface{}  
}  

type WorkerFunc func(id int) Worker  

func NewDispatcher(queues, works int, wf WorkerFunc) (*Dispatcher, error)  

func (d *Dispatcher) Start()  
func (d *Dispatcher) Add(v interface{})  
func (d *Dispatcher) Wait()  
func (d *Dispatcher) Close()  

type ResultFunc func(v interface{}) error  
func (d *Dispatcher) ResultWait(rsf ResultFunc)  


example
```
package main

import (
    "fmt"
    "github.com/akitanoyo/dispatcher"
)

type MyWorker struct {
    dispatcher.Worker

    // user members
    id int
}

func (w *MyWorker) Init(id int) {
    fmt.Printf("test Init() %d\n", id)
    w.id = id
}

func (w *MyWorker) Proc(v interface{}) interface{} {
    fmt.Printf("test Proc(%d) %v\n", w.id, v)
    return nil
}

func main() {
    // queue 100, workers 5
    d, err := dispatcher.NewDispatcher(100, 5, func(id int) dispatcher.Worker {
        w := &MyWorker{}
        return w
    })
    if err != nil {
        panic(err)
    }

    d.Start()

    for i := 0; i < 100; i++ {
        // send worker
        d.Add(fmt.Sprintf("test %d", i))
    }

    d.Wait()
    d.Close()
}
```

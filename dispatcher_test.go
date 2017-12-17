
package dispatcher

import "testing"
import "fmt"
import "time"
import "math/rand"

type MyWorker struct {
    Worker
    // user members
    id int
}

func (w *MyWorker) Init(id int) {
    fmt.Printf("test Init() %d\n", id)
    w.id = id
}

func (w *MyWorker) Proc(v interface{}) {
    fmt.Printf("test Proc(%d) %v\n", w.id, v)
    n := rand.Intn(100)
    time.Sleep(time.Duration(n) * time.Millisecond)
}

func TestDispatcher(t *testing.T) {
    d, err := NewDispatcher(100, 5, func(id int) Worker {
        return &MyWorker{/* init members */}
    })
    if err != nil {
        t.Errorf("%v\n", err)
    }

    d.Start()

    for i := 0; i < 100; i++ {
        d.Add(fmt.Sprintf("test %d", i))
    }

    d.Wait()
}


func TestDispatcherStop(t *testing.T) {
    d, err := NewDispatcher(100, 5, func(id int) Worker {
        return &MyWorker{/* init members */}
    })
    if err != nil {
        t.Errorf("%v\n", err)
    }

    d.Start()

    for i := 0; i < 100; i++ {
        d.Add(fmt.Sprintf("test %d", i))
    }

    d.Stop()
    time.Sleep(300 * time.Millisecond)
    d.Wait() // 無くていいが 呼ばれても問題ないようにした
}


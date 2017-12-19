
package dispatcher

import "testing"
import "fmt"
import "time"
import "math/rand"
import "errors"
import "runtime"
import "os"

type MyWorker struct {
    Worker
    // user members
    id int
}

func (w *MyWorker) Init(id int) {
    fmt.Printf("test Init() %d\n", id)
    w.id = id
}

func (w *MyWorker) Proc(v interface{}) interface{} {
    fmt.Printf("test Proc(%d) %v\n", w.id, v)
    n := rand.Intn(100)
    time.Sleep(time.Duration(n) * time.Millisecond)

    return nil
}

func TestDispatcher(t *testing.T) {
    goroutinesPrint(t, "start", 0)

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
    d.Close()
    
    gn := goroutinesPrint(t, "end", 300)
    if gn > 2 {
        t.Errorf("goroutine not ended remaining %d\n", gn - 2)
    }
}


func TestDispatcherStop(t *testing.T) {
    goroutinesPrint(t, "start", 0)

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

    d.Close()

    gn := goroutinesPrint(t, "end", 300)
    if gn > 2 {
        t.Errorf("goroutine not ended remaining %d\n", gn - 2)
    }
}

type MyResWorker struct {
    Worker  // interface
    // user members
    id int
    // res  chan string
}

func (w *MyResWorker) Init(id int) {
    fmt.Printf("test Init() %d\n", id)
    w.id = id
}

type ResWork struct {
    nn int
    mes string
    req string
}


func (w *MyResWorker) Proc(v interface{}) interface{} {
    fmt.Printf("test Proc(%d) %v\n", w.id, v)
    n := rand.Intn(100)
    time.Sleep(time.Duration(n) * time.Millisecond)

    return ResWork {
        nn: n,
        mes: "myresw successful",
        req: v.(string),
    }
}

func TestDispatcherRes(t *testing.T) {
    goroutinesPrint(t, "start", 0)

    d, err := NewDispatcher(100, 5, func(id int) Worker {
        w := &MyResWorker{}
        return w
    })
    if err != nil {
        t.Errorf("%v\n", err)
    }

    d.Start()

    for i := 0; i < 100; i++ {
        d.Add(fmt.Sprintf("test %d", i))
    }
    
    ct := 0
    d.ResultWait(func(r interface{}) error {
        s := r.(ResWork) // Proc return type
        fmt.Printf("%s[%d](%s) (%d)\n", s.mes, s.nn, s.req, ct)
        ct++
        if ct > 50 {
            return errors.New("stop")
        }
        return nil
    })

    d.Close()
    
    gn := goroutinesPrint(t, "end", 300)
    if gn > 2 {
        t.Errorf("goroutine not ended remaining %d\n", gn - 2)
    }
}

func TestMain(m *testing.M) {
    n := runtime.NumCPU()
	runtime.GOMAXPROCS(n)
    fmt.Printf("CPUs %d\n", n)

    os.Exit(m.Run())
}

func goroutinesPrint(t *testing.T, mes string, delay int) (gn int) {
    if delay > 0 {
        time.Sleep(time.Duration(delay) * time.Millisecond)
    }
    gn = runtime.NumGoroutine()
    fmt.Printf("==== goroutines %d (%s.%s)\n", gn, t.Name(), mes)
    return
}


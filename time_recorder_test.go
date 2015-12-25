package pprof

import (
    // "bytes"
    "strconv"
    "testing"
    "time"
)

func Test_TimeRecorder(t *testing.T) {
    recorder := NewTimeRecorderDefault()

    for i := 0; i < 10; i++ {
        t := time.Now()
        doSomething(100)
        recorder.AddRecord("doSomething(100)", time.Since(t))
    }

    for i := 0; i < 10; i++ {
        t := time.Now()
        doSomething(1000)
        recorder.AddRecord("doSomething(1000)", time.Since(t))
    }

    // for i := 0; i < 10; i++ {
    // 	t := time.Now()
    // 	doSomething(10000)
    // 	recorder.AddRecord("doSomething(10000)", time.Since(t))
    // }

    // for i := 0; i < 10; i++ {
    // 	t := time.Now()
    // 	doSomething(100000)
    // 	recorder.AddRecord("doSomething(100000)", time.Since(t))
    // }

    // buffer := new(bytes.Buffer)
    // recorder.WriteCSV(buffer)
    // println()
    // println(buffer.String())
    // println()

    recorder.PrintOut()
}

func doSomething(n int) int {
    m := 0
    for i := 0; i < n; i++ {
        m += i
        for j := 0; j < 30; j++ {
            strconv.Itoa(m) // cost some CPU time
        }
    }
    return m
}

package pprof

import (
    "bufio"
    "fmt"
    "io"
    "os"
    "sort"

    // "sync"
    // "sync/atomic"
    "time"
)

var (
    DEFAULT_CACHE_SIZE = 4096
)

type TimeRecorder struct {
    // mutex   sync.RWMutex
    records   map[string]*timeRecord
    cacheChan chan map[string]time.Duration
}

func NewTimeRecorderDefault() *TimeRecorder {
    return NewTimeRecorder(DEFAULT_CACHE_SIZE)
}

func NewTimeRecorder(size int) *TimeRecorder {
    recorder := &TimeRecorder{
        records:   make(map[string]*timeRecord),
        cacheChan: make(chan map[string]time.Duration, size),
    }
    go func() {
        for {
            select {
            case v := <-recorder.cacheChan:
                recorder.addTimeStamp(v)
            }
        }
    }()
    return recorder
}
func (t *TimeRecorder) AddRecord(name string, duration time.Duration) {
    m := make(map[string]time.Duration)
    m[name] = duration
    t.cacheChan <- m
}

func (tr *TimeRecorder) SaveCSV(filename string) error {
    file, err := os.Create(filename)
    if err != nil {
        return err
    }
    defer file.Close()

    return tr.writeCSV(file)
}

func (t *TimeRecorder) PrintOut() (l []string) {

    records := t.getRecords()

    l = append(l, fmt.Sprintf("%20s  %20s %20s %20s  %20s  %20s  %20s", "name", "times", "median", "avg", "min", "max", "total"))
    for _, record := range records {
        l = append(l, fmt.Sprintf("%20s  %20d %20.0f  %20.0f  %20.0f  %20.0f  %20.0f",
            record.Name, record.Times, record.MedianTime, record.AvgUsedTime, record.MinUsedTime, record.MaxUsedTime, record.TotalUsedTime))
    }

    for _, s := range l {
        fmt.Println(s)
    }
    return
}

func (t *TimeRecorder) addTimeStamp(m map[string]time.Duration) {
    for name, usedTime := range m {
        t.setRecord(name, usedTime)
    }
}

func (t *TimeRecorder) setRecord(name string, usedTime time.Duration) {
    usedNano := usedTime.Nanoseconds()

    var exists bool
    var r *timeRecord
    r, exists = t.records[name]
    if !exists {
        r = NewRecord(name)
        t.records[name] = r
    }
    r.AddRecord(usedNano)
}

func (tr *TimeRecorder) writeCSV(writer io.Writer) error {
    results := tr.getRecords()

    buf := bufio.NewWriter(writer)

    if _, err := fmt.Fprintln(writer, "name,times,avg,min,max,total"); err != nil {
        return err
    }

    for _, r := range results {
        if _, err := fmt.Fprintf(writer,
            "%s,%d,%d,%d,%d,%d\n",
            r.Name,
            r.Times,
            r.AvgUsedTime,
            r.MinUsedTime,
            r.MaxUsedTime,
            r.TotalUsedTime,
        ); err != nil {
            return err
        }
    }

    return buf.Flush()
}

func (tr *TimeRecorder) getRecords() timeRecordList {
    list := timeRecordList{}
    for _, record := range tr.records {
        list = append(list, record)
    }
    list.summary()
    sort.Sort(list)
    return list
}

// type sortTimeRecord struct {
//     Name          string
//     Times         int64
//     AvgUsedTime   int64
//     MinUsedTime   int64
//     MaxUsedTime   int64
//     TotalUsedTime int64
// }

// type sortTimeRecords []*sortTimeRecord

// func (this sortTimeRecords) Len() int {
//     return len(this)
// }

// func (this sortTimeRecords) Swap(i, j int) {
//     this[i], this[j] = this[j], this[i]
// }

// func (this sortTimeRecords) Less(i, j int) bool {
//     return this[i].AvgUsedTime > this[j].AvgUsedTime || (this[i].AvgUsedTime == this[j].AvgUsedTime && this[i].Times < this[j].Times)
// }

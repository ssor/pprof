package pprof

import (
    "sort"
)

type timeRecord struct {
    Name          string
    Times         int
    TotalUsedTime float64
    MaxUsedTime   float64
    MinUsedTime   float64
    AvgUsedTime   float64
    MedianTime    float64
    History       []float64
}

func NewRecord(name string) *timeRecord {
    return &timeRecord{
        Name:    name,
        History: []float64{},
    }
}

func (t *timeRecord) AddRecord(period int64) {
    t.Times += 1
    t.History = append(t.History, float64(period))
}

func (t *timeRecord) summary() {
    if t.Times <= 0 {
        return
    }
    if len(t.History) <= 0 {
        return
    }

    sort.Float64s(t.History)
    t.MinUsedTime = t.History[0]
    t.MaxUsedTime = t.History[len(t.History)-1]
    t.TotalUsedTime = sumFload64(t.History)
    t.AvgUsedTime = t.TotalUsedTime / float64(t.Times)
    t.MedianTime = t.History[t.Times/2]
}

func sumFload64(l []float64) float64 {
    var sum float64 = 0
    for _, f := range l {
        sum += f
    }
    return sum
}

type timeRecordList []*timeRecord

func (t timeRecordList) summary() {
    for _, record := range t {
        record.summary()
    }
}

func (this timeRecordList) Len() int {
    return len(this)
}

func (this timeRecordList) Swap(i, j int) {
    this[i], this[j] = this[j], this[i]
}

func (this timeRecordList) Less(i, j int) bool {
    return this[i].AvgUsedTime > this[j].AvgUsedTime || (this[i].AvgUsedTime == this[j].AvgUsedTime && this[i].Times < this[j].Times)
}

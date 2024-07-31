package vnstat

import "time"

type Report struct {
	VNStatVersion string   `json:"vnstatversion"`
	JSONVersion   string   `json:"jsonversion"`
	IFaces        []*IFace `json:"interfaces"`
}

type IFace struct {
	Name    string           `json:"name"`
	Alias   string           `json:"alias"`
	Created *DateTime        `json:"created"`
	Updated *DateTime        `json:"updated"`
	Traffic *TrafficOverview `json:"traffic"`
}

type Date struct {
	Year  int `json:"year"`
	Month int `json:"month,omitempty"`
	Day   int `json:"day,omitempty"`
}

type Time struct {
	Hour   int `json:"hour"`
	Minute int `json:"minute"`
}

type DateTime struct {
	Date      *Date `json:"date"`
	Time      *Time `json:"time"`
	Timestamp int64 `json:"timestamp"`
}

func (dt *DateTime) TimeParsed() time.Time {
	return time.Unix(dt.Timestamp, 0)
}

type TrafficOverview struct {
	Total      *Traffic   `json:"total"`
	FiveMinute []*Traffic `json:"fiveminute"`
	Hour       []*Traffic `json:"hour"`
	Day        []*Traffic `json:"day"`
	Month      []*Traffic `json:"month"`
	Year       []*Traffic `json:"year"`
	Top        []*Traffic `json:"top"`
}

type Traffic struct {
	DateTime
	ID int64  `json:"id"`
	Rx uint64 `json:"rx"`
	Tx uint64 `json:"tx"`
}

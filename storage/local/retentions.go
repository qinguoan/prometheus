package local

import (
	"fmt"
	"math"
	"regexp"
	"strconv"
	"strings"

	"github.com/prometheus/common/model"

	"time"

	"github.com/qinguoan/prometheus/storage/local/chunk"
)

const (
	DefaultPointBytes  = 8
	DefaultLimitFactor = 0.5
	DefaultRetetions   = "1s:3d,10s:30d,60s:180d,300s:900d,1800s:5400d,3h:32400d"
	TestRetentions     = "1s:5m,10s:50m,60s:300m"
)

const (
	Seconds = 1
	Minutes = 60
	Hours   = 3600
	Days    = 86400
	Weeks   = 86400 * 7
	Years   = 86400 * 365
)

var DefaultRetentionList Retentions

func unitMultiplier(s string) (int64, error) {
	switch {
	case strings.HasPrefix(s, "s"):
		return Seconds, nil
	case strings.HasPrefix(s, "m"):
		return Minutes, nil
	case strings.HasPrefix(s, "h"):
		return Hours, nil
	case strings.HasPrefix(s, "d"):
		return Days, nil
	case strings.HasPrefix(s, "w"):
		return Weeks, nil
	case strings.HasPrefix(s, "y"):
		return Years, nil
	}
	return 0, fmt.Errorf("Invalid unit multiplier [%v]", s)
}

type Retentions []*Retention

func (r Retentions) Less(i, j int) bool {
	return r[i].interval < r[j].interval
}

func (r Retentions) Len() int {
	return len(r)
}

func (r Retentions) Swap(i, j int) {
	r[i], r[j] = r[j], r[i]
}

var retentionRegexp *regexp.Regexp = regexp.MustCompile("^(\\d+)([smhdwy]+)$")

func parseRetentionPart(rp string) (int64, error) {
	part, err := strconv.ParseInt(rp, 10, 64)
	if err == nil {
		return part, nil
	}
	if !retentionRegexp.MatchString(rp) {
		return 0, fmt.Errorf("%v", rp)
	}
	matches := retentionRegexp.FindStringSubmatch(rp)
	value, err := strconv.ParseInt(matches[1], 10, 64)
	if err != nil {
		return 0, fmt.Errorf("Regex on %v is borked, %v cannot be parsed as int", rp, matches[1])
	}
	multiplier, err := unitMultiplier(matches[2])
	return multiplier * value, err
}

func parseRetentionDef(r string) (*Retention, error) {
	parts := strings.Split(r, ":")
	if len(parts) != 2 {
		return nil, fmt.Errorf("Not enough parts in retentionDef [%v]", r)
	}

	interval, err := parseRetentionPart(parts[0])
	if err != nil {
		return nil, err
	}

	persist, err := parseRetentionPart(parts[1])
	if err != nil {
		return nil, err
	}

	return New(interval, persist), nil
}

func ParseRetentionDefs(rs string) (Retentions, error) {
	retentions := make(Retentions, 0)
	rss := strings.Split(rs, ",")
	if len(rss) >= math.MaxUint8 {
		return nil, fmt.Errorf("retentions can not above %d", math.MaxUint8)
	}
	for _, def := range rss {
		retention, err := parseRetentionDef(def)
		if err != nil {
			return nil, err
		}

		retentions = append(retentions, retention)
	}

	return retentions, nil
}

type Retention struct {
	interval      int64
	reservedSecs  int64
	reservedChunk int64
	requestedSize int64
	limitedSize   int64
}

func New(interval, reserve int64) *Retention {
	chunkCapacity := int64(chunk.ChunkLen / DefaultPointBytes)
	points := int64(reserve / interval)
	reservedChunk := int64(math.Ceil(float64(points) / float64(chunkCapacity)))
	size := reservedChunk * int64(chunkLenWithHeader)
	limit := int64((DefaultLimitFactor + 1) * float64(size))
	return &Retention{interval, reserve, reservedChunk, size, limit}
}

func (r *Retention) RequestedBytes() int64 {
	return r.requestedSize
}

func (r *Retention) LimitedBytes() int64 {
	return r.limitedSize
}

func (r *Retention) ReservedChunks() int64 {
	return r.reservedChunk
}

func (r *Retention) Interval() time.Duration {
	return time.Duration(r.interval * int64(time.Second))
}

func (r *Retention) MaxReserved() time.Duration {
	return time.Duration(r.reservedSecs * int64(time.Second))
}

func GetInterval(start, stop model.Time) time.Duration {
	d := stop.Sub(start)
	for _, r := range DefaultRetentionList {
		if r.MaxReserved() > d {
			return r.Interval()
		}
	}

	return DefaultRetentionList[len(DefaultRetentionList)-1].Interval()
}

func nextInterval(idx int) time.Duration {

	if len(DefaultRetentionList)-1 < idx+1 {
		return DefaultRetentionList[idx].Interval() * 10
	} else {
		return DefaultRetentionList[idx+1].Interval()
	}

}

func init() {

	DefaultRetentionList, _ = ParseRetentionDefs(TestRetentions)
}

package bench_test

import (
	"fmt"
	"math/rand"
	"os"
	"reflect"
	"runtime"
	"strconv"
	"strings"
	"sync"
	"testing"
	"time"
	"unsafe"

	"github.com/google/uuid"
	"golang.org/x/exp/constraints"

	"github.com/EinfachAndy/hashmaps"
	cornelk "github.com/cornelk/hashmap"
	"github.com/dolthub/swiss"
	g "github.com/zyedidia/generic"
	gmap "github.com/zyedidia/generic/hashmap"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

func getRanges() []int {
	r := os.Getenv("RANGES")
	if r == "" {
		r = "50000 100000 200000 400000 600000 800000 1000000 1200000 1400000 1600000 1800000 2000000 2200000 2400000 2600000 2800000 3000000"
	}
	rangesStr := strings.Split(r, " ")
	rangesInt := make([]int, len(rangesStr))

	for i := range rangesStr {
		var err error
		rangesInt[i], err = strconv.Atoi(rangesStr[i])
		if err != nil {
			panic(err)
		}
	}
	return rangesInt
}

func handleElem[K comparable, V any](key K, val V) bool {
	return false
}
func handleElem2[K comparable, V any](key K, val V) {}

func handleElem3(key any, value any) bool {
	return false
}

func getMapNames() []string {
	m := os.Getenv("MAPS")
	if m == "" {
		m = "std robin robinLowLoad unordered swiss generic flat"
	}
	return strings.Split(m, " ")
}

func createMap[K hashmaps.Ordered, V any](n int, mapName string) hashmaps.IHashMap[K, V] {
	switch mapName {
	case "std":
		m := make(map[K]V, n)
		return hashmaps.IHashMap[K, V]{
			Put: func(k K, v V) bool {
				m[k] = v
				return false
			},
			Get: func(k K) (V, bool) {
				v, ok := m[k]
				return v, ok
			},
			Remove: func(k K) bool {
				delete(m, k)
				return true
			},
			Each: func(callback func(key K, val V) bool) {
				for k, v := range m {
					if callback(k, v) {
						return
					}
				}
			},
			Load: func() float32 {
				return -1.0 //unknown
			},
		}
	case "robin":
		m := hashmaps.NewRobinHood[K, V]()
		m.Reserve(uintptr(n))
		return hashmaps.IHashMap[K, V]{
			Get:     m.Get,
			Reserve: m.Reserve,
			Put:     m.Put,
			Remove:  m.Remove,
			Clear:   m.Clear,
			Size:    m.Size,
			Each:    m.Each,
			Load:    m.Load,
		}
	case "unordered":
		m := hashmaps.NewUnordered[K, V]()
		m.Reserve(uintptr(n))
		return hashmaps.IHashMap[K, V]{
			Get:     m.Get,
			Reserve: m.Reserve,
			Put:     m.Put,
			Remove:  m.Remove,
			Clear:   m.Clear,
			Size:    m.Size,
			Each:    m.Each,
			Load:    m.Load,
		}
	case "robinLowLoad":
		m := hashmaps.NewRobinHood[K, V]()
		m.Reserve(uintptr(n))
		m.MaxLoad(0.5)
		return hashmaps.IHashMap[K, V]{
			Get:     m.Get,
			Reserve: m.Reserve,
			Put:     m.Put,
			Remove:  m.Remove,
			Clear:   m.Clear,
			Size:    m.Size,
			Each:    m.Each,
			Load:    m.Load,
		}
	case "flat":
		m := hashmaps.NewFlat[K, V]()
		m.Reserve(uintptr(n))
		return hashmaps.IHashMap[K, V]{
			Get:     m.Get,
			Reserve: m.Reserve,
			Put:     m.Put,
			Remove:  m.Remove,
			Clear:   m.Clear,
			Size:    m.Size,
			Each:    m.Each,
			Load:    m.Load,
		}
	case "swiss":
		m := swiss.NewMap[K, V](uint32(n))
		return hashmaps.IHashMap[K, V]{
			Get: m.Get,
			Put: func(k K, v V) bool {
				m.Put(k, v)
				return true
			},
			Remove: m.Delete,
			Size:   m.Count,
			Each:   m.Iter,
			Load: func() float32 {
				return -1.0 //unknown
			},
		}
	case "generic":
		var (
			key K
			m   *gmap.Map[K, V]
		)
		kind := reflect.ValueOf(&key).Elem().Type().Kind()
		switch kind {
		case reflect.Uint32:
			var x = g.HashUint32
			m = gmap.New[K, V](uint64(n), g.Equals[K], *(*func(K) uint64)(unsafe.Pointer(&x)))
		case reflect.Uint64:
			var x = g.HashUint64
			m = gmap.New[K, V](uint64(n), g.Equals[K], *(*func(K) uint64)(unsafe.Pointer(&x)))
		case reflect.String:
			var x = g.HashString
			m = gmap.New[K, V](uint64(n), g.Equals[K], *(*func(K) uint64)(unsafe.Pointer(&x)))
		default:
			panic("type not supported")
		}
		return hashmaps.IHashMap[K, V]{
			Get: m.Get,
			Put: func(k K, v V) bool {
				m.Put(k, v)
				return true
			},
			Remove: func(k K) bool {
				m.Remove(k)
				return true
			},
			Size: m.Size,
			Each: func(callback func(key K, val V) bool) {
				m.Each(handleElem2[K, V])
			},
			Load: func() float32 {
				return -1.0 //unknown
			},
		}
	case "cornelk":
		// very slow
		m := cornelk.New[K, V]()
		m.Grow(uintptr(n))
		return hashmaps.IHashMap[K, V]{
			Get:    m.Get,
			Put:    m.Insert,
			Remove: m.Del,
			Size:   m.Len,
			Each:   m.Range,
			Load: func() float32 {
				return float32(m.FillRate()) / 100.0
			},
		}
	case "sync":
		m := &sync.Map{}
		return hashmaps.IHashMap[K, V]{
			Get: func(k K) (V, bool) {
				v, ok := m.Load(k)
				return v.(V), ok
			},
			Put: func(k K, v V) bool {
				m.Store(k, v)
				return true
			},
			Remove: func(k K) bool {
				_, ok := m.LoadAndDelete(k)
				return ok
			},
			Each: func(callback func(key K, val V) bool) {
				m.Range(handleElem3)
			},
			Load: func() float32 {
				return -1.0 //unknown
			},
		}

	default:
		panic(fmt.Sprintln("unknown map:", mapName))
	}
}

func genRandIntArray[V constraints.Integer](n int) []V {
	values := make(map[V]bool, n)
	values[0] = true
	arr := make([]V, n)
	for i := 0; i < n; {
		x := V(rand.Uint64())
		_, found := values[x]
		if !found {
			values[x] = true
			arr[i] = x
			i++
		}
	}
	return arr
}

func genShuffledIntArray[V constraints.Integer](n int) []V {
	arr := make([]V, n)
	for i := range arr {
		arr[i] = V(i + 1)
	}
	rand.Seed(time.Now().UnixNano())
	rand.Shuffle(len(arr), func(i, j int) { arr[i], arr[j] = arr[j], arr[i] })
	return arr
}

func genDifferentRandIntArray[V constraints.Integer](in []V) []V {
	out := make([]V, len(in))
	values := make(map[V]bool, len(in))
	for _, x := range in {
		values[x] = true
	}
	values[0] = true

	for j := 0; j < len(out); {
		y := V(rand.Uint64())
		_, found := values[y]
		if !found {
			out[j] = y
			j++
		}
	}

	return out
}

func genUUIDArray(n int) []string {
	arr := make([]string, n)
	for i := range arr {
		arr[i] = uuid.NewString()
	}
	return arr
}

func report(b *testing.B, n int, load float32) {
	b.ReportAllocs()
	b.ReportMetric(float64(n), "N-runs")
	b.ReportMetric(float64(load), "Load")
	var mem runtime.MemStats
	runtime.ReadMemStats(&mem)
	b.ReportMetric(float64(mem.Alloc), "Bytes")
}

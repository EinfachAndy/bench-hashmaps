package bench_test

import (
	"fmt"
	"math/rand"
	"testing"
)

// Before the test, a vector with the random uuids is generated.
// Then for each value in the vector, the key-value pair (k, 1) is inserted into the hash map.
func BenchmarkUUIDInserts(b *testing.B) {
	for _, r := range getRanges() {
		arr := genUUIDArray(r)
		for _, mapName := range getMapNames() {
			b.Run(fmt.Sprintf("%s-%d", mapName, r), func(b *testing.B) {
				load := float32(-1.0)
				for i := 0; i < b.N; i++ {
					b.StopTimer()

					m := createMap[string, uint64](0, mapName)

					b.StartTimer()
					for j := range arr {
						m.Put(arr[j], 1)
					}
					b.StopTimer()

					load = m.Load()
				}
				report(b, r, load)
			})
		}
	}
}

func BenchmarkUUIDInsertsWithReserve(b *testing.B) {
	for _, r := range getRanges() {
		arr := genUUIDArray(r)
		for _, mapName := range getMapNames() {
			b.Run(fmt.Sprintf("%s-%d", mapName, r), func(b *testing.B) {
				load := float32(-1.0)
				for i := 0; i < b.N; i++ {
					b.StopTimer()

					m := createMap[string, uint64](r, mapName)

					b.StartTimer()
					for j := range arr {
						m.Put(arr[j], 1)
					}
					b.StopTimer()

					load = m.Load()
				}
				report(b, r, load)
			})
		}
	}
}

// Before the test, n elements in the same way as in the uuid insert test are added.
// Each key is deleted one by one in a different and random order than the one they were inserted.
func BenchmarkUUIDRandomFullDeletes(b *testing.B) {
	for _, r := range getRanges() {
		arr := genUUIDArray(r)
		for _, mapName := range getMapNames() {
			b.Run(fmt.Sprintf("%s-%d", mapName, r), func(b *testing.B) {
				load := float32(-1.0)
				for i := 0; i < b.N; i++ {
					b.StopTimer()

					m := createMap[string, uint64](r, mapName)
					for j := range arr {
						m.Put(arr[j], 1)
					}
					rand.Shuffle(len(arr), func(i, j int) { arr[i], arr[j] = arr[j], arr[i] })

					b.StartTimer()
					for j := range arr {
						m.Remove(arr[j])
					}
					b.StopTimer()

					load = m.Load()
				}
				report(b, r, load)
			})
		}
	}
}

// Before the test, n elements are inserted in the same way as in the uuid insert test.
// Read each key-value pair is look up in a different and random order than the one they were inserted.
func BenchmarkUUIDRandomReads(b *testing.B) {
	for _, r := range getRanges() {
		arr := genUUIDArray(r)
		for _, mapName := range getMapNames() {
			b.Run(fmt.Sprintf("%s-%d", mapName, r), func(b *testing.B) {
				load := float32(-1.0)
				for i := 0; i < b.N; i++ {
					b.StopTimer()

					m := createMap[string, uint64](0, mapName)
					for j := range arr {
						m.Put(arr[j], 1)
					}
					rand.Shuffle(len(arr), func(i, j int) { arr[i], arr[j] = arr[j], arr[i] })

					b.StartTimer()
					for j := range arr {
						_, found := m.Get(arr[j])
						if !found {
							b.Fatal("inserted key not found")
						}
					}
					b.StopTimer()

					load = m.Load()
				}
				report(b, r, load)
			})
		}
	}
}

// Before the test, n elements are inserted in the same way as in the uuid insert test.
// Then a another vector of n random elements different from the inserted elements is generated
// which is tried to search in the hash map.
func BenchmarkUUIDReadsMisses(b *testing.B) {
	for _, r := range getRanges() {
		arr := genUUIDArray(r)
		other := genUUIDArray(r)
		for _, mapName := range getMapNames() {
			b.Run(fmt.Sprintf("%s-%d", mapName, r), func(b *testing.B) {
				load := float32(-1.0)
				for i := 0; i < b.N; i++ {
					b.StopTimer()

					m := createMap[string, uint64](0, mapName)

					for j := range arr {
						m.Put(arr[j], 1)
					}

					b.StartTimer()
					for j := range arr {
						_, found := m.Get(other[j])
						if found {
							b.Fatal("missed key was found")
						}
					}
					b.StopTimer()

					load = m.Load()
				}
				report(b, r, load)
			})
		}
	}
}

// Before the test, we insert n elements in the same way as in the uuid insert test
// before deleting half of these values randomly. We then try to read all the original values
// in a different order which will lead to 50% hits and 50% misses.
func BenchmarkUUIDRandomFullReadsAfterDeletingHalf(b *testing.B) {
	for _, r := range getRanges() {
		arr := genUUIDArray(r)
		for _, mapName := range getMapNames() {
			b.Run(fmt.Sprintf("%s-%d", mapName, r), func(b *testing.B) {
				load := float32(-1.0)
				for i := 0; i < b.N; i++ {
					b.StopTimer()

					m := createMap[string, uint64](0, mapName)
					for j := range arr {
						m.Put(arr[j], 1)
					}
					rand.Shuffle(len(arr), func(i, j int) { arr[i], arr[j] = arr[j], arr[i] })
					numRemoved := len(arr) / 2
					for j := 0; j < numRemoved; j++ {
						m.Remove(arr[j])
					}
					rand.Shuffle(len(arr), func(i, j int) { arr[i], arr[j] = arr[j], arr[i] })

					b.StartTimer()
					ac := 0
					for j := range arr {
						_, found := m.Get(arr[j])
						x := 0
						if !found {
							x = 1
						}
						ac = ac + x
					}
					b.StopTimer()
					if ac != numRemoved {
						b.Fatal("unexpected lookup accumulation:", ac)
					}

					load = m.Load()
				}
				report(b, r, load)
			})
		}
	}
}

// Before the test, n elements are inserted in the same way as in the uuid insert test.
// Then use the hash map iterators to read all the key-value pairs.
func BenchmarkUUIDRandomFullIteration(b *testing.B) {
	for _, r := range getRanges() {
		arr := genUUIDArray(r)
		for _, mapName := range getMapNames() {
			b.Run(fmt.Sprintf("%s-%d", mapName, r), func(b *testing.B) {
				load := float32(-1.0)
				for i := 0; i < b.N; i++ {
					b.StopTimer()

					m := createMap[string, uint64](0, mapName)

					for j := range arr {
						m.Put(arr[j], 1)
					}

					b.StartTimer()
					m.Each(handleElem[string, uint64])
					b.StopTimer()

					load = m.Load()
				}
				report(b, r, load)
			})
		}
	}
}

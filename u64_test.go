// package bench is inspired by https://tessil.github.io/2016/08/29/benchmark-hopscotch-map.html
package bench_test

import (
	"fmt"
	"math/rand"
	"runtime"
	"testing"
)

// Before the test, a vector with the values [0, N) is generated and shuffled.
// Then for each value in the vector, the key-value pair (k, 1) is inserted into the hash map.
func BenchmarkU64RandomShuffleInserts(b *testing.B) {
	for _, r := range getRanges() {
		arr := genShuffledIntArray[uint64](r)
		for _, mapName := range getMapNames() {
			b.Run(fmt.Sprintf("%s-%d", mapName, r), func(b *testing.B) {
				load := float32(-1.0)
				for i := 0; i < b.N; i++ {
					b.StopTimer()

					m := createMap[uint64, uint64](0, mapName)

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

// Before the test, a vector with random values in range [0, 2^64-1) is generated.
// Then for each value in the vector, the key-value pair (k, 1) is inserted into the hash map.
func BenchmarkU64RandomFullInserts(b *testing.B) {
	for _, r := range getRanges() {
		arr := genRandIntArray[uint64](r)
		for _, mapName := range getMapNames() {
			b.Run(fmt.Sprintf("%s-%d", mapName, r), func(b *testing.B) {
				load := float32(-1.0)
				for i := 0; i < b.N; i++ {
					b.StopTimer()

					m := createMap[uint64, uint64](0, mapName)

					b.StartTimer()
					for j := range arr {
						m.Put(arr[j], 1)
					}
					b.StopTimer()
					runtime.GC() // more accurate memory tracking

					load = m.Load()
				}
				report(b, r, load)
			})
		}
	}
}

// Same as the random full inserts test but the reserve method of the hash map is called beforehand
// to avoid any rehash during the insertion. It provides a fair comparison even if the growth factor
// of each hash map is different.
func BenchmarkU64RandomFullWithReserveInserts(b *testing.B) {
	for _, r := range getRanges() {
		arr := genRandIntArray[uint64](r)
		for _, mapName := range getMapNames() {
			b.Run(fmt.Sprintf("%s-%d", mapName, r), func(b *testing.B) {
				load := float32(-1.0)
				for i := 0; i < b.N; i++ {
					b.StopTimer()

					m := createMap[uint64, uint64](r, mapName)

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

// Before the test, n elements in the same way as in the random full insert test are added.
// Each key is deleted one by one in a different and random order than the one they were inserted.
func BenchmarkU64RandomFullDeletes(b *testing.B) {
	for _, r := range getRanges() {
		arr := genRandIntArray[uint64](r)
		for _, mapName := range getMapNames() {
			b.Run(fmt.Sprintf("%s-%d", mapName, r), func(b *testing.B) {
				load := float32(-1.0)
				for i := 0; i < b.N; i++ {
					b.StopTimer()

					m := createMap[uint64, uint64](r, mapName)
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

// Before the test, n elements are inserted in the same way as in the random shuffle inserts test.
// Read each key-value pair is look up in a different and random order than the one they were inserted.
func BenchmarkU64RandomShuffleReads(b *testing.B) {
	for _, r := range getRanges() {
		arr := genShuffledIntArray[uint64](r)
		for _, mapName := range getMapNames() {
			b.Run(fmt.Sprintf("%s-%d", mapName, r), func(b *testing.B) {
				load := float32(-1.0)
				for i := 0; i < b.N; i++ {
					b.StopTimer()

					m := createMap[uint64, uint64](0, mapName)
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

// Before the test, n elements are inserted in the same way as in the random full inserts test.
// Read each key-value pair is look up in a different and random order than the one they were inserted.
func BenchmarkU64FullReads(b *testing.B) {
	for _, r := range getRanges() {
		arr := genRandIntArray[uint64](r)
		for _, mapName := range getMapNames() {
			b.Run(fmt.Sprintf("%s-%d", mapName, r), func(b *testing.B) {
				load := float32(-1.0)
				for i := 0; i < b.N; i++ {
					b.StopTimer()

					m := createMap[uint64, uint64](0, mapName)

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

// Before the test, n elements are inserted in the same way as in the random full inserts test.
// Then a another vector of n random elements different from the inserted elements is generated
// which is tried to search in the hash map.
func BenchmarkU64FullReadsMisses(b *testing.B) {
	for _, r := range getRanges() {
		arr := genRandIntArray[uint64](r)
		other := genDifferentRandIntArray[uint64](arr)
		for _, mapName := range getMapNames() {
			b.Run(fmt.Sprintf("%s-%d", mapName, r), func(b *testing.B) {
				load := float32(-1.0)
				for i := 0; i < b.N; i++ {
					b.StopTimer()

					m := createMap[uint64, uint64](0, mapName)

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

// Before the test, we insert n elements in the same way as in the random full inserts test
// before deleting half of these values randomly. We then try to read all the original values
// in a different order which will lead to 50% hits and 50% misses.
func BenchmarkU64RandomFullReadsAfterDeletingHalf(b *testing.B) {
	for _, r := range getRanges() {
		arr := genRandIntArray[uint64](r)
		for _, mapName := range getMapNames() {
			b.Run(fmt.Sprintf("%s-%d", mapName, r), func(b *testing.B) {
				load := float32(-1.0)
				for i := 0; i < b.N; i++ {
					b.StopTimer()

					m := createMap[uint64, uint64](0, mapName)
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

// Before the test, n elements are inserted in the same way as in the random full inserts test.
// Then use the hash map iterators to read all the key-value pairs.
func BenchmarkU64RandomFullIteration(b *testing.B) {
	for _, r := range getRanges() {
		arr := genRandIntArray[uint64](r)
		for _, mapName := range getMapNames() {
			b.Run(fmt.Sprintf("%s-%d", mapName, r), func(b *testing.B) {
				load := float32(-1.0)
				for i := 0; i < b.N; i++ {
					b.StopTimer()

					m := createMap[uint64, uint64](0, mapName)

					for j := range arr {
						m.Put(arr[j], 1)
					}

					b.StartTimer()
					m.Each(handleElem[uint64, uint64])
					b.StopTimer()

					load = m.Load()
				}
				report(b, r, load)
			})
		}
	}
}

// Before the test, n/2 elements are inserted in the same way as the random full inserts test.
// Then the vector is shuffled and processed (50% reads, 25% inserts, 25% deletes).
// This tests combines all operations with a successful vs unsuccessful rate that is about 50/50.
func BenchmarkU64_50Reads_25Inserts_25Deletes(b *testing.B) {
	for _, r := range getRanges() {
		arr := genRandIntArray[uint64](r)
		for _, mapName := range getMapNames() {
			b.Run(fmt.Sprintf("%s-%d", mapName, r), func(b *testing.B) {
				load := float32(-1.0)
				for i := 0; i < b.N; i++ {
					b.StopTimer()

					m := createMap[uint64, uint64](0, mapName)

					for j := 0; j < len(arr)/2; j++ {
						m.Put(arr[j], arr[j])
					}
					rand.Shuffle(len(arr), func(i, j int) { arr[i], arr[j] = arr[j], arr[i] })

					b.StartTimer()
					for _, key := range arr {
						switch rand.Intn(4) {
						case 0:
							fallthrough
						case 1:
							m.Get(key)
						case 2:
							m.Put(key, key)
						case 3:
							m.Remove(key)
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

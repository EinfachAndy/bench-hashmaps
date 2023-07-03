// package bench is inspired by https://tessil.github.io/2016/08/29/benchmark-hopscotch-map.html
package bench_test

import (
	"fmt"
	"math/rand"
	"runtime"
	"testing"
)

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

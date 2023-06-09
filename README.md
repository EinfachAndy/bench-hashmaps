# Golang Hash Map Benchmark

[![MIT License](https://img.shields.io/badge/license-MIT-blue.svg)](https://github.com/EinfachAndy/bench-hashmaps/blob/main/LICENSE)

# How to run:

```bash
git clone git@github.com:EinfachAndy/bench-hashmaps.git
cd bench-hashmaps
make run-bench
make charts
```

## Run custom benchmark

following environment variables can be configure:

- `RANGES` list of integers (n)
- `MAPS` list of map names

```bash
MAPS="swiss std" RANGES="50000 100000 200000 400000" make run-bench
make charts
```

### Supported hash maps

| Name          | Module                |
|---------------|-----------------------|
| std           | golang map |
| robin         | https://pkg.go.dev/github.com/EinfachAndy/hashmaps#RobinHood (lf 0.8) |
| robinLowLoad  | https://pkg.go.dev/github.com/EinfachAndy/hashmaps#RobinHood (lf 0.5) |
| unordered     | https://pkg.go.dev/github.com/EinfachAndy/hashmaps#Unordered |
| swiss         | https://pkg.go.dev/github.com/dolthub/swiss#Map |
| generic       | https://pkg.go.dev/github.com/zyedidia/generic/hashmap#Map |
| cornelk       | https://pkg.go.dev/github.com/cornelk/hashmap#Map |
| sync          | https://pkg.go.dev/sync#Map |

# Contributing

If you would like to add a new benchmark or hash map, feel free to contribute.

### Note:
This benchmark is inspired from [Benchmark of major hash maps implementations](https://tessil.github.io/2016/08/29/benchmark-hopscotch-map.html).

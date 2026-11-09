Median of 3 rounds (after 1 unrecorded warm-up round), excluding round r1 at RATE=60, DURATION=2m. Host: 12 threads, 3703 MiB RAM, 6.6.87.2-microsoft-standard-WSL2. API limited to 2 CPUs / 512m.

| Metric | gin-gorm | nethttp-gorm | nethttp-sqlc |
| --- | ---: | ---: | ---: |
| Browse p50 (ms) | 4.3 | 4.4 | 2.9 |
| Browse p95 (ms) | 36.4 | 34.4 | 29.1 |
| Browse p99 (ms) | 47.0 | 46.4 | 41.3 |
| Submit p95 (ms) | 50.5 | 54.2 | 47.2 |
| Login p95 (ms) | 268 | 263 | 261 |
| Throughput (req/s) | 220 | 220 | 220 |
| Failed requests (%) | 0.00 | 0.00 | 0.00 |
| API CPU avg (%) | 108 | 109 | 100 |
| API memory avg (MiB) | 23 | 23 | 21 |
| API memory max (MiB) | 26 | 27 | 23 |
| Binary size (MB) | 45.5 | 34.7 | 33.1 |
| Image size (MB) | 152.4 | 137.9 | 135.6 |

# Branching factors

Implemented a general solution at e014880b4712d4617f40c1c261105c5261aa276b.

```
Running tool: /usr/local/go/bin/go test -benchmem -run=^$ github.com/vsekhar/mmr -bench ^(BenchmarkHeight)$ -v

goos: darwin
goarch: amd64
pkg: github.com/vsekhar/mmr
BenchmarkHeight/branching-2-4         	76398361	        17.6 ns/op	       0 B/op	       0 allocs/op
BenchmarkHeight/branching-3-4         	 1546952	       795 ns/op	       0 B/op	       0 allocs/op
BenchmarkHeight/branching-4-4         	   32583	     36309 ns/op	       0 B/op	       0 allocs/op
BenchmarkHeight/branching-7-4         	  615481	      2153 ns/op	       0 B/op	       0 allocs/op
BenchmarkHeight/branching-8-4         	   53630	     21217 ns/op	       0 B/op	       0 allocs/op
BenchmarkHeight/branching-12-4        	  364978	      2984 ns/op	       0 B/op	       0 allocs/op
BenchmarkHeight/branching-16-4        	   86502	     14406 ns/op	       0 B/op	       0 allocs/op
BenchmarkHeight/branching-32-4        	   61646	     20869 ns/op	       0 B/op	       0 allocs/op
BenchmarkHeight/branching-100-4       	   68580	     15083 ns/op	       0 B/op	       0 allocs/op
BenchmarkHeight/branching-1000-4      	   33778	     40108 ns/op	       0 B/op	       0 allocs/op
PASS
ok  	github.com/vsekhar/mmr	15.670s
Success: Benchmarks passed.
```

Conclusion: binary trees are way fast.

Removed general solution.

# Notes

## TODOs

--

## Open questions

### API

Decisions needed to turn MMR math into a data structure:

* Contentful vs. content-free tree
* How to hash payload/node
* Require client to store salt?

General solution?

### Proofs

Proof sequence is fixed for any `{pos, n}` pair. So the proof for `{pos, n}` need only be a list of values to be combined as the application needs.

So inclusion proof derivation need only return a list of positions to include in the proof, which the application can extract and send off. Validating an inclusion proof need only accept the `{pos, n}` pair and the list of values at the proof positions.

```go
func Inclusion(pos, n int) []int // returns list of positions to include in proof
```

Generalized validator? Something like the sort.Interface with index-based methods like `Swap(i, j int)`. Issue is managing the stack. Maybe pass a stack machine? Don't even need to specify which values, they are already in order of use.

```go
type Combiner interface {
    Combine(pos int, accumulatorFirst bool)
}

func Validate(pos, n int, v Combiner) // TODO: where is the result?
```

What the library code will do is call those functions in the right order based on the input. Library code doesn't even need the input, and the calls to methods of Validator don't need positions (just use the next value in the input).

Is there a way for the application Validator to not have to manage a stack? Need a way to tell it what to do with the proof data it got.

Thoughts: the stack machine is evaluating a nested expression. Every `{pos, n}` pair produces a machine that requires a specific input sequence of values.

BUT: evaluation of a proof has all proof values, could use random access.

## Appendix: Old stuff

### Branching factors other than 2

Implemented a general solution at `e014880b4712d4617f40c1c261105c5261aa276b`.

```shell
Running tool: /usr/local/go/bin/go test -benchmem -run=^$ github.com/vsekhar/mmr -bench ^(BenchmarkHeight)$ -v

goos: darwin
goarch: amd64
pkg: github.com/vsekhar/mmr
BenchmarkHeight/branching-2-4          76398361       17.6 ns/op        0 B/op        0 allocs/op
BenchmarkHeight/branching-3-4           1546952        795 ns/op        0 B/op        0 allocs/op
BenchmarkHeight/branching-4-4             32583      36309 ns/op        0 B/op        0 allocs/op
BenchmarkHeight/branching-7-4            615481       2153 ns/op        0 B/op        0 allocs/op
BenchmarkHeight/branching-8-4             53630      21217 ns/op        0 B/op        0 allocs/op
BenchmarkHeight/branching-12-4           364978       2984 ns/op        0 B/op        0 allocs/op
BenchmarkHeight/branching-16-4            86502      14406 ns/op        0 B/op        0 allocs/op
BenchmarkHeight/branching-32-4            61646      20869 ns/op        0 B/op        0 allocs/op
BenchmarkHeight/branching-100-4           68580      15083 ns/op        0 B/op        0 allocs/op
BenchmarkHeight/branching-1000-4          33778      40108 ns/op        0 B/op        0 allocs/op
PASS
ok   github.com/vsekhar/mmr 15.670s
Success: Benchmarks passed.
```

Conclusion: binary trees are way fast.

Removed general solution.

### Contentful vs. content-free trees

Contentful trees have user content in non-leaf nodes, whereas content-free trees have (or refer to) user content only in leaf nodes.

For a binary tree, half of nodes are leaf nodes, other half are non-leaf. Content-free trees need only one hash in each non-leaf node (that of its priors/children). Contentful trees need two hashes (that of the user content and that of its priors/children) in order to recreate the hash used in subsequent nodes.

But, if we use the priors/children as the "salt" and require clients to store the salt, then we don't need to store the payloads, we just store `hash(hash([priors...]), payload)` and return `hash([priors...])` to the client.

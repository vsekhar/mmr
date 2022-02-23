package mmr

type Combiner interface {
	// Combine combines the values at positions i and j (in that order) and places
	// the combined result at position i.
	Combine(i, j int)
}

// There are three kinds of proofs
//
//   1) Tree proof: a proof that a value at position p1 is included in the subtree anchored
//      at some later value p2
//                   - p1 is usually a peak of some digest
//                   - Proof consists of pre- and post-hashing different values starting with
//                     p1 end producing p2.
//   2) Peak proof: a proof that a peak at position p2 is included in the digest d1 of the MMR
//      with size n.
//                   - Proof consists of an ordered list of values with a designated spot for
//                     p2 that produces d1.
//   3) Digest proof: a proof that digest of the MMR of size n, d1, is included in the digest of
//      the MMR of size m>n, d2.
//        - Proof consists of an ordered list of values that produces d1, an inclusion proof
//          of the last peak of d1 in an appropriate peak p2 of d2, and an inclusion proof of p2 in d2.
//
// Proofs fit together as follows:
//
//   A) Inclusion proof of p1 in d1 = 1 + 2
//   B) Consistency proof of d1 in d2 = 3 + ???

//

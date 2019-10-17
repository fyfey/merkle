### Merkle Tree

The beginnings of a merkle tree file transfer system a-la bit-torrent.

The idea is to chunk a file up and store the merkle tree. Then send chunks and use the tree to verify the chunks as they arrive the other end.

```go
go run main.go arrival_in_nara.txt 16
```
As you can see, out of 661 chunks, it only takes 9 verifications from the merkle proof to prove that a certain chunk is valid in the tree.
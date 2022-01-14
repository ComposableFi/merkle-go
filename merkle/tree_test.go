package merkle_test

import (
	"testing"

	"github.com/ComposableFi/merkle-go/merkle"
	"github.com/stretchr/testify/require"
)

func TestNewMerkleTree(t *testing.T) {
	merkle.NewTree(Sha256Hasher{})
}

func TestFromLeaves(t *testing.T) {
	aHash, err := Sha256Hasher{}.Hash([]byte("a"))
	require.NoError(t, err)
	bHash, err := Sha256Hasher{}.Hash([]byte("b"))
	require.NoError(t, err)
	cHash, err := Sha256Hasher{}.Hash([]byte("c"))
	require.NoError(t, err)

	leaves := []merkle.Hash{aHash, bHash, cHash}
	mtree := merkle.NewTree(Sha256Hasher{})
	mtree, err = mtree.FromLeaves(leaves)
	require.NoError(t, err)
	require.Equal(t, []merkle.Hash{}, mtree.UncommittedLeaves)
}

func TestRoot(t *testing.T) {
	aHash, err := Sha256Hasher{}.Hash([]byte("a"))
	require.NoError(t, err)
	bHash, err := Sha256Hasher{}.Hash([]byte("b"))
	require.NoError(t, err)
	cHash, err := Sha256Hasher{}.Hash([]byte("c"))
	require.NoError(t, err)

	leaves := []merkle.Hash{aHash, bHash, cHash}
	mtree := merkle.NewTree(Sha256Hasher{})
	mtree, err = mtree.FromLeaves(leaves)
	require.NoError(t, err)

	indicesToProve := []uint32{0, 1}
	leavesToProve := leaves[0:2]
	proof := mtree.Proof(indicesToProve)
	root := mtree.GetRoot()

	leafTuples := merkle.MapIndiceAndLeaves(indicesToProve, leavesToProve)
	verified, err := proof.Verify(root, leafTuples, len(leaves))
	require.NoError(t, err)
	require.True(t, verified)
}

func TestProof(t *testing.T) {
	values := []string{"a", "b", "c", "d", "e", "f"}
	var leaves []merkle.Hash
	for i := 0; i < len(values); i++ {
		h, _ := Sha256Hasher{}.Hash([]byte(values[i]))
		leaves = append(leaves, h)
	}
	mtree := merkle.NewTree(Sha256Hasher{})
	mtree, err := mtree.FromLeaves(leaves)
	require.NoError(t, err)

	indicesToProve := []uint32{3, 4}
	leavesToProve := leaves[3:5]
	proof := mtree.Proof(indicesToProve)
	root := mtree.GetRoot()

	leafTuples := merkle.MapIndiceAndLeaves(indicesToProve, leavesToProve)
	verified, err := proof.Verify(root, leafTuples, len(leaves))
	require.NoError(t, err)
	require.True(t, verified)
}

func TestCorrectTreeRoot(t *testing.T) {
	testData := setup()

	merkleTree, err := merkle.NewTree(Sha256Hasher{}).FromLeaves(testData.leafHashes)
	require.NoError(t, err)
	rootHex := merkleTree.GetRootHex()
	require.Equal(t, testData.expectedRootHex, rootHex)
}

func TestCorrectTreeDepth(t *testing.T) {
	testData := setup()

	merkleTree, err := merkle.NewTree(Sha256Hasher{}).FromLeaves(testData.leafHashes)
	require.NoError(t, err)
	depth := merkleTree.GetDepth()
	require.Equal(t, 3, depth)
}

func TestCorrectProofRoot(t *testing.T) {
	testData := setup()
	indicesToProve := []uint32{3, 4}
	expectedProofHashes := []string{
		"2e7d2c03a9507ae265ecf5b5356885a53393a2029d241394997265a1a25aefc6",
		"252f10c83610ebca1a059c0bae8255eba2f95be4d1d7bcfa89d7248a82d9f111",
		"e5a01fee14e0ed5c48714f22180f25ad8365b53f9779f79dc4a3d7e93963f94a",
	}
	merkleTree, err := merkle.NewTree(Sha256Hasher{}).FromLeaves(testData.leafHashes)
	require.NoError(t, err)
	proof := merkleTree.Proof(indicesToProve)
	require.Equal(t, expectedProofHashes, proof.ProofHashesHex())
}

func TestGetCorrectRootAfterCommit(t *testing.T) {
	testData := setup()
	expectedRoot := testData.expectedRootHex
	leafHashes := testData.leafHashes

	merkleTree, err := merkle.NewTree(Sha256Hasher{}).FromLeaves([]merkle.Hash{})
	require.NoError(t, err)
	merkleTree2, err := merkle.NewTree(Sha256Hasher{}).FromLeaves(leafHashes)
	require.NoError(t, err)

	merkleTree.Append(leafHashes)

	root, err := merkleTree.UncommittedRootHex()
	require.NoError(t, err)

	require.Equal(t, expectedRoot, merkleTree2.GetRootHex())
	require.Equal(t, expectedRoot, root)

	expectedRoot = "e2a80e0e872a6c6eaed37b4c1f220e1935004805585b5f99617e48e9c8fe4034"
	leaf, err := Sha256Hasher{}.Hash([]byte("g"))
	require.NoError(t, err)
	merkleTree.Insert(leaf)

	uncommittedRoot, err := merkleTree.UncommittedRootHex()
	require.NoError(t, err)
	require.Equal(t, expectedRoot, uncommittedRoot)

	require.Equal(t, merkle.Hash{}, merkleTree.GetRoot())

	merkleTree.Commit()

	hHash, _ := Sha256Hasher{}.Hash([]byte("h"))
	kHash, _ := Sha256Hasher{}.Hash([]byte("k"))
	newLeaves := []merkle.Hash{hHash, kHash}

	merkleTree.Append(newLeaves)

	require.Equal(t, "e2a80e0e872a6c6eaed37b4c1f220e1935004805585b5f99617e48e9c8fe4034", merkleTree.GetRootHex())
	uncommittedRootHex, err := merkleTree.UncommittedRootHex()
	require.NoError(t, err)
	require.Equal(t, "09b6890b23e32e607f0e5f670ab224e36af8f6599cbe88b468f4b0f761802dd6", uncommittedRootHex)

	merkleTree.Commit()

	leaves := merkleTree.GetLeaves()
	reconstructedTree, err := merkleTree.FromLeaves(leaves)
	require.NoError(t, err)
	require.Equal(t, "09b6890b23e32e607f0e5f670ab224e36af8f6599cbe88b468f4b0f761802dd6", reconstructedTree.GetRootHex())

}

func TestChangeTheResultWenCalledTwice(t *testing.T) {
	leafValues := []string{"a", "b", "c", "d", "e", "f"}
	var leaves []merkle.Hash
	for _, v := range leafValues {
		h, _ := Sha256Hasher{}.Hash([]byte(v))
		leaves = append(leaves, h)
	}

	merkleTree := merkle.NewTree(Sha256Hasher{})

	// Appending leaves to the tree without committing
	merkleTree.Append(leaves)

	require.Equal(t, merkle.Hash{}, merkleTree.GetRoot())

	uncommittedRootHex, err := merkleTree.UncommittedRootHex()
	require.NoError(t, err)
	require.Equal(t, "1f7379539707bcaea00564168d1d4d626b09b73f8a2a365234c62d763f854da2", uncommittedRootHex)

	merkleTree.Commit()

	require.Equal(t, "1f7379539707bcaea00564168d1d4d626b09b73f8a2a365234c62d763f854da2", merkleTree.GetRootHex())

	uncommittedRootHex, err = merkleTree.UncommittedRootHex()
	require.NoError(t, err)
	require.Equal(t, "", uncommittedRootHex)

	gHash, _ := Sha256Hasher{}.Hash([]byte("g"))
	merkleTree.Insert(gHash)

	uncommittedRootHex, err = merkleTree.UncommittedRootHex()
	require.NoError(t, err)
	require.Equal(t, "e2a80e0e872a6c6eaed37b4c1f220e1935004805585b5f99617e48e9c8fe4034", uncommittedRootHex)

	merkleTree.Commit()

	require.Equal(t, "e2a80e0e872a6c6eaed37b4c1f220e1935004805585b5f99617e48e9c8fe4034", merkleTree.GetRootHex())

	hHash, _ := Sha256Hasher{}.Hash([]byte("h"))
	kHash, _ := Sha256Hasher{}.Hash([]byte("k"))
	merkleTree.Append([]merkle.Hash{hHash, kHash})

	merkleTree.Commit()
	merkleTree.Commit()

	require.Equal(t, "09b6890b23e32e607f0e5f670ab224e36af8f6599cbe88b468f4b0f761802dd6", merkleTree.GetRootHex())

	merkleTree.Rollback()
	require.Equal(t, "e2a80e0e872a6c6eaed37b4c1f220e1935004805585b5f99617e48e9c8fe4034", merkleTree.GetRootHex())

	merkleTree.Rollback()
	require.Equal(t, "1f7379539707bcaea00564168d1d4d626b09b73f8a2a365234c62d763f854da2", merkleTree.GetRootHex())
}

func TestRollbackPreviousCommit(t *testing.T) {
	leafValues := []string{"a", "b", "c", "d", "e", "f"}
	var leaves []merkle.Hash
	for _, v := range leafValues {
		h, _ := Sha256Hasher{}.Hash([]byte(v))
		leaves = append(leaves, h)
	}

	merkleTree := merkle.NewTree(Sha256Hasher{})
	merkleTree.Append(leaves)

	require.Equal(t, merkle.Hash{}, merkleTree.GetRoot())

	merkleTree.Commit()

	require.Equal(t, "1f7379539707bcaea00564168d1d4d626b09b73f8a2a365234c62d763f854da2", merkleTree.GetRootHex())

	gHash, _ := Sha256Hasher{}.Hash([]byte("g"))
	merkleTree.Insert(gHash)

	uncommittedRootHex, err := merkleTree.UncommittedRootHex()
	require.NoError(t, err)

	require.Equal(t, "e2a80e0e872a6c6eaed37b4c1f220e1935004805585b5f99617e48e9c8fe4034", uncommittedRootHex)

	merkleTree.Commit()

	require.Equal(t, "e2a80e0e872a6c6eaed37b4c1f220e1935004805585b5f99617e48e9c8fe4034", merkleTree.GetRootHex())

	hHash, _ := Sha256Hasher{}.Hash([]byte("h"))
	kHash, _ := Sha256Hasher{}.Hash([]byte("k"))
	merkleTree.Append([]merkle.Hash{hHash, kHash})

	uncommittedRootHex, err = merkleTree.UncommittedRootHex()
	require.NoError(t, err)
	require.Equal(t, "09b6890b23e32e607f0e5f670ab224e36af8f6599cbe88b468f4b0f761802dd6", uncommittedRootHex)

	require.Equal(t, "e2a80e0e872a6c6eaed37b4c1f220e1935004805585b5f99617e48e9c8fe4034", merkleTree.GetRootHex())

	merkleTree.Commit()

	require.Equal(t, "09b6890b23e32e607f0e5f670ab224e36af8f6599cbe88b468f4b0f761802dd6", merkleTree.GetRootHex())

	merkleTree.Rollback()
	require.Equal(t, "e2a80e0e872a6c6eaed37b4c1f220e1935004805585b5f99617e48e9c8fe4034", merkleTree.GetRootHex())

	merkleTree.Rollback()
	require.Equal(t, "1f7379539707bcaea00564168d1d4d626b09b73f8a2a365234c62d763f854da2", merkleTree.GetRootHex())
}

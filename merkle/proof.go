package merkle

import (
	"bytes"
	"encoding/hex"
	"math"

	"github.com/ComposableFi/merkle-go/helpers"
)

// func (p Proof) fromBytes(bytes []byte) (PartialTree, error) {
// 	return p.deserialize(bytes)
// }

// func (p Proof) deserialize(bytes []byte) (PartialTree, error) {
// 	return p.serializer.Deserialize(bytes)
// }

func (p Proof) Verify(root Hash, leafTuples []Leaf, totalLeavesCount int) bool {
	extractedRoot := p.GetRoot(leafTuples, int(totalLeavesCount))
	return bytes.Equal(extractedRoot, root)
}

func (p Proof) GetRoot(leafTuples []Leaf, totalLeavesCount int) Hash {
	treeDepth := getTreeDepth(totalLeavesCount)
	sortLeavesByIndex(leafTuples)
	var leafIndices []uint32
	for _, l := range leafTuples {
		leafIndices = append(leafIndices, l.Index)
	}
	proofIndicesLayers := proofIndeciesByLayers(leafIndices, totalLeavesCount)
	var proofLayers [][]Leaf
	for _, proofIndices := range proofIndicesLayers {
		var proofHashes []Hash
		for i := 0; i < len(proofIndices); i++ {
			proofHashes = append(proofHashes, p.proofHashes[i])
		}
		m := MapIndiceAndLeaves(proofIndices, proofHashes)
		proofLayers = append(proofLayers, m)
	}

	if len(proofLayers) > 0 {
		firstLayer := proofLayers[0]
		firstLayer = append(firstLayer, leafTuples...)
		sortLeavesByIndex(firstLayer)
		proofLayers[0] = firstLayer

	} else {
		proofLayers = append(proofLayers, leafTuples)
	}
	partialTree := NewPartialTree(p.hasher)
	PartialTree, err := partialTree.build(proofLayers, treeDepth)
	if err != nil {
		return Hash{}
	}
	return PartialTree.GetRoot()
}

func (p Proof) GetRootHex(leafTuples []Leaf, totalLeavesCount int) string {
	return hex.EncodeToString(p.GetRoot(leafTuples, totalLeavesCount))
}

func (p Proof) ProofHashes() []Hash {
	return p.proofHashes
}

func proofIndeciesByLayers(sortedLeafIndices []uint32, leavsCount int) [][]uint32 {
	depth := getTreeDepth(leavsCount)
	unevenLayers := unevenLayers(leavsCount)
	var proofIndices [][]uint32
	for layerIndex := 0; layerIndex < depth; layerIndex++ {
		siblingIndices := helpers.GetSiblingIndecies(sortedLeafIndices)
		leavesCount := unevenLayers[layerIndex]
		layerLastNodeIndex := sortedLeafIndices[len(sortedLeafIndices)-1]
		if layerLastNodeIndex == uint32(leavesCount)-1 {
			_, siblingIndices = helpers.PopFromUint32Queue(siblingIndices)
		}

		proofNodesIndices := helpers.Difference(siblingIndices, sortedLeafIndices)
		proofIndices = append(proofIndices, proofNodesIndices)
		sortedLeafIndices = helpers.GetParentIndecies(sortedLeafIndices)
	}
	return proofIndices

}

func unevenLayers(treeLeavesCount int) map[int]int {
	depth := getTreeDepth(treeLeavesCount)
	unevenLayers := make(map[int]int)
	for i := 0; i < depth; i++ {
		unevenLayer := treeLeavesCount%2 != 0
		if unevenLayer {
			unevenLayers[i] = treeLeavesCount
		}
		treeLeavesCount = int(math.Ceil(float64(treeLeavesCount) / 2))
	}
	return unevenLayers
}

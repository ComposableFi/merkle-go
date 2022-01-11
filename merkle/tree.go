package merkle

import (
	"encoding/hex"
	"errors"
	"math"
	"sort"

	"github.com/ComposableFi/merkle-go/helpers"
)

func (t Tree) FromLeaves(leaves []Hash) (Tree, error) {
	t.append(leaves)
	err := t.commit()
	if err != nil {
		return Tree{}, err
	}
	return t, nil
}

func (t *Tree) GetRoot() Hash {
	layers := t.layerTuples()
	lastLayer := layers[len(layers)-1]
	firstItem := lastLayer[0]
	return firstItem.Hash
}

func (t *Tree) GetRootHex() string {
	root := t.GetRoot()
	return hex.EncodeToString([]byte(root))
}

func (t *Tree) HelperNodes(leafIndices []uint32) []Hash {
	var helperNodes []Hash
	for _, layer := range t.HelperNodeTuples(leafIndices) {
		for _, li := range layer {
			helperNodes = append(helperNodes, li.Hash)
		}
	}
	return helperNodes
}
func (t *Tree) HelperNodeTuples(leafIndeceis []uint32) [][]Leaf {
	var helperNodes [][]Leaf
	for _, treeLayer := range t.layerTuples() {
		siblings := helpers.GetSiblingIndecies(leafIndeceis)
		helperIndices := helpers.Difference(siblings, leafIndeceis)

		var helpersLayer []Leaf
		for _, idx := range helperIndices {
			leaf, found := getLeafAtIndex(treeLayer, idx)
			if found {
				helpersLayer = append(helpersLayer, leaf)
			}
		}

		if helpersLayer != nil {
			helperNodes = append(helperNodes, helpersLayer)
		}

		leafIndeceis = helpers.GetParentIndecies(leafIndeceis)
	}
	return helperNodes
}

func (t *Tree) Proof(leafIndices []uint32) Proof {
	return NewProof(t.HelperNodes(leafIndices), t.hasher)
}

func (t *Tree) insert(leaf Hash) {
	t.UncommittedLeaves = append(t.UncommittedLeaves, leaf)
}

func (t *Tree) append(leaves []Hash) {
	t.UncommittedLeaves = append(t.UncommittedLeaves, leaves...)
}

func (t *Tree) commit() error {
	diff, err := t.uncommittedDiff()
	if err != nil {
		return err
	}
	t.history = append(t.history, diff)
	t.currentWorkingTree.mergeUnverified(diff)
	t.UncommittedLeaves = []Hash{}
	return nil
}

func (t *Tree) uncommittedRoot() (Hash, error) {
	shadowTree, err := t.uncommittedDiff()
	if err != nil {
		return Hash{}, err
	}
	return shadowTree.GetRoot(), nil
}

func (t *Tree) uncommittedRootHex() (string, error) {
	root, err := t.uncommittedRoot()
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(root), err
}

func (t *Tree) abortCommitted() {
	t.UncommittedLeaves = make([]Hash, 0)
}

func (t *Tree) depth() int {
	return len(t.layerTuples()) - 1
}

func (t *Tree) leaves() []Hash {
	layers := t.layers()
	if len(layers) > 0 {
		return []Hash{}
	}
	return layers[0]
}

func (t *Tree) leavesLen() int {
	leaves := t.leavesTuples()
	return len(leaves)
}

func (t *Tree) leavesTuples() []Leaf {
	if len(t.layerTuples()) > 0 {
		return t.layerTuples()[0]
	}
	return []Leaf{}
}

func (t *Tree) layers() [][]Hash {
	return t.currentWorkingTree.layerNodes()
}

func (t *Tree) layerTuples() [][]Leaf {
	return t.currentWorkingTree.layers
}

func (t *Tree) uncommittedDiff() (PartialTree, error) {
	if len(t.UncommittedLeaves) == 0 {
		return PartialTree{}, errors.New("leaves can not be empty!")
	}
	commitedLeavesCount := t.leavesLen()
	var shadowIndecies []uint32
	for i, _ := range t.UncommittedLeaves {
		shadowIndecies = append(shadowIndecies, uint32(commitedLeavesCount+i))
	}
	var shadowNodeTuples []Leaf
	for _, idx := range shadowIndecies {
		x := Leaf{Index: idx, Hash: t.UncommittedLeaves[idx]}
		shadowNodeTuples = append(shadowNodeTuples, x)
	}
	partialTreeTuples := t.HelperNodeTuples(shadowIndecies)
	leavesInNewTree := t.leavesLen() + len(t.UncommittedLeaves)
	uncommittedTreeDepth := getTreeDepth(leavesInNewTree)
	if len(partialTreeTuples) == 0 {
		partialTreeTuples = append(partialTreeTuples, shadowNodeTuples)
	} else {
		firstLayer := partialTreeTuples[0]
		firstLayer = append(firstLayer, shadowNodeTuples...)
		sortLeavesByIndex(firstLayer)
	}
	tree := NewPartialTree(t.hasher)
	return tree.build(partialTreeTuples, uncommittedTreeDepth)
}

func getLeafAtIndex(leavesAndHash []Leaf, index uint32) (Leaf, bool) {
	for _, l := range leavesAndHash {
		if l.Index == index {
			return l, true
		}
	}
	return Leaf{}, false
}

func getLayerAtIndex(layers [][]Leaf, index uint32) ([]Leaf, bool) {
	if len(layers) > int(index) {
		return layers[index], true
	}
	return []Leaf{}, false
}

func sortLeavesByIndex(li []Leaf) {
	sort.Slice(li, func(i, j int) bool { return li[i].Index < li[j].Index })

}
func getTreeDepth(leaves_count int) int {
	if leaves_count == 1 {
		return 1
	} else {
		return int(math.Ceil(math.Log2(float64(leaves_count))))
	}
}

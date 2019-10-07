/*
 *  Brown University, CS138, Spring 2018
 *
 *  Purpose: Implements functions that are invoked by other nodes over RPC.
 */

package tapestry

import (
	"errors"

	"golang.org/x/net/context"
)

/**
 * RPC receiver functions
 */

func (local *Node) HelloCaller(ctx context.Context, n *NodeMsg) (*NodeMsg, error) {
	return local.node.toNodeMsg(), nil
}

func (local *Node) GetNextHopCaller(ctx context.Context, id *IdMsg) (*NextHop, error) {
	idVal, err := ParseID(id.Id)
	if err != nil {
		return nil, err
	}
	hasNext, next, err := local.GetNextHop(idVal)
	rsp := &NextHop{
		HasNext: hasNext,
		Next:    next.toNodeMsg(),
	}
	return rsp, err
}

//STUDENT IMPLEMENTED
// Calls register on the local node with the correct arguments and returns a response
// containing whether the local node is a root and a reason???
func (local *Node) RegisterCaller(ctx context.Context, r *Registration) (*Ok, error) {
	// func (local *Node) Register(key string, replica RemoteNode) (isRoot bool, err error)
	isRoot, err := local.Register(r.Key, r.FromNode.toRemoteNode())

	rsp := &Ok{
		Ok:     isRoot,
		Reason: "",
	}

	return rsp, err
}

func (local *Node) FetchCaller(ctx context.Context, key *Key) (*FetchedLocations, error) {
	isRoot, values, err := local.Fetch(key.Key)

	rsp := &FetchedLocations{
		Values: remoteNodesToNodeMsgs(values),
		IsRoot: isRoot,
	}
	return rsp, err
}

func (local *Node) RemoveBadNodesCaller(ctx context.Context, nodes *Neighbors) (*Ok, error) {
	err := local.RemoveBadNodes(nodeMsgsToRemoteNodes(nodes.Neighbors))
	rsp := &Ok{
		Ok: true,
	}
	return rsp, err
}

func (local *Node) AddNodeCaller(ctx context.Context, n *NodeMsg) (*Neighbors, error) {
	neighbors, err := local.AddNode(n.toRemoteNode())

	rsp := &Neighbors{
		Neighbors: remoteNodesToNodeMsgs(neighbors),
	}
	return rsp, err
}

// STUDENT WRITTEN
// Calls AddNodeMulticast, passes the error straight on
func (local *Node) AddNodeMulticastCaller(ctx context.Context, m *MulticastRequest) (*Neighbors, error) {
	neighbors, err := local.AddNodeMulticast(m.NewNode.toRemoteNode(), int(m.Level))

	rsp := &Neighbors{
		Neighbors: remoteNodesToNodeMsgs(neighbors),
	}
	return rsp, err
}

func (local *Node) TransferCaller(ctx context.Context, td *TransferData) (*Ok, error) {
	parsedData := make(map[string][]RemoteNode)
	for key, set := range td.Data {
		parsedData[key] = nodeMsgsToRemoteNodes(set.Neighbors)
	}
	err := local.Transfer(td.From.toRemoteNode(), parsedData)

	rsp := &Ok{
		Ok: true,
	}
	return rsp, err
}

func (local *Node) AddBackpointerCaller(ctx context.Context, n *NodeMsg) (*Ok, error) {
	err := local.AddBackpointer(n.toRemoteNode())
	rsp := &Ok{
		Ok: err == nil,
	}
	return rsp, err
}

func (local *Node) RemoveBackpointerCaller(ctx context.Context, n *NodeMsg) (*Ok, error) {
	err := local.RemoveBackpointer(n.toRemoteNode())
	rsp := &Ok{
		Ok: true,
	}
	return rsp, err
}

func (local *Node) GetBackpointersCaller(ctx context.Context, br *BackpointerRequest) (*Neighbors, error) {
	n, err := local.GetBackpointers(br.From.toRemoteNode(), int(br.Level))
	rsp := &Neighbors{
		Neighbors: remoteNodesToNodeMsgs(n),
	}
	return rsp, err
}

func (local *Node) NotifyLeaveCaller(ctx context.Context, ln *LeaveNotification) (*Ok, error) {
	replacement := ln.Replacement.toRemoteNode()
	err := local.NotifyLeave(ln.From.toRemoteNode(), &replacement)
	rsp := &Ok{
		Ok: true,
	}
	return rsp, err
}

func (local *Node) BlobStoreFetchCaller(ctx context.Context, key *Key) (*DataBlob, error) {
	data, isOk := local.blobstore.Get(key.Key)
	var err error
	if !isOk {
		err = errors.New("Key not found")
	}
	return &DataBlob{
		Key:  key.Key,
		Data: data,
	}, err
}

func (local *Node) TapestryLookupCaller(ctx context.Context, key *Key) (*Neighbors, error) {
	nodes, err := local.Lookup(key.Key)
	return &Neighbors{remoteNodesToNodeMsgs(nodes)}, err
}

func (local *Node) TapestryStoreCaller(ctx context.Context, blob *DataBlob) (*Ok, error) {
	return &Ok{Ok: true}, local.Store(blob.Key, blob.Data)
}

func remoteNodesToNodeMsgs(remoteNodes []RemoteNode) []*NodeMsg {
	nodeMsgs := make([]*NodeMsg, len(remoteNodes))
	for i, thing := range remoteNodes {
		nodeMsgs[i] = thing.toNodeMsg()
	}
	return nodeMsgs
}

func nodeMsgsToRemoteNodes(nodeMsgs []*NodeMsg) []RemoteNode {
	remoteNodes := make([]RemoteNode, len(nodeMsgs))
	for i, thing := range nodeMsgs {
		remoteNodes[i] = thing.toRemoteNode()
	}
	return remoteNodes
}

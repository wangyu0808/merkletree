package merkletree

import (
	"bytes"
	"crypto/md5"
	"errors"
	"fmt"
)

//MerkleTree 树
type MerkleTree struct {
	Root *MerkleNode //根节点
}
//MerkleNode 节点
type MerkleNode struct{
	Left *MerkleNode  //左节点
	Right *MerkleNode //右节点
	Data []byte       //节点数据
}

//GetMerkleRoot 生成默克尔树根节点
func GetMerkleRoot(data [][]byte) *MerkleTree {
	//节点
	var nodes []*MerkleNode
	//奇数，添加最后一个数，变为偶数
	if len(data) % 2 !=0{
		data = append(data,data[len(data)-1])
	}
	for _,v := range data{
		//所有叶子节点
		hash := md5.Sum(v)
		node := MerkleNode{nil, nil,hash[:]}
		nodes = append(nodes, &node)
		fmt.Printf("原始数据：%v，hash:%x\n",string(v),node.Data)
	}
	nodes = GetMerkleNode(nodes,true)
	//取根节点返回
	return &MerkleTree{nodes[0]}
}
// GetMerkleNode 生成节点
func GetMerkleNode(merkleNodes []*MerkleNode,isLeaf bool)  (returnMerkleNodes []*MerkleNode) {
	//如果不是叶子，且只有一个节点传入，这是根节点，返回
	if !isLeaf && len(merkleNodes) == 1 {
		return merkleNodes
	}
	var nodeLen = len(merkleNodes)
	//判断节点奇偶，奇数复制最后一份
	var isHaveCopy = false
	if nodeLen % 2 != 0 {
		copyNode := *merkleNodes[nodeLen - 1]
		merkleNodes = append(merkleNodes, &copyNode)
		nodeLen = nodeLen + 1
		isHaveCopy = true
	}

	for i:=0; i<nodeLen; i=i+2 {
		data := append(merkleNodes[i].Data, merkleNodes[i+1].Data...)
		hash := md5.Sum(data)

		//生成新的merkleNode，并加入到返回值
		node := MerkleNode{merkleNodes[i],merkleNodes[i+1],hash[:]}
		returnMerkleNodes = append(returnMerkleNodes, &node)
	}

	//处理最后一个,如果是复制前一个的，则子节点置空
	if  isHaveCopy {
		merkleNodes[nodeLen - 1].Right = nil
		merkleNodes[nodeLen - 1].Left = nil
	}
	//递归进行
	return GetMerkleNode(returnMerkleNodes,false)
}
//GetProve 根据节点原始数据返回默克尔证明路径
func (t *MerkleTree)GetProve(node []byte) (prove [][]byte,err error){
	if t == nil{
		err = errors.New("默克尔树不存在")
		return
	}
	if len(node)==0{
		err = errors.New("原始数据不能为空")
		return
	}
	nodeHash := md5.Sum(node)
	hash := nodeHash[:]
	prove = t.Root.GetNodeProve(&hash,nil)
	return
}
func (n *MerkleNode)GetNodeProve(hash *[]byte,prove [][]byte) [][]byte{
	//相等，说明是原始数据
	if  bytes.Equal(n.Data,*hash){
		return append(prove,n.Data)
	}
	prove = append(prove,n.Data)
	//左节点或者右节点为空，说明是叶子节点
	if n.Left == nil || n.Right == nil{
		prove = nil
	}
	if l := n.Left;l != nil {
		p := l.GetNodeProve(hash,prove)
		if p != nil{
			return p
		}
	}
	if r := n.Right;r != nil {
		p := r.GetNodeProve(hash,prove)
		if p != nil{
			return p
		}
	}
	return nil
}
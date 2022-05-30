package merkletree

import (
	"fmt"
	"testing"
)

func TestMerkle(t *testing.T){
	//原始数据
	var data [][]byte
	data = [][]byte{
		[]byte("1"),
		[]byte("2"),
		[]byte("3"),
		[]byte("4"),
		[]byte("5"),
		[]byte("6"),
		[]byte("7"),
		[]byte("8"),
		[]byte("9"),
		[]byte("10"),
	}
	root := GetMerkleRoot(data)
	fmt.Printf("默克尔根：%x\n",root.Root.Data)
	//获取默克尔证明
	arrProve,err := root.GetProve([]byte("8"))
	if err != nil{
		fmt.Println("error:",err.Error())
		return
	}
	prove := []string{}
	for _,v := range arrProve{
		prove = append(prove,fmt.Sprintf("%x",v))
	}
	fmt.Println("默克尔证明路径:",prove)
}
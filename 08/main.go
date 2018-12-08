package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"strconv"
	"strings"
)

type Node struct {
	MetaData []int
	Children []*Node
}

func (n Node) WriteTo(w io.Writer, indent int) error {
	var metadata []string
	for _, x := range n.MetaData {
		metadata = append(metadata, strconv.Itoa(x))
	}
	_, err := fmt.Fprintf(w,
		"%sNode(children=%d, metadata=[%s])\n",
		strings.Repeat("  ", indent),
		len(n.Children),
		strings.Join(metadata, ", "),
	)
	if err != nil {
		return err
	}
	for _, child := range n.Children {
		if err := child.WriteTo(w, indent+1); err != nil {
			return err
		}
	}
	return nil
}

func (n Node) String() string {
	var b strings.Builder
	n.WriteTo(&b, 0)
	return b.String()
}

type Parser struct {
	nums  []int
	index int
	err   error
}

func NewParser(nums []int) *Parser {
	return &Parser{nums: nums}
}

func (p *Parser) Err() error { return p.err }

func (p *Parser) Done() bool {
	return p.err != nil || p.index >= len(p.nums)
}

func (p *Parser) Int() int {
	if p.err != nil {
		return 0
	}
	if p.Done() {
		p.err = io.EOF
		return 0
	}
	num := p.nums[p.index]
	p.index++
	return num
}

func (p *Parser) Node() *Node {
	if p.err != nil {
		return nil
	}
	var (
		node      = &Node{}
		children  = p.Int()
		metadatas = p.Int()
	)
	for i := 0; i < children; i++ {
		node.Children = append(node.Children, p.Node())
	}
	for i := 0; i < metadatas; i++ {
		node.MetaData = append(node.MetaData, p.Int())
	}
	return node
}

func (p *Parser) Root() *Node {
	node := p.Node()
	if !p.Done() {
		p.err = fmt.Errorf("trailing numbers after node")
	}
	return node
}

func ReadInput(file string) ([]int, error) {
	data, err := ioutil.ReadFile(file)
	if err != nil {
		return nil, err
	}
	var nums []int
	for _, s := range strings.Fields(string(data)) {
		n, err := strconv.Atoi(s)
		if err != nil {
			return nil, err
		}
		nums = append(nums, n)
	}
	return nums, nil
}

func main() {
	nums, err := ReadInput("input.txt")
	if err != nil {
		log.Fatal(err)
	}
	parser := NewParser(nums)
	root := parser.Root()
	if err := parser.Err(); err != nil {
		log.Fatal(err)
	}
	fmt.Println(root)
}

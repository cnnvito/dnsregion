package dnsregion

import (
	_ "embed"
	"encoding/json"
	"fmt"
	"strings"
)

//go:embed data/client_ips.json
var clientIPs []byte

var DefaultSubnetResource, _ = NewSubnetResouce(clientIPs)

func NewSubnetResouce(content []byte) (*SubnetResoucre, error) {
	var r SubnetResoucre
	if err := json.Unmarshal(content, &r); err != nil {
		return nil, err
	}
	return &r, nil
}

type groups struct {
	Name      string   `json:"name"`
	Childrens []string `json:"childrens"`
}

type node struct {
	Name string `json:"name"`
	Ip   string `json:"ip"`
}

func (n node) String() string {
	return fmt.Sprintf("\t%s\t%s", n.Name, n.Ip)
}

type SubnetResoucre struct {
	Groups    []groups        `json:"groups"`
	Childrens map[string]node `json:"childrens"`
}

func (r *SubnetResoucre) String() string {
	var result strings.Builder
	for _, group := range r.Groups {
		result.WriteString(group.Name + "\n")
		for _, node := range r.SearchChildren(group.Childrens) {
			result.WriteString(node.String() + "\n")
		}
	}
	return result.String()
}

func (r *SubnetResoucre) SearchChildren(keys []string) []node {
	var result []node
	for _, key := range keys {
		if node, ok := r.Childrens[key]; ok {
			result = append(result, node)
		}
	}
	return result
}

func (r *SubnetResoucre) SearchGroupChildrenKeys(name string) []string {
	var result []string
	for _, group := range r.Groups {
		if group.Name == name {
			return append(result, group.Childrens...)
		}
	}
	return result
}

func (r *SubnetResoucre) SearchGroupChildrenNodes(name string) []node {
	keys := r.SearchGroupChildrenKeys(name)
	if len(keys) == 0 {
		return nil
	}
	return r.SearchChildren(keys)
}

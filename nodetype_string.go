// generated by stringer -type=nodeType; DO NOT EDIT

package haddoque

import "fmt"

const _nodeType_name = "nodeChainnodeListnodeBoolnodeTextnodeNumbernodeWherenodeAndnodeOrnodeInnodeContains"

var _nodeType_index = [...]uint8{0, 9, 17, 25, 33, 43, 52, 59, 65, 71, 83}

func (i nodeType) String() string {
	if i < 0 || i+1 >= nodeType(len(_nodeType_index)) {
		return fmt.Sprintf("nodeType(%d)", i)
	}
	return _nodeType_name[_nodeType_index[i]:_nodeType_index[i+1]]
}

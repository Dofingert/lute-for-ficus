// Lute - 一款结构化的 Markdown 引擎，支持 Go 和 JavaScript
// Copyright (c) 2019-present, b3log.org
//
// Lute is licensed under Mulan PSL v2.
// You can use this software according to the terms and conditions of the Mulan PSL v2.
// You may obtain a copy of Mulan PSL v2 at:
//         http://license.coscl.org.cn/MulanPSL2
// THIS SOFTWARE IS PROVIDED ON AN "AS IS" BASIS, WITHOUT WARRANTIES OF ANY KIND, EITHER EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO NON-INFRINGEMENT, MERCHANTABILITY OR FIT FOR A PARTICULAR PURPOSE.
// See the Mulan PSL v2 for more details.

package render

import (
	"encoding/json"

	"github.com/Dofingert/lute-for-ficus/ast"
	"github.com/Dofingert/lute-for-ficus/parse"
	"github.com/Dofingert/lute-for-ficus/util"
)

type JSONRenderer struct {
	*BaseRenderer
}

func NewJSONRenderer(tree *parse.Tree, options *Options) Renderer {
	// 渲染器剔除语法树块级 IAL 节点 https://github.com/88250/protyle/issues/1
	var ials []*ast.Node
	ast.Walk(tree.Root, func(n *ast.Node, entering bool) ast.WalkStatus {
		if !entering || ast.NodeKramdownBlockIAL != n.Type {
			return ast.WalkContinue
		}
		ials = append(ials, n)
		return ast.WalkContinue
	})
	for _, ial := range ials {
		ial.Unlink()
	}

	ret := &JSONRenderer{NewBaseRenderer(tree, options)}
	ret.DefaultRendererFunc = ret.renderNode
	return ret
}

func (r *JSONRenderer) renderNode(node *ast.Node, entering bool) ast.WalkStatus {
	if entering {
		if nil != node.Previous {
			r.WriteString(",")
		}
		node.Data, node.TypeStr = util.BytesToStr(node.Tokens), node.Type.String()
		node.Properties = ial2Map(node.KramdownIAL)
		delete(node.Properties, "refcount")
		data, err := json.Marshal(node)
		node.Data, node.TypeStr = "", ""
		node.Properties = nil
		if nil != err {
			panic("marshal node to json failed: " + err.Error())
			return ast.WalkStop
		}
		n := util.BytesToStr(data)
		n = n[:len(n)-1] // 去掉结尾的 }
		r.WriteString(n)
		if nil != node.FirstChild {
			r.WriteString(",\"Children\":[")
		} else {
			r.WriteString("}")
		}
	} else {
		if nil != node.FirstChild {
			r.WriteByte(']')
			r.WriteString("}")
		}
	}
	return ast.WalkContinue
}

func ial2Map(ial [][]string) (ret map[string]string) {
	ret = map[string]string{}
	for _, kv := range ial {
		ret[kv[0]] = kv[1]
	}
	return
}

package stmt

import (
	"fmt"
	"io"

	"github.com/z7zmey/php-parser/node"
	"github.com/z7zmey/php-parser/token"
)

type If struct {
	node.SimpleNode
	token  token.Token
	cond   node.Node
	stmts  []node.Node
	elseIf []node.Node
	_else  node.Node
}

func NewIf(token token.Token, cond node.Node, stmts []node.Node) node.Node {
	return If{
		node.SimpleNode{Name: "If", Attributes: make(map[string]string)},
		token,
		cond,
		stmts,
		nil,
		nil,
	}
}

func (n If) AddElseIf(elseIf node.Node) node.Node {
	if (n.elseIf == nil) {
		n.elseIf = make([]node.Node, 0)
	}

	n.elseIf = append(n.elseIf, elseIf)

	return n
}

func (n If) SetElse(_else node.Node) node.Node {
	n._else = _else

	return n
}

func (n If) Print(out io.Writer, indent string) {
	fmt.Fprintf(out, "\n%v%v [%d %d] %q", indent, n.Name, n.token.StartLine, n.token.EndLine, n.token.Value)

	if n.cond != nil {
		fmt.Fprintf(out, "\n%vcond:", indent+"  ")
		n.cond.Print(out, indent+"    ")
	}

	if n.stmts != nil {
		fmt.Fprintf(out, "\n%vstmts:", indent+"  ")
		for _, nn := range n.stmts {
			nn.Print(out, indent+"    ")
		}
	}
	
	if n.elseIf != nil {
		fmt.Fprintf(out, "\n%velseIfs:", indent+"  ")
		for _, nn := range n.elseIf {
			nn.Print(out, indent+"    ")
		}
	}
	if n._else != nil {
		fmt.Fprintf(out, "\n%velse:", indent+"  ")
		n._else.Print(out, indent+"    ")
	}
}
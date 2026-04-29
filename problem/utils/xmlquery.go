package utils

import (
	"fmt"
	"strings"

	"github.com/antchfx/xmlquery"
)

func XmlqueryPrettyPrint(node *xmlquery.Node, indent string) string {
	var b strings.Builder
	switch node.Type {
	case xmlquery.ElementNode:
		b.WriteString(fmt.Sprintf("%s<%s", indent, node.Data))
		for _, attr := range node.Attr {
			b.WriteString(fmt.Sprintf(` %s="%s"`, attr.Name.Local, attr.Value))
		}
		if node.FirstChild == nil {
			b.WriteString("/>\n")
		} else {
			b.WriteString(">\n")
			for child := node.FirstChild; child != nil; child = child.NextSibling {
				b.WriteString(XmlqueryPrettyPrint(child, indent+"  "))
			}
			b.WriteString(fmt.Sprintf("%s</%s>\n", indent, node.Data))
		}
	case xmlquery.TextNode:
		text := strings.TrimSpace(node.Data)
		if text != "" {
			b.WriteString(fmt.Sprintf("%s%s\n", indent, text))
		}
	case xmlquery.CommentNode:
		b.WriteString(fmt.Sprintf("%s<!-- %s -->\n", indent, node.Data))
	}
	return b.String()
}

package gedview

import (
	"bufio"
	"io"
	"regexp"
	"strconv"
	"strings"
)

type Node struct {
	Raw      string
	Level    int
	Type     string
	Value    string
	Children []*Node

	parent *Node
}

// CreateAST creates based off of a GEDCOM formated file a AST which can be
// used to interpret the structure and ultimately work with the data more
// appropriately.
func CreateAST(contents io.Reader) (*Node, error) {
	s := bufio.NewScanner(contents)

	root := &Node{Level: -1, Children: []*Node{}, Raw: "ROOT"}
	current := root

	for s.Scan() {
		line := s.Text()
		re := regexp.MustCompile(`^(\d+)\s([a-zA-Z\@1-9]+)(\s(.*))?$`)
		matches := re.FindStringSubmatch(line)

		if len(matches) > 0 {
			lvl, err := strconv.Atoi(matches[1])
			if err != nil {
				// this should never happend since the regex already
				// confirmed this is a proper numeric value
				return nil, err
			}

			n := &Node{
				Raw:   line,
				Level: lvl,
				Type:  matches[2],
				Value: strings.Trim(matches[3], " "),
			}

			switch {
			case lvl > current.Level:
				n.parent = current
				current.Children = append(current.Children, n)
			case lvl == current.Level:
				n.parent = current.parent
				current.parent.Children = append(current.parent.Children, n)
			case lvl < current.Level:
				p := current
				for lvl <= p.Level {
					p = p.parent
				}

				n.parent = p
				p.Children = append(p.Children, n)
			}

			current = n
		}
	}

	return root, nil
}

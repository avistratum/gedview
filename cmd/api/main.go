package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"gedview"
	"log"
	"net/http"
	"os"
	"regexp"
)

var path *string = flag.String("path", "", "path to GEDCOM file to derive the API from")

type PersonNode struct {
	ID   string `json:"id"`
	Name string `json:"string"`
}

type FamilyNode struct {
	ID       string       `json:"id"`
	Spouses  []PersonNode `json:"spouses`
	Children []PersonNode `json:"children"`
}

func main() {
	flag.Parse()
	file, err := os.Open(*path)
	if err != nil {
		panic(err)
		return
	}

	tree, err := gedview.CreateAST(file)
	if err != nil {
		panic(err)
		return
	}

	idRe := regexp.MustCompile(`^@I(\d+)@$`)
	indi := tree.FindBy(func(n *gedview.Node) bool {
		return idRe.MatchString(n.Type)
	})

	persons := map[string]PersonNode{}
	for _, n := range indi {
		m := idRe.FindStringSubmatch(n.Type)
		persons[m[1]] = PersonNode{
			ID:   m[1],
			Name: n.FindByType("NAME")[0].Value,
		}
	}

	fidRe := regexp.MustCompile(`^@F(\d+)@$`)
	fams := tree.FindBy(func(n *gedview.Node) bool {
		return fidRe.MatchString(n.Type)
	})

	families := map[string]FamilyNode{}
	for _, n := range fams {
		m := fidRe.FindStringSubmatch(n.Type)

		chil := n.FindByType("CHIL")
		children := []PersonNode{}
		for _, n := range chil {
			m := idRe.FindStringSubmatch(n.Value)
			children = append(children, persons[m[1]])
		}

		sp := n.FindBy(func(n *gedview.Node) bool {
			switch n.Type {
			case "WIFE":
				fallthrough
			case "HUSB":
				return true
			default:
				return false
			}
		})
		spouses := []PersonNode{}
		for _, n := range sp {
			m := idRe.FindStringSubmatch(n.Value)
			spouses = append(spouses, persons[m[1]])
		}

		families[m[1]] = FamilyNode{
			ID:       m[1],
			Children: children,
			Spouses:  spouses,
		}
	}

	http.HandleFunc("/families", func(w http.ResponseWriter, r *http.Request) {
		enc := json.NewEncoder(w)
		enc.Encode(families)
	})

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Println(fmt.Sprintf("server started; listening on ':%v'", port))

	err = http.ListenAndServe(":"+port, nil)
	if err != nil {
		log.Printf("%w", err)
	}
}

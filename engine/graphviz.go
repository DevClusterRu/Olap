package engine

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/goccy/go-graphviz"
	"github.com/goccy/go-graphviz/cgraph"
	"log"
	"math/rand"
	"strconv"
	"time"
)



func DoGraph(agr AggregationResult) string {

	g := graphviz.New()
	graph, err := g.Graph()
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		if err := graph.Close(); err != nil {
			log.Fatal(err)
		}
		g.Close()
	}()

	nodesPointers := make(map[string]*cgraph.Node)

	nodesSorted:=agr.NodesSorting
	nodes:=agr.NodesMap
	links:=agr.LinksMap

	for _,v:=range nodesSorted{

		if v==""{
			continue
		}

		nodeName:=v
		if nodes[v]>0{
			nodeName+=` [`+strconv.Itoa(nodes[v])+`]`
		}

		nodesPointers[v], err = graph.CreateNode(nodeName)
		if err != nil {
			log.Println(err)
		}
	}
	fmt.Println("All ok")

	for k,v:=range links {
		names:=bytes.Split([]byte(k),[]byte("_"))

		e, err := graph.CreateEdge(k, nodesPointers[string(names[0])], nodesPointers[string(names[1])])
		if err != nil {
			log.Println(err)
		}
		if v>0 {
			e.SetLabel(strconv.Itoa(v))
		}
	}

	var buf bytes.Buffer
	if err := g.Render(graph, "dot", &buf); err != nil {
		log.Println(err)
	}
	fmt.Println(buf.String())

	tm:=time.Now()
	salt:=strconv.Itoa(rand.Intn(9999))

	//way:=""
	fname := tm.Format("20060102150405"+salt)+".svg"

	if err := g.RenderFilename(graph, graphviz.SVG, Way+fname); err != nil {
		log.Println(err)
	}

	result:=JsonReturn{
		Filename: fname,
		Events: agr.AllEvents,
		Objects: agr.AllObjects,
	}

	b, err:=json.Marshal(result)
	if err!=nil{
		log.Println("Error when marshall result -->",err.Error())
	}


	return string(b)

}
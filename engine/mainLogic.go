package engine

import "go.mongodb.org/mongo-driver/bson"

func setStartNode(i int, nodesMap map[string]int, linksMap map[string]int, nodesSorting []string, data []bson.M) ([]string, map[string]int, map[string]int) {
	nodesSorting = append(nodesSorting, "[Start]")
	nodeFrom := "[Start]"
	nodeTo := data[i]["event"].(string)
	nodesMap[nodeFrom] = 0
	nodesMap[nodeTo] = 0
	linksMap[nodeFrom+"_"+nodeTo] = 0
	nodesSorting = append(nodesSorting, nodeTo)

	return nodesSorting, nodesMap, linksMap
}

func setEndNode(i int,  data []bson.M, linksMap map[string]int, nodesSorting []string) ([]string, map[string]int){
	nodeFrom := data[i]["event"].(string)
	nodeTo := "[End]"
	nodesSorting = append(nodesSorting, nodeFrom)
	nodesSorting = append(nodesSorting, nodeTo)
	linksMap[nodeFrom+"_"+nodeTo] = 0

	return nodesSorting, linksMap
}



func BaseCalculation(agr AggregationResult) AggregationResult {

	var nodesSorting []string //Слайс для нод, сложенных по порядку
	var nodesMap, linksMap map[string]int
	var data []bson.M

	data = agr.Nodes

	nodesMap = make(map[string]int) //Карта нод хаотичная, имя в ключе и количество в значении
	linksMap = make(map[string]int) //Карта линков хаотичная, имя в ключе как node1_node2 и количество в значении

	for i := -1; i < len(data)-1; i++ {

		if i == -1 { //Определяем стартовую ноду и устанавливаем связь между ней и первой реальной нодой
			nodesSorting, nodesMap, linksMap = setStartNode(i + 1, nodesMap, linksMap, nodesSorting, data)
			continue
		}

		if data[i]["object"] != data[i+1]["object"] {
			nodesSorting, nodesMap, linksMap = setStartNode(i + 1, nodesMap, linksMap, nodesSorting, data)
			nodesSorting, linksMap = setEndNode(i,data, linksMap, nodesSorting)
			continue
		}

		if data[i]["event"] != data[i+1]["event"] {
			//Берем следующий 2 ноды
			nodeFrom := data[i]["event"].(string)
			nodeTo := data[i+1]["event"].(string)

			//Если в картах нет этих нод, создаем их
			if _, found := nodesMap[nodeFrom]; found == false {
				nodesSorting = append(nodesSorting, nodeFrom)
				nodesMap[nodeFrom] = 0
			}
			if _, found := nodesMap[nodeTo]; found == false {
				nodesSorting = append(nodesSorting, nodeTo)
				nodesMap[nodeTo] = 0
			}

			if _, found := linksMap[nodeFrom+"_"+nodeTo]; found == true {
				//Если уже есть такая связь, инкрементируем, если нет - устанавливаем
				linksMap[nodeFrom+"_"+nodeTo]++
			} else {
				linksMap[nodeFrom+"_"+nodeTo] = 1
			}
			//Теперь важно! Нода, в которую ведет связь должна быть инкрементирована!
			nodesMap[nodeTo] += 1
		}

	}

	if len(data) >= 1 {
		nodesSorting, linksMap = setEndNode(len(data) - 1,data, linksMap, nodesSorting)
	}

	//fmt.Println(nodesMap)
	//fmt.Println(linksMap)
	agr.NodesMap = nodesMap
	agr.LinksMap = linksMap
	agr.NodesSorting = nodesSorting

	return agr
}

package engine


func BaseCalculation(agr AggregationResult) AggregationResult{

	data:=agr.Nodes

	nodesSorting := make([]string, len(data)+3) //Массив для нод, сложенных по порядку
	nodesMap := make(map[string]int) //Карта нод хаотичная, имя в ключе и количество в значении
	linksMap := make(map[string]int) //Карта линков хаотичная, имя в ключе как node1_node2 и количество в значении

	currentPositionInNodesSorting:=0

	for i := -1; i < len(data)-1; i++ {

		if i==-1{ //Определяем стартовую ноду и устанавливаем связь между ней и первой реальной нодой
			nodesSorting[currentPositionInNodesSorting] = "[Start]"
			currentPositionInNodesSorting++
			nodeFrom := "[Start]"
			nodeTo := data[i+1]["event"].(string)
			nodesMap[nodeFrom] = 0
			nodesMap[nodeTo] = 0
			linksMap[nodeFrom+"_"+nodeTo] = 0
			nodesSorting[currentPositionInNodesSorting] = nodeTo
			currentPositionInNodesSorting++
			continue
		}

		//Берем следующий 2 ноды
		nodeFrom := data[i]["event"].(string)
		nodeTo := data[i+1]["event"].(string)

		//Если в картах нет этих нод, создаем их
		if  _, found := nodesMap[nodeFrom]; found==false{
			nodesSorting[currentPositionInNodesSorting] = nodeFrom //Пишем в мамссив для сохранения порядка
			currentPositionInNodesSorting++
			nodesMap[nodeFrom]=0
		}
		if  _, found := nodesMap[nodeTo]; found==false{
			nodesSorting[currentPositionInNodesSorting] = nodeTo //Пишем в мамссив для сохранения порядка
			currentPositionInNodesSorting++
			nodesMap[nodeTo]=0
		}

		if  _, found := linksMap[nodeFrom+"_"+nodeTo]; found==true{
			//Если уже есть такая связь, инкрементируем, если нет - устанавливаем
			linksMap[nodeFrom+"_"+nodeTo]++
		} else {
			linksMap[nodeFrom+"_"+nodeTo] = 1
		}

		//Теперь важно! Нода, в которую ведет связь должна быть инкрементирована!
		nodesMap[nodeTo]+=1

	}

	if len(data)>=1 {

		nodeFrom := data[len(data)-1]["event"].(string)
		nodeTo := "[End]"

		nodesSorting[len(data)+1] = nodeFrom
		nodesSorting[len(data)+2] = nodeTo

		linksMap[nodeFrom+"_"+nodeTo] = 0
	}


	//fmt.Println(nodesMap)
	//fmt.Println(linksMap)
	agr.NodesMap = nodesMap
	agr.LinksMap = linksMap
	agr.NodesSorting = nodesSorting


	return agr
}

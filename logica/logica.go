package logica

import (
		"log"
		"strings"
		"regexp"
)
// Logica is a structure for child
type Logica struct { 
	Label string				 `json:"Label"`
	state int
	origin int
	visibleChildren int
	Children map[string] *Logica  
}

type (
	// Scenario is a DICT type struture for holding instances of state variable
	Scenario map[string] string
	// ScenarioList holds an array of scenarios
	ScenarioList []Scenario
)

// PatternLogicaNodeSuffix used to identify prefix - dimension and suffix (residual)
const PatternLogicaNodeSuffix = "^(\\.|\\!)([a-zA-Z_\\[\\]0-9]+)([.|!]\\S+)?$"

const ( // ioat is reset to 0
	// LogicaLoggingOn = 0
	LogicaLoggingOn = 0
	// LogicaLoggingOff = 1
	LogicaLoggingOff = 1
)

const (
	// LogicaOriginal when node was available at logging and not deleted
	LogicaOriginal = 0
	// LogicaDeleted when node was available at logging and deleted
	LogicaDeleted = 1
	// LogicaCreated when node wasn't available at logging
	LogicaCreated = 2
)

func find(path string, pattern string) (string, string, string) {
	var re = regexp.MustCompile(pattern)
	var results = re.FindStringSubmatch(path)
	//log.Println(results)
	//log.Println(len(results))
	if len(results)==4 {return results[1],results[2],results[3]}
	return "","",""
} 

// CreateLogica is a Logica constructor
func CreateLogica(name string) Logica{
	//name = strings.Replace(name,".","",-1)
	return Logica{Label:name, state: LogicaLoggingOff, visibleChildren: 0, origin: LogicaOriginal, Children: make(map[string]*Logica)}
}

// CreateScenario creates a new dict for scenario structure
func CreateScenario() Scenario {
	return make(Scenario)
}

// CreateScenarioList creates a list for scenario structures
func CreateScenarioList() ScenarioList {
	return 	make(ScenarioList,0)
}

// Copy a Scenario so not overwriting in lists
func (l Scenario) Copy() Scenario {
	var newScene = CreateScenario()
	for k,v := range l {
		newScene[k] = v
	}
	return newScene
}

// Output for Scenario
func (l Scenario) Output() string {
	var output = "{"
	for k,v := range l {
		output += " "+k+":"+v+" "
	}
	output += "}\n"
	return output
}

func (l ScenarioList) hasParameter(paramname string) bool {
	if (len(l)>0) {
		for _, scene := range l {
			_, ok := scene[paramname]
			return ok
		}
	}
	return false
}

// Output for Scenario
func (l ScenarioList) Output() string {
	var output = "[\n"
	for _,k := range l {
		output += " "+k.Output()
	}
	output += "]\n"
	return output
}

// StartLogging initiates logging of all changes
func (l *Logica) StartLogging() bool {
	l.state = LogicaLoggingOn
	l.visibleChildren = len(l.Children)
	// log.Printf("%s is logging .. \n",l.Label)
	for _,v := range l.Children {
		v.StartLogging()
	}
	return false
}

// Revert reverts Logica to state before StartLogging
func (l *Logica) Revert() bool {
	l.state = LogicaLoggingOff
	for k,v := range l.Children {
		v.Revert()
		switch v.origin {
		case LogicaDeleted:
			v.origin = LogicaOriginal
			if strings.Index(v.Label,"__") == 0 {
				strings.Replace(v.Label,"__","",1)
				l.Children[v.Label] = v
				delete(l.Children,k)
			}

			break
		case LogicaCreated:
			l.deleteChild(k)
			break
		}
	}
	return false
}

// Has returns true if Logica has path
func (l *Logica) Has(path string) bool {
	var logic = l.Get(path)
	return (logic != nil)
}

// Get returns Logica Object if Logica is on path
func (l *Logica) Get( path string) *Logica {
	prefix,dimension,suffix := find(path,PatternLogicaNodeSuffix)
	if prefix != "" {
		value, ok := l.Children[dimension]
		if !ok || value.origin == LogicaDeleted {return nil}
		if l.state == LogicaLoggingOff {l.visibleChildren = len(l.Children)}
		if (prefix=="!" && l.visibleChildren!=1) {return nil}
		if (suffix != ""){
			return value.Get(suffix)
		} 
		return value
	} 
	log.Printf("Wrong path format:\n\t %s",path)
	panic("Wrong format")
	//return nil
}

// Add adds a Logica Child Path and returns leaf node
func (l *Logica) Add(path string) *Logica {
	prefix,dimension,suffix := find(path,PatternLogicaNodeSuffix)
	if prefix != "" {
		//log.Println(dimension)

		if prefix == "!" {
			l.Clear(dimension)
		}
		value, ok := l.Children[dimension]
		if ok && value.origin != LogicaDeleted {
			if suffix != "" {
				return value.Add(suffix)
			}
			return l
		}
		//newchild := CreateLogica(dimension)
		newchild := l.createChild(dimension)
		l.Children[dimension] = &newchild
		if suffix != "" {
			return newchild.Add(suffix)
		}
		return &newchild
	}
	log.Printf("Wrong path format:\n\t %s",path)
	panic("Wrong format")
}

// Pop deletes a Logica Child Path
func (l *Logica) Pop( path string) {
	prefix,dimension,suffix := find(path,PatternLogicaNodeSuffix)
	if prefix != "" {
		value, ok := l.Children[dimension]
		if ok && value.origin != LogicaDeleted {
			if suffix != "" {
				value.Pop(suffix)
			} else {
				value.Clear("")
				l.deleteChild(dimension)
			}
		} else {
			log.Printf("Could not find:\n\t %s of %s in %s",path,dimension,l.Label)
		}
	} else {
		log.Printf("Wrong path format:\n\t %s",path)
		panic("Wrong format")
	}
}

func (l *Logica) createChild(name string) Logica {
	l.visibleChildren++
	newChild := CreateLogica(name)
	if l.state == LogicaLoggingOn {
		newChild.origin = LogicaCreated
		newChild.state = LogicaLoggingOn
		oldChild, oldOk := l.Children[name]
		_, bkOk := l.Children["__"+name]
		if oldOk && !bkOk {
			l.Children["__"+name] = oldChild
			delete(l.Children,name)
		}
	}
	return newChild
}

func (l *Logica) deleteChild(name string) {
	if (l.Children[name].origin != LogicaDeleted) {l.visibleChildren--}
	child , ok := l.Children[name]
	if ok {
		if l.state == LogicaLoggingOn && child.origin != LogicaCreated  {
			child.origin = LogicaDeleted
			l.Children[name] = child
		} else {
			delete(l.Children,name)
		}
	}
}

func (l *Logica) hasChild( childname string) bool {
	value, ok := l.Children[childname]
	return ok && value.origin != LogicaDeleted
}

// Clear removes all children except for exception
func (l *Logica) Clear(exception string) {
	for k,v := range l.Children {
		if (k!= exception) {
			v.Clear("")
			l.deleteChild(k)
		}
	}
}

// Output returns a string of the Logica Tree structure
func (l *Logica) Output(input string, level int) string {
	if l.origin == LogicaDeleted {input +="\033[0;31m"}
	if l.origin == LogicaCreated {input +="\033[0;32m"}
	input += strings.Repeat("  ",level)+l.Label+"\n"
	if l.origin != LogicaOriginal {input +="\033[0m"}
	level++
	for _, v := range l.Children { 
		input += v.Output("",level)
	}
	return input	
}

func isParameter(text string) (string, bool) {
	if (text[0] == '[' && text[len(text)-1] == ']') {
        return text[1:len(text)-1], true
    }
    return text, false
}

// Parameters collects all scenarios from given Logica Path
func (l *Logica) Parameters( path string, space ScenarioList ) ScenarioList {
	if len(space)==0 {
		var emptyScene = CreateScenario()
		return l.parameters(path,emptyScene,space)
	} else {
		var tmp = CreateScenarioList()
		for _,scene := range space {
			tmp = l.parameters(path,scene,tmp)
			//log.Println(tmp.Output())
		}
		space = tmp
	}
	return space
}

func (l *Logica) parameters( path string, scene Scenario, space ScenarioList) ScenarioList {
	prefix,branch,nextpath := find(path,PatternLogicaNodeSuffix)
	if (prefix != ""){
		// There is a prefix and thus a branch
		if branch, ok := isParameter(branch); ok {
			//The branch is an parameter
			/*
			if space.hasParameter(branch) {
				// We already have the parameter - removal of all scenes not in branch dimension
				for _, scene := range space {
					parmvar := scene[branch]
					child := l.Children[parmvar]
					if child != nil {
						// only add children in paramspace and in child dimension
						child.Parameters(nextpath,scene,space)
					}
				}
				
			} else {*/
				// New parameter - defined by branch dimensions
				if l.state == LogicaLoggingOff {l.visibleChildren = len(l.Children)}
				if prefix != "!" || l.visibleChildren == 1 {
					for key, child := range l.Children {
						if (child.origin != LogicaDeleted || l.state == LogicaLoggingOff) {
							value,ok := scene[branch]
							if ok==false || (key==value) {
								//log.Println("true," + scene.Output())
								var scene2 = scene.Copy()
								scene2[branch] = key
								space = child.parameters(nextpath,scene2,space)
							} else {
								//log.Println("false," + scene.Output())
							}
						}
					}
				}
			//}
		} else if l.hasChild(branch){
			//The branch is a child and we continue down the branch
			if (prefix != "!" || l.visibleChildren == 1){
				space = l.Children[branch].parameters(nextpath, scene, space)
			}
		}
	} else {
		// We are at end of branch - at a leaf - and add scene
		if len(scene) > 0 {
			space = append(space,scene.Copy())
		}
	}
	return space
}

func isSameScenario(sceneA Scenario, sceneB Scenario, parameters []string) bool {
	for _,parameter := range parameters {
		if sceneA[parameter] != sceneB[parameter] {
			return false
		}
	}
	return true
}


/*
bool isSameScenario(scenario elementA, scenario elementB, std::list<string> vars){
    for (auto varname : vars){
        if (elementA[varname] != elementB[varname]) return false;
    }
    return true;
}

scenarioList* intersection(scenarioList listA, scenarioList listB, std::list<string> vars, std::list<string> uniquevars)
{
    scenarioList *output = new scenarioList();
    if (!listA.empty())
    {
        for (auto elementA : listA)
        {
            for (auto elementB : listB)
            {
                if (isSameScenario(elementA, elementB, vars))
                {
                    scenario newElementA = elementA;
                    for (auto xtra : uniquevars)
                    {
                        newElementA[xtra] = elementB[xtra];
                    }
                    //output->push_back(newElementA);
                    output->insert(output->end(), newElementA);
                }
            }
        }
    }
    else
        {
            for (auto elementB : listB)
            {
                //output->push_back(elementB);
                output->insert(output->end(), elementB);
            }
        }
        return output;

    }

    void splitvars(scenarioList listA, scenarioList listB, std::list<string> & commonvars, std::list<string> & uniquevars)
    {
        if (!listB.empty() && !listA.empty())
        {
            auto elementB = listB.begin();
            auto elementA = listA.begin();
            for (auto key : *elementB)
            {
                if (elementA->find(key.first.c_str()) != elementA->end())
                {
                    //commonvars.push_back(key.first);
                    commonvars.insert(commonvars.end(), key.first);
                }
                else
                {
                    //uniquevars.push_back(key.first);
                    uniquevars.insert(uniquevars.end(), key.first);
                }
            }
    }
    else if (listA.empty())
    {

        //uniquevars.splice(uniquevars.end(), listB);
         
        auto elementB = listB.begin();
        for (auto key : *elementB)
        {
            //uniquevars.push_back(key.first);
            uniquevars.insert(uniquevars.end(), key.first);
        }
        
    }
}

scenarioList* combineLists(scenarioList listA, scenarioList listB)
{
    std::list<string> commonvars, uniquevars;
    splitvars(listA, listB, commonvars, uniquevars);
    scenarioList *newList = intersection(listA, listB, commonvars, uniquevars);
    //multiplyList(newListA,newListB,commonvars,uniquevars);
    return newList;
}


*/
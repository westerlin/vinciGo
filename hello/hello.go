package main

import (
	"fmt"
	"log"
	fruit "fruit"
	"logica"
	//"vincireader"
)




func main() {
	fmt.Printf("hello, world\n")
	fruit.Myfunc()
	fmt.Printf(" Test %d\n", fruit.Myinteger())
	log.Printf("My test")
	var mylogica = logica.CreateLogica("<root>")
	mylogica.Add(".actors").Add(".Lucy").Add(".Gender").Add(".Female")
	mylogica.Add(".actors.Miranda.Gender.Female")
	mylogica.Add(".actors.Lucretia.Gender.Female")
	mylogica.Add(".actors.Silvia.Gender.Female")
	mylogica.Add(".locations").Add(".church")
	mylogica.Add(".locations.barn")
	mylogica.Add(".locations.house.Miranda")
	mylogica.Add(".locations.farm.Lucretia")
	mylogica.Add(".locations.house.Silvia")
	mylogica.Add(".locations.house.Peter")
	mylogica.StartLogging()
	mylogica.Add(".locations.barn!old_one.planks")
	mylogica.Add(".locations.barn!old_one.wood")
	mylogica.Add(".locations.house")

	//fmt.Println(mylogica)
	fmt.Println(mylogica.Output("",0)) 
	if mylogica.Has(".locations.house") {
		fmt.Println("House was there")
	} else {
		fmt.Println("House was not there")
	}
	mylogica.Clear("actors")
	mylogica.Add(".actors.Miranda.Age.28")
	fmt.Println(mylogica.Output("",0)) 
	mylogica.Pop(".actors.Miranda.Gender")
	fmt.Println(mylogica.Output("",0)) 
	//vincireader.ReadFile()
	// fmt.Println("",mylogica.Has(".actors.Miranda!Age"))
	// fmt.Println("",mylogica.Has(".actors.Miranda.Age"))
	// fmt.Println("",mylogica.Has(".actors.Miranda!Gender"))
	// fmt.Println("",mylogica.Has(".actors.Lucy!Gender"))

	mylogica.Revert()
	fmt.Println(mylogica.Output("",0)) 

	var scene = logica.CreateScenario()
	scene["Actors"] = "Lucy"
	scene["Location"] = "Barn"
	scene["Creator"] = "Rose"
	var scene2 = logica.CreateScenario()
	scene2["Actors"] = "Amanda"
	scene2["Location"] = "Garden"
	scene2["Creator"] = "Charles"
	var scenelist = make(logica.ScenarioList,0)
	scenelist = append(scenelist,scene)
	scenelist = append(scenelist,scene2)
	fmt.Println(scene.Output())
	fmt.Println(scenelist.Output())
	scene = logica.CreateScenario()
	scenelist = logica.CreateScenarioList()
	var scene1 = logica.CreateScenario()
	scene1["murderess"] = "Miranda"
	scene1["actor"] = "Miranda"
	scenelist = append(scenelist,scene1)

	scene1 = logica.CreateScenario()
	scene1["murderess"] = "Miranda"
	scene1["actor"] = "Lucretia"
	scenelist = append(scenelist,scene1)

	mylogica.Add(".actors.Charles.Gender.Male")
	mylogica.Add(".actors.Peter.Gender.Male")
	mylogica.Add(".actors.Peter.Gender.Female")
	scenelist = mylogica.Parameters(".actors.[actor].Gender.Female",scenelist)
	fmt.Println(scenelist.Output())
	scenelist = mylogica.Parameters(".locations.house.[actor]",scenelist)
	fmt.Println(scenelist.Output())
}

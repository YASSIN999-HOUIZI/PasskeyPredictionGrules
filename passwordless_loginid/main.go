package main

import (
	"fmt"

	"github.com/hyperjumptech/grule-rule-engine/ast"
	"github.com/hyperjumptech/grule-rule-engine/builder"

	"github.com/hyperjumptech/grule-rule-engine/engine"
	"github.com/hyperjumptech/grule-rule-engine/pkg"
)

func main(){
	devices, err := readDeviceInfoFromJSON("./devices.json")
	if err != nil {
		panic(err)
	}
	

	drls := rules

	knowledgeLibrary := ast.NewKnowledgeLibrary()
	ruleBuilder := builder.NewRuleBuilder(knowledgeLibrary)
	bs := pkg.NewBytesResource([]byte(drls))
	err = ruleBuilder.BuildRuleFromResource("PasswordLessRules", "0.0.2", bs)
	if err != nil {
		panic(err)
	}

	knowledgeBase, _ := knowledgeLibrary.NewKnowledgeBaseInstance("PasswordLessRules", "0.0.2")
	engine := engine.NewGruleEngine()

	
	for _, myFact := range devices {
		// Initialize the intermediate predictions for each iteration
		
		helper1 := &IntermediateCloudPrediction{
			Prediction: "",
		}
		helper2 := &IntermediatePasskeyPrediction{
			Prediction: "",
		}
		dataCtx := ast.NewDataContext()
		err = dataCtx.Add("DF", myFact)
		if err != nil {
			panic(err)
		}

		err = dataCtx.Add("IntermediateCloudPrediction", helper1)
		if err != nil {
			panic(err)
		}
		err = dataCtx.Add("IntermediatePasskeyPrediction", helper2)
		if err != nil {
			panic(err)
		}


		err = engine.Execute(dataCtx, knowledgeBase)
		if err != nil {
			panic(err)
		}
		myFact.Output.MatchProbability = myFact.MatchProbability()
		
	}

	err = writeDeviceInfoToCSV("device_info.csv", devices[:]) // Pass slice to avoid modifying original array
	if err != nil {
		panic(err)
	}

	fmt.Println("Device information saved to device_info.csv")
}
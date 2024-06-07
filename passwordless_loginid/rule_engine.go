package main

import (
	"fmt"

	"github.com/hyperjumptech/grule-rule-engine/ast"
	"github.com/hyperjumptech/grule-rule-engine/builder"
	"github.com/hyperjumptech/grule-rule-engine/engine"
	"github.com/hyperjumptech/grule-rule-engine/pkg"
)

// ProcessDevices processes the devices using the rule engine
func ProcessDevices(devices []*DeviceFact, rules string) error {
	
	knowledgeLibrary := ast.NewKnowledgeLibrary()
	ruleBuilder := builder.NewRuleBuilder(knowledgeLibrary)
	bs := pkg.NewBytesResource([]byte(rules))
	err := ruleBuilder.BuildRuleFromResource("PasswordLessRules", "0.0.2", bs)
	if err != nil {
		return err
	}

	knowledgeBase, _ := knowledgeLibrary.NewKnowledgeBaseInstance("PasswordLessRules", "0.0.2")
	ruleEngine := engine.NewGruleEngine()

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
			return err
		}

		err = dataCtx.Add("IntermediateCloudPrediction", helper1)
		if err != nil {
			return err
		}
		err = dataCtx.Add("IntermediatePasskeyPrediction", helper2)
		if err != nil {
			return err
		}

		err = ruleEngine.Execute(dataCtx, knowledgeBase)
		if err != nil {
			return err
		}
		myFact.Output.MatchProbability = myFact.MatchProbability()

		fmt.Println(myFact.Auth.DeviceID," Processed")
	}

	return nil
}


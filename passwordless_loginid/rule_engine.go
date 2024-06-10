package main

import (
	"github.com/hyperjumptech/grule-rule-engine/ast"
	"github.com/hyperjumptech/grule-rule-engine/builder"
	"github.com/hyperjumptech/grule-rule-engine/engine"
	"github.com/hyperjumptech/grule-rule-engine/pkg"
)

// ProcessDevices processes the devices using the rule engine
func ProcessDevice(devicesinfos *DeviceInfo, userPasskeyHistory []*UserPasskeyHistory, deviceFeatures DeviceFeature) (Output,error) {

	device := DeviceData{
		Auth: *devicesinfos,
		DeviceFeatures: deviceFeatures,
		UserPasskeyHistory: userPasskeyHistory,
	}
	helper1 := &IntermediateCloudPrediction{
		Prediction: "",
	}
	helper2 := &IntermediatePasskeyPrediction{
		Prediction: "",
	}
	
	knowledgeLibrary := ast.NewKnowledgeLibrary()
	ruleBuilder := builder.NewRuleBuilder(knowledgeLibrary)
	bs := pkg.NewBytesResource([]byte(rules))
	err := ruleBuilder.BuildRuleFromResource("PasswordLessRules", "0.0.2", bs)
	if err != nil {
		return Output{},err
	}

	knowledgeBase, _ := knowledgeLibrary.NewKnowledgeBaseInstance("PasswordLessRules", "0.0.2")
	ruleEngine := engine.NewGruleEngine()

	dataCtx := ast.NewDataContext()
	err = dataCtx.Add("DF", &device)
	if err != nil {
		return Output{},err
	}

	err = dataCtx.Add("IntermediateCloudPrediction", helper1)
	if err != nil {
		return Output{},err
	}
	err = dataCtx.Add("IntermediatePasskeyPrediction", helper2)
	if err != nil {
		return Output{},err
	}

	err = ruleEngine.Execute(dataCtx, knowledgeBase)
	if err != nil {
		return Output{},err
	}
	device.Output.MatchProbability = device.MatchProbability()

	return device.Output,nil
}


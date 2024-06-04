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
	
	
	var rules = `
	rule CheckPassKeysNonExistance "Check devices with no passkeys" salience 100 {
        when 
            !DF.IsPassKeyExisting()
        then
            DF.AssignOutPut("nopasskhey",1.0);
            Retract("CheckPassKeysNonExistance");
    }
	rule CheckCloudCompatibleVersion "Check compatible cloud browser" salience 50 {
        when 
            DF.MatchCloudCompatibleBrowser()
        then
            DF.AssignOutPut("dzadzdzdezdezdezdezdez",5.0);
            Retract("CheckCloudCompatibleVersion");
    }
	`

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
		

		dataCtx := ast.NewDataContext()
		err = dataCtx.Add("DF", myFact)
		if err != nil {
			panic(err)
		}
		

		err = engine.Execute(dataCtx, knowledgeBase)
		if err != nil {
			panic(err)
		}
		
		
	}

	err = writeDeviceInfoToCSV("device_info.csv", devices[:]) // Pass slice to avoid modifying original array
	if err != nil {
		panic(err)
	}

	fmt.Println("Device information saved to device_info.csv")
}
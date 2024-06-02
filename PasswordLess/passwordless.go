package main

import (
	"encoding/csv"
	"fmt"
	"os"
	"time"

	"github.com/hyperjumptech/grule-rule-engine/ast"
	"github.com/hyperjumptech/grule-rule-engine/builder"
	"github.com/hyperjumptech/grule-rule-engine/engine"
	"github.com/hyperjumptech/grule-rule-engine/pkg"
)

func writeDeviceInfoToCSV(fileName string, devices []*DeviceInfo) error {
	csvFile, err := os.OpenFile(fileName, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer csvFile.Close()

	csvWriter := csv.NewWriter(csvFile)

	fileInfo, err := csvFile.Stat()
	if err != nil {
		return err
	}

	if fileInfo.Size() == 0 {
		header := []string{"Label", "OSName", "Output"}
		err = csvWriter.Write(header)
		if err != nil {
			return err
		}
	}

	for _, device := range devices {
		record := []string{device.Label, device.OSName, device.OutPut}
		err = csvWriter.Write(record)
		if err != nil {
			return err
		}
	}

	csvWriter.Flush()
	return nil
}

type DeviceFeature struct {
	NavigLanguage string
	IsCMA         bool
	IsUVPPA       bool
	PasskeyType   string
}

type DeviceInfo struct {
	Label                       string
	CreatedAt                   time.Time
	OSName                      string
	OSVersion                   string
	BrowserName                 string
	BrowserVersion              string
	DeviceScreenSize_ColorDepth [2]float64
	IpAddress                   string
	DeviceFeature               DeviceFeature
	MatchOS                     bool
	MatchDeviceProperties       bool
	OutPut                      string
}

type IntermediatePasskeyPrediction struct {
	Prediction string
}

type IntermediateCloudPrediction struct {
	Prediction string
}

func (df *DeviceInfo) IsPassKeyExisting() bool {
	return df.DeviceFeature.IsUVPPA
}

func (df *DeviceInfo) IsMatchingDeviceProperties() bool {
	return df.MatchDeviceProperties
}

func (df *DeviceInfo) MatchOs() bool {
	return df.MatchOS
}

func (df *DeviceInfo) IsCloud() bool {
	return df.DeviceFeature.PasskeyType == "cloud"
}

func (df *DeviceInfo) AssignOutPut(res string) string {
	return res
}

func main() {

	myFacts := [4]*DeviceInfo{
		{
			Label:                       "Ezzoubeir",
			OSName:                      "Windows",
			OSVersion:                   "10",
			BrowserName:                 "Chrome",
			BrowserVersion:              "90.0",
			DeviceScreenSize_ColorDepth: [2]float64{1920, 1080},
			IpAddress:                   "192.168.1.1",
			DeviceFeature: DeviceFeature{
				NavigLanguage: "English",
				IsCMA:         true,
				IsUVPPA:       false,
				PasskeyType:   "local",
			},
			MatchOS:               true,
			MatchDeviceProperties: false,
		},
		{
			Label:                       "Anabelle",
			OSName:                      "macOS",
			OSVersion:                   "11.2",
			BrowserName:                 "Safari",
			BrowserVersion:              "14.0",
			DeviceScreenSize_ColorDepth: [2]float64{2560, 1600},
			IpAddress:                   "192.168.1.2",
			DeviceFeature: DeviceFeature{
				NavigLanguage: "French",
				IsCMA:         false,
				IsUVPPA:       true,
				PasskeyType:   "device",
			},
			MatchOS:               true,
			MatchDeviceProperties: true,
		},
		{
			Label:                       "Jasper",
				OSName:                      "Linux",
			OSVersion:                   "20.04",
			BrowserName:                 "Firefox",
			BrowserVersion:              "89.0",
			DeviceScreenSize_ColorDepth: [2]float64{1366, 768},
			IpAddress:                   "192.168.1.3",
			DeviceFeature: DeviceFeature{
				NavigLanguage: "Spanish",
				IsCMA:         true,
				IsUVPPA:       true,
				PasskeyType:   "cloud",
			},
			MatchOS:               false,
			MatchDeviceProperties: false,
		},
		{
			Label:                       "joe",
			OSName:                      "Android",
			OSVersion:                   "11",
			BrowserName:                 "Chrome",
			BrowserVersion:              "108",
			DeviceScreenSize_ColorDepth: [2]float64{1080, 1920},
			IpAddress:                   "192.168.1.4",
			DeviceFeature: DeviceFeature{
				NavigLanguage: "German",
				IsCMA:         true,
				IsUVPPA:       true,
				PasskeyType:   "cloud",
			},
			MatchOS:               false,
			MatchDeviceProperties: false,
		},
	}

	var rules = PasskeysNonExistanceRules + DevicePropertiesMatchRules + CloudBasedPasskeysRules

	drls := rules

	knowledgeLibrary := ast.NewKnowledgeLibrary()
	ruleBuilder := builder.NewRuleBuilder(knowledgeLibrary)
	bs := pkg.NewBytesResource([]byte(drls))
	err := ruleBuilder.BuildRuleFromResource("PasswordLessRules", "0.0.2", bs)
	if err != nil {
		panic(err)
	}

	knowledgeBase, _ := knowledgeLibrary.NewKnowledgeBaseInstance("PasswordLessRules", "0.0.2")
	engine := engine.NewGruleEngine()

	for _, myFact := range myFacts {
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
	}

	err = writeDeviceInfoToCSV("device_info.csv", myFacts[:]) // Pass slice to avoid modifying original array
	if err != nil {
		panic(err)
	}

	fmt.Println("Device information saved to device_info.csv")
}

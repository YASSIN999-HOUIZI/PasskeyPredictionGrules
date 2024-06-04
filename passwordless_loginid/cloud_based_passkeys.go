package main

var CloudBasedPasskeysRules = `
    rule CheckCloud "Check if the passkey is cloud-based" salience 70 {
        when 
            DF.IsPassKeyExisting() &&
            ((DF.MatchOs() && !DF.MatchDeviceProperties()) || !DF.MatchOs())
        then
            IntermediateCloudPrediction.Prediction = "Proceed";
            Retract("CheckCloud");
    }

    rule CheckIsNoIsCompatibleDevice "Check if the passkey is not cloud-based for Windows" salience 60 {
        when 
            DF.IsPassKeyExisting() &&
            IntermediateCloudPrediction.Prediction == "Proceed" &&
            !DF.IsCompatibleDevice()
        then
            DF.AssignOutPut("nopasskey",1.0);
            Retract("CheckIsNoIsCompatibleDevice");
    }

    rule CheckIsCompatibleDevice "Check if the passkey is not cloud-based for iOS" salience 60 {
		when 
			DF.IsPassKeyExisting() &&
			IntermediateCloudPrediction.Prediction == "Proceed" &&
			DF.IsCompatibleDevice()
		then
			IntermediatePasskeyPrediction.Prediction = "Proceed";
			Retract("CheckIsCompatibleDevice");
	}
	
	rule CloudPasskeyPredictionAssignNoPasskey "Assign no passkey if DeviceFeature IsCMA is false" salience 40 {
		when
			IntermediatePasskeyPrediction.Prediction == "Proceed" &&
			!DF.IsConditionalMediationAvailable() 
		then
			DF.AssignOutPut("nopasskey",1.0);
			Retract("CloudPasskeyPredictionAssignNoPasskey");
	}
	rule CloudPasskeyPredictionAssignCheckCloudIsNotCompatibleBrowsers "Assign CheckCloudCompatibleBrowsers if DeviceFeature IsCMA is true" salience 30 {
		when
			IntermediatePasskeyPrediction.Prediction == "Proceed" &&
			DF.IsConditionalMediationAvailable() &&
			!DF.MatchCloudCompatibleBrowser() 
		then
			DF.AssignOutPut("nopasskey",1.0);
			Retract("CloudPasskeyPredictionAssignCheckCloudIsNotCompatibleBrowsers");
	}

	rule CloudPasskeyPredictionAssignCheckCloudCompatibleBrowsers "Assign CheckCloudCompatibleBrowsers if DeviceFeature IsCMA is true" salience 30 {
		when
			IntermediatePasskeyPrediction.Prediction == "Proceed" &&
			DF.IsConditionalMediationAvailable() &&
			DF.MatchCloudCompatibleBrowser() 
		then
			DF.AssignOutPut("cloud",1.0);
			Retract("CloudPasskeyPredictionAssignCheckCloudCompatibleBrowsers");
	}

`

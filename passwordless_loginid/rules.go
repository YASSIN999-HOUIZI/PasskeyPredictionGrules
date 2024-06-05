package main

var rules = `
    rule CheckPassKeysNonExistance "Check devices with no passkeys" salience 100 {
        when 
            !DF.IsPassKeyExisting()
        then
            DF.AssignOutPut("nopasskhey",1.0);
            Retract("CheckPassKeysNonExistance");
    }
	rule CheckMatch_OS_DeviceProperties1 "Check match Device Properties" salience 90 {
        when 
            DF.IsPassKeyExisting() && DF.MatchOS() && DF.MatchDeviceProperties() 
        then
            DF.AssignOutPut(DF.GetDeviceOrLocal(),1.0);
            Retract("CheckMatch_OS_DeviceProperties1");
    }

    rule CheckMatch_OS_DeviceProperties2 "When device properties do not match" salience 80 {
        when 
            DF.IsPassKeyExisting() && DF.MatchOS() && !DF.MatchDeviceProperties()
        then
            IntermediateCloudPrediction.Prediction = "Proceed";
            Retract("CheckMatch_OS_DeviceProperties2");
    }

    rule CheckMatchOS "When the OS does not match" salience 80 {
        when 
            DF.IsPassKeyExisting() && !DF.MatchOS()
        then
            IntermediateCloudPrediction.Prediction = "Proceed";
            Retract("CheckMatchOS");
    }

	rule CheckCloud "Check if the passkey is cloud-based" salience 70 {
        when 
            DF.IsPassKeyExisting() &&
            ((DF.MatchOS() && !DF.MatchDeviceProperties()) || !DF.MatchOS())
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

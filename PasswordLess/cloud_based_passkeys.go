package main

var CloudBasedPasskeysRules = `
    rule CheckCloud "Check if the passkey is cloud-based" salience 70 {
        when 
            DF.IsPassKeyExisting() &&
            ((DF.MatchOs() && !DF.IsMatchingDeviceProperties()) || !DF.MatchOs()) && 
            DF.IsCloud() 
        then
            IntermediateCloudPrediction.Prediction = "Proceed";
            Retract("CheckCloud");
    }

	rule IsNotCloud "Check if the passkey is cloud-based" salience 65 {
        when 
            DF.IsPassKeyExisting() &&
            ((DF.MatchOs() && !DF.IsMatchingDeviceProperties()) || !DF.MatchOs()) && 
            !DF.IsCloud()
        then
			DF.OutPut = DF.AssignOutPut("NoPasskey");
            Retract("IsNotCloud");
    }

    rule CheckIsNotCloudForOtherOS "Check if the passkey is not cloud-based for Windows" salience 60 {
        when 
            DF.IsPassKeyExisting() &&
            IntermediateCloudPrediction.Prediction == "Proceed" &&
            DF.IsCloud() &&
            ( DF.OSName != "iOS" && DF.OSName != "macOS" && DF.OSName != "Android")
        then
            DF.OutPut = DF.AssignOutPut("NoPasskey");
            Retract("CheckIsNotCloudForOtherOS");
    }

    rule CheckIsNotCloudForiOS "Check if the passkey is not cloud-based for iOS" salience 60 {
		when 
			DF.IsPassKeyExisting() &&
			IntermediateCloudPrediction.Prediction == "Proceed" &&
			DF.IsCloud() &&
			DF.OSName == "iOS"
		then
			IntermediatePasskeyPrediction.Prediction = "Proceed";
			Retract("CheckIsNotCloudForiOS");
	}
	rule CheckIsNotCloudForMacOS "Check if the passkey is not cloud-based for macOS" salience 50 {
		when 
			DF.IsPassKeyExisting() &&
			IntermediateCloudPrediction.Prediction == "Proceed" &&
			DF.IsCloud() &&
			DF.OSName == "macOS"
		then
			IntermediatePasskeyPrediction.Prediction = "Proceed";
			Retract("CheckIsNotCloudForMacOS");
	}
	rule CheckIsNotCloudForAndroid "Check if the passkey is not cloud-based for Android" salience 50 {
		when 
			DF.IsPassKeyExisting() &&
			IntermediateCloudPrediction.Prediction == "Proceed" &&
			DF.IsCloud() &&
			DF.OSName == "Android"
		then
			IntermediatePasskeyPrediction.Prediction = "Proceed";
			Retract("CheckIsNotCloudForAndroid");
	}
	rule CloudPasskeyPredictionAssignNoPasskey "Assign no passkey if DeviceFeature IsCMA is false" salience 40 {
		when
			DF.IsCloud() &&
			DF.DeviceFeature.IsCMA == false &&
			IntermediatePasskeyPrediction.Prediction == "Proceed"
		then
			DF.OutPut = DF.AssignOutPut("NoPasskey");
			Retract("CloudPasskeyPredictionAssignNoPasskey");
	}
	rule CloudPasskeyPredictionAssignCheckCloudCompatibleBrowsers "Assign CheckCloudCompatibleBrowsers if DeviceFeature IsCMA is true" salience 30 {
		when
			DF.IsCloud() &&
			DF.DeviceFeature.IsCMA == true &&
			(DF.OSName == "Android" || DF.OSName == "iOS" || DF.OSName == "macOS") &&
			IntermediatePasskeyPrediction.Prediction == "Proceed"
		then
			IntermediatePasskeyPrediction.Prediction = "CheckCloudCompatibleBrowsers";
			Retract("CloudPasskeyPredictionAssignCheckCloudCompatibleBrowsers");
	}
	rule CheckCloudCompatibleBrowsersForiOS "Check if the browser is cloud-compatible for iOS" salience 20 {
		when
			DF.IsCloud() &&
			DF.DeviceFeature.IsCMA == true &&
			IntermediatePasskeyPrediction.Prediction == "CheckCloudCompatibleBrowsers" &&
			DF.OSName == "iOS" &&
			(DF.BrowserName == "Safari" || DF.BrowserName == "Chrome" || DF.BrowserName == "Edge" || DF.BrowserName == "Firefox") &&
			DF.BrowserVersion == "16"
		then
			DF.OutPut = DF.AssignOutPut("HasPasskey (cloud-based)");
			Retract("CheckCloudCompatibleBrowsersForiOS");
	}
	rule CheckCloudCompatibleBrowsersForMacOS "Check if the browser is cloud-compatible for macOS" salience 20 {
		when
			DF.IsCloud() &&
			DF.DeviceFeature.IsCMA == true &&
			IntermediatePasskeyPrediction.Prediction == "CheckCloudCompatibleBrowsers" &&
			DF.OSName == "macOS" &&
			((DF.BrowserName == "Safari" && DF.BrowserVersion == "16.1") ||
			(DF.BrowserName == "Chrome" && DF.BrowserVersion == "118") ||
			((DF.BrowserName == "Edge" || DF.BrowserName == "Firefox") && DF.BrowserVersion == "122"))
		then
			DF.OutPut = DF.AssignOutPut("HasPasskey (cloud-based)");
			Retract("CheckCloudCompatibleBrowsersForMacOS");
	}
	rule CheckCloudCompatibleBrowsersForAndroid "Check if the browser is cloud-compatible for Android" salience 20 {
		when
			DF.IsCloud() &&
			DF.DeviceFeature.IsCMA == true &&
			IntermediatePasskeyPrediction.Prediction == "CheckCloudCompatibleBrowsers" &&
			DF.OSName == "Android" &&
			((DF.BrowserName == "Chrome" && DF.BrowserVersion == "108") ||
			(DF.BrowserName == "Edge" && DF.BrowserVersion == "122") ||
			(DF.BrowserName == "Samsung" && DF.BrowserVersion == "21"))
		then
			DF.OutPut = DF.AssignOutPut("HasPasskey (cloud-based)");
			Retract("CheckCloudCompatibleBrowsersForAndroid");
	}
	rule CheckCloudCompatibleBrowsersFallback "Fallback for cloud-compatible browser check" salience 10 {
		when
			DF.IsCloud() &&
			DF.DeviceFeature.IsCMA == true &&
			IntermediatePasskeyPrediction.Prediction == "CheckCloudCompatibleBrowsers" &&
			!(
				(DF.OSName == "iOS" && (DF.BrowserName == "Safari" || DF.BrowserName == "Chrome" || DF.BrowserName == "Edge" || DF.BrowserName == "Firefox") && DF.BrowserVersion == "16") ||
				(DF.OSName == "macOS" && ((DF.BrowserName == "Safari" && DF.BrowserVersion == "16.1") || (DF.BrowserName == "Chrome" && DF.BrowserVersion == "118") || ((DF.BrowserName == "Edge" || DF.BrowserName == "Firefox") && DF.BrowserVersion == "122"))) ||
				(DF.OSName == "Android" && ((DF.BrowserName == "Chrome" && DF.BrowserVersion == "108") || (DF.BrowserName == "Edge" && DF.BrowserVersion == "122") || (DF.BrowserName == "Samsung" && DF.BrowserVersion == "21")))
			)
		then
			DF.OutPut = DF.AssignOutPut("NoPasskey");
			Retract("CheckCloudCompatibleBrowsersFallback");
	}
`

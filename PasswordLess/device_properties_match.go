package main

var DevicePropertiesMatchRules = `
    rule CheckMatch_OS_DeviceProperties1 "Check match Device Properties" salience 90 {
        when 
            DF.IsPassKeyExisting() && DF.MatchOs() && DF.IsMatchingDeviceProperties()  
        then
            DF.OutPut = DF.AssignOutPut("HasPassKey(device/local)");
            Retract("CheckMatch_OS_DeviceProperties1");
    }

    rule CheckMatch_OS_DeviceProperties2 "When device properties do not match" salience 80 {
        when 
            DF.IsPassKeyExisting() && DF.MatchOs() && !DF.IsMatchingDeviceProperties()
        then
            IntermediateCloudPrediction.Prediction = "Proceed";
            Retract("CheckMatch_OS_DeviceProperties2");
    }

    rule CheckMatchOS "When the OS does not match" salience 80 {
        when 
            DF.IsPassKeyExisting() && !DF.MatchOs()
        then
            IntermediateCloudPrediction.Prediction = "Proceed";
            Retract("CheckMatchOS");
    }
`

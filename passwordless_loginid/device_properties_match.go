package main

var DevicePropertiesMatchRules = `
    rule CheckMatch_OS_DeviceProperties1 "Check match Device Properties" salience 90 {
        when 
            DF.IsPassKeyExisting() && DF.MatchOS() && DF.MatchDeviceProperties()  
        then
            DF.AssignOutPut("device",1.0);
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
`

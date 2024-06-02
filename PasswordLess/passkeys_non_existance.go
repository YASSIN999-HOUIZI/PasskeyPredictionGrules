package main

var PasskeysNonExistanceRules = `
    rule CheckPassKeysNonExistance "Check devices with no passkeys" salience 100 {
        when 
            !DF.IsPassKeyExisting()
        then
            DF.OutPut = DF.AssignOutPut("NoPasskey");
            Retract("CheckPassKeysNonExistance");
    }
`

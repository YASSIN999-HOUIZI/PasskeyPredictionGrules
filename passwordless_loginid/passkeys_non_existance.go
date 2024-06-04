package main

var PasskeysNonExistanceRules = `
    rule CheckPassKeysNonExistance "Check devices with no passkeys" salience 100 {
        when 
            !DF.IsPassKeyExisting()
        then
            DF.AssignOutPut("nopasskhey",1.0);
            Retract("CheckPassKeysNonExistance");
    }
`

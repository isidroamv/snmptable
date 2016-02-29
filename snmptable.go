package main

import (
	"fmt"
	"time"
	"encoding/json"
	"encoding/hex"
	wsnmp "github.com/tiebingzhang/WapSNMP"
)


type MIBObject struct {
	oid string
	name string
	value string
	tableName string
	dataType wsnmp.BERType
}


func main(){
	
	target := ""
	community := ""
	version := wsnmp.SNMPv2c

	fmt.Printf("Contacting %v %v %v\n", target, community, version)
	snmp, err := wsnmp.NewWapSNMP(target, community, version, 2*time.Second, 5)
	if err != nil {
		fmt.Printf("Error creating wsnmp => %v\n", snmp)
		return
	}

	mibObjects := [] MIBObject {
		MIBObject {
			oid: ".1.3.6.1.4.1.14179.2.2.13.1.1",
			name: "bsnAPIfLoadRxUtilization",
			dataType: 3},
		MIBObject {
			oid: ".1.3.6.1.4.1.14179.2.2.1.1.30",
			name: "bsnAPGroupVlanName",
			dataType: 3},
		MIBObject {
			oid: ".1.3.6.1.4.1.14179.2.2.1.1.3",
			name: "bsnAPName",
			dataType: 3},
		MIBObject {
			oid: ".1.3.6.1.4.1.14179.2.2.1.1.24",
			name: "bsnAPTertiaryMwarName",
			dataType: 3},
	}

	mibObjectKey := MIBObject {
		oid: ".1.3.6.1.4.1.14179.2.2.1.1.1",
		name: "bsnAPDot3MacAddress",
		dataType: 4,
	}
	
	ta := getMIBTable(mibObjectKey, mibObjects, snmp)
	jsonString, err := json.Marshal(ta)
	fmt.Println("json error", err)
	fmt.Printf("json success: %s",jsonString)
}

/**
* Return a map with multiples fields
*
*/
func getMIBTable (mibKey MIBObject, mibs []MIBObject, snmp *wsnmp.WapSNMP) map[string] map[string] string {
	var snmpTable = make(map[string] map[string] string)
	var keyObject string

	oid, err := wsnmp.ParseOid(mibKey.oid)
	table, err := snmp.GetTable(oid)
	for key, value := range table {
		if str, ok := value.(string); ok {
			keyObject = getKeyObject(mibKey.oid, key)
		    snmpTable[keyObject] = make(map[string]string)
		    snmpTable[keyObject][mibKey.name] = typeConvert(str, mibKey.dataType);
		}
	}

	for _, mib := range mibs {
		fmt.Println("")
		oid, err = wsnmp.ParseOid(mib.oid)
		if (err!=nil) {
			fmt.Println("Can't parse Oid", err)
		}

		table, err = snmp.GetTable(oid)
		if (err!=nil) {
			fmt.Println("Can't get table", err)
		}

		for oid, value := range table {
			if value, ok := value.(string); ok {
				keyObject = getKeyObject(mib.oid, oid)
				if (snmpTable[keyObject]!=nil) {
					snmpTable[keyObject][mib.name] = typeConvert(value, mib.dataType) 
				}
			}
		}
	}
	return snmpTable
}

func getKeyObject(key string, oid string) string{
	return oid[len(key):]
}

func typeConvert(str string, dataType wsnmp.BERType) string {
	switch(dataType) {
	case wsnmp.AsnBitStr:
		return str
	case wsnmp.AsnOctetStr:
		bytes := []byte(str)
		return hex.EncodeToString(bytes)
	}
	return str
}
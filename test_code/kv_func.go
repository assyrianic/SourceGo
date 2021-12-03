package main

import (
	"sourcemod"
)

func KeyValuesToStringMap(kv KeyValues, stringmap map[string][]char, hide_top bool, depth int, prefix *[]char) {
	type SectStr = [128]char
	for {
		var section_name SectStr
		kv.GetSectionName(section_name, len(section_name))
		if kv.GotoFirstSubKey(false) {
			var new_prefix SectStr
			switch {
				case depth==0 && hide_top:
					new_prefix = ""
				case prefix[0] == 0:
					new_prefix = section_name
				default:
					FormatEx(new_prefix, len(new_prefix), "%s.%s", prefix, section_name)
			}
			KeyValuesToStringMap(kv, stringmap, hide_top, depth+1, new_prefix)
            kv.GoBack()
		} else {
			if kv.GetDataType(NULL_STRING) != KvData_None {
				var key SectStr
				if prefix[0] == 0 {
					key = section_name
				} else {
					FormatEx(key, len(key), "%s.%s", prefix, section_name)
				}
				
				// lowercaseify the key
				keylen := strlen(key)
				for i := 0; i < keylen; i++ {
					bytes := IsCharMB(key[i])
					if bytes==0 {
						key[i] = CharToLower(key[i])
					} else {
						i += (bytes - 1)
					}
				}
				
				var value SectStr
				kv.GetString(NULL_STRING, value, len(value), NULL_STRING)
				stringmap[key] = value
			}
		}
		
		if !kv.GotoNextKey(false) {
			break
		}
	}
}

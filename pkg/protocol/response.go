package protocol

import (
	"strconv"
)

type HamlibError string

func (e HamlibError) Error() string {
	message, ok := HamlibErrorMessages[string(e)]
	if ok {
		return message
	}
	return "unknown hamlib error: " + string(e)
}

const (
	InvalidParameter             HamlibError = "-1"
	InvalidConfiguration         HamlibError = "-2"
	MemoryShortage               HamlibError = "-3"
	FeatureNotImplemented        HamlibError = "-4"
	CommunicationTimedOut        HamlibError = "-5"
	IOError                      HamlibError = "-6"
	InternalHamlibError          HamlibError = "-7"
	ProtocolError                HamlibError = "-8"
	CommandRejectedByTheRig      HamlibError = "-9"
	CommandPerformedButTruncated HamlibError = "-10"
	FeatureNotAvailable          HamlibError = "-11"
	TargetVFOUnaccessible        HamlibError = "-12"
	CommunicationBusError        HamlibError = "-13"
	CommunicationBusCollision    HamlibError = "-14"
	NullRigHandle                HamlibError = "-15"
	InvalidVFO                   HamlibError = "-16"
	ArgumentOutOfDomain          HamlibError = "-17"
	FunctionDeprecated           HamlibError = "-18"
	SecurityError                HamlibError = "-19"
	RigNotPoweredOn              HamlibError = "-20"
)

var HamlibErrorMessages = map[string]string{
	"0":   "Command completed successfully",
	"-1":  "Invalid parameter",
	"-2":  "Invalid configuration",
	"-3":  "Memory shortage",
	"-4":  "Feature not implemented",
	"-5":  "Communication timed out",
	"-6":  "IO error",
	"-7":  "Internal Hamlib error",
	"-8":  "Protocol error",
	"-9":  "Command rejected by the rig",
	"-10": "Command performed, but arg truncated, result not guaranteed",
	"-11": "Feature not available",
	"-12": "Target VFO unaccessible",
	"-13": "Communication bus error",
	"-14": "Communication bus collision",
	"-15": "NULL RIG handle or invalid pointer parameter",
	"-16": "Invalid VFO",
	"-17": "Argument out of domain of func",
	"-18": "Function deprecated",
	"-19": "Security error password not provided or crypto failure",
	"-20": "Rig is not powered on",
}

func OKResponse(cmd CommandKey) Response {
	return Response{Command: cmd, Result: "0"}
}

func ErrorResponse(cmd CommandKey, err HamlibError) Response {
	return Response{Command: cmd, Result: string(err)}
}

func GetFreqResponse(frequency int) Response {
	return Response{
		Command: "get_freq",
		Data:    []string{strconv.Itoa(frequency)},
		Keys:    []string{"Frequency"},
		Result:  "0",
	}
}

func GetVFOResponse(vfo string) Response {
	return Response{
		Command: "get_vfo",
		Data:    []string{vfo},
		Keys:    []string{"VFO"},
		Result:  "0",
	}
}

func GetModeResponse(mode string, passband int) Response {
	return Response{
		Command: "get_mode",
		Data:    []string{mode, strconv.Itoa(passband)},
		Keys:    []string{"Mode", "Passband"},
		Result:  "0",
	}
}

func GetSplitFreqResponse(frequency int) Response {
	return Response{
		Command: "get_split_freq",
		Data:    []string{strconv.Itoa(frequency)},
		Keys:    []string{"TX Frequency"},
		Result:  "0",
	}
}

func GetSplitVFOResponse(enabled bool, txVFO string) Response {
	splitEnabled := "0"
	if enabled {
		splitEnabled = "1"
	}
	return Response{
		Command: "get_split_vfo",
		Data:    []string{splitEnabled, txVFO},
		Keys:    []string{"Split", "TX VFO"},
		Result:  "0",
	}
}

func GetSplitModeResponse(mode string, passband int) Response {
	return Response{
		Command: "get_split_mode",
		Data:    []string{mode, strconv.Itoa(passband)},
		Keys:    []string{"TX Mode", "TX Passband"},
		Result:  "0",
	}
}

func GetPTTResponse(enabled bool) Response {
	pttEnabled := "0"
	if enabled {
		pttEnabled = "1"
	}
	return Response{
		Command: "get_ptt",
		Data:    []string{pttEnabled},
		Keys:    []string{"PTT"},
		Result:  "0",
	}
}

func GetLevelKeyspdResponse(wpm int) Response {
	return Response{
		Command: "get_level_keyspd",
		Data:    []string{strconv.Itoa(wpm)},
		Keys:    []string{"KEYSPD"},
		Result:  "0",
	}
}

func GetLockModeResponse(enabled bool) Response {
	lockModeEnabled := "0"
	if enabled {
		lockModeEnabled = "1"
	}
	return Response{
		Command: "get_lock_mode",
		Data:    []string{lockModeEnabled},
		Keys:    []string{"Locked"},
		Result:  "0",
	}
}

var NoResponse = Response{}

var ChkVFOResponse = Response{
	Command: "chk_vfo",
	Data:    []string{"0"},
	Keys:    []string{""},
	Result:  "0",
}

var DumpStateResponse = Response{
	Command: "dump_state",
	Data: []string{`0
1
2
150000.000000 1500000000.000000 0x1ff -1 -1 0x10000003 0x3
0 0 0 0 0 0 0
0 0 0 0 0 0 0
0x1ff 1
0x1ff 0
0 0
0x1e 2400
0x2 500
0x1 8000
0x1 2400
0x20 15000
0x20 8000
0x40 230000
0 0
9990
9990
10000
0
10 
10 20 30 
0xffffffffffffffff
0xffffffffffffffff
0xfffffffff7ffffff
0xffffffff83ffffff
0xffffffffffffffff
0xffffffffffffffbf
`},
	Result: "0",
}

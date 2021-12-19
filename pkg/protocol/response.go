package protocol

import "strconv"

func OKResponse(cmd CommandKey) Response {
	return Response{Command: cmd, Result: "0"}
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

package protocol

var (
	ShortCommands = make(map[string]Command)
	LongCommands  = make(map[string]Command)
	Commands      = []Command{
		{
			Short:                "F",
			Long:                 "set_freq",
			InvalidatesCommand:   "get_freq",
			SupportsExtendedMode: true,
		},
		{
			Short:                "f",
			Long:                 "get_freq",
			SupportsExtendedMode: true,
			Cacheable:            true,
		},
		{
			Short:              "M",
			Long:               "set_mode",
			InvalidatesCommand: "get_mode",
		},
		{
			Short:     "m",
			Long:      "get_mode",
			Cacheable: true,
		},
		{
			Short:              "V",
			Long:               "set_vfo",
			InvalidatesCommand: "get_vfo",
		},
		{
			Short:     "v",
			Long:      "get_vfo",
			Cacheable: true,
		},
		{
			Short:              "J",
			Long:               "set_rit",
			InvalidatesCommand: "get_rit",
		},
		{
			Short:     "j",
			Long:      "get_rit",
			Cacheable: true,
		},
		{
			Short:              "Z",
			Long:               "set_xit",
			InvalidatesCommand: "get_xit",
		},
		{
			Short:     "z",
			Long:      "get_xit",
			Cacheable: true,
		},
		{
			Short:              "T",
			Long:               "set_ptt",
			InvalidatesCommand: "get_ptt",
		},
		{
			Short:     "t",
			Long:      "get_ptt",
			Cacheable: true,
		},
		{
			Short:     "\x8b",
			Long:      "get_dcd",
			Cacheable: true,
		},
		{
			Short:              "R",
			Long:               "set_rprt_shift",
			InvalidatesCommand: "get_rprt_shift",
		},
		{
			Short:     "r",
			Long:      "get_rprt_shift",
			Cacheable: true,
		},
		{
			Short:              "O",
			Long:               "set_rprt_offs",
			InvalidatesCommand: "get_rprt_offs",
		},
		{
			Short:     "o",
			Long:      "get_rprt_offset",
			Cacheable: true,
		},
		{
			Short:              "C",
			Long:               "set_ctcss_tone",
			InvalidatesCommand: "get_ctcss_tone",
		},
		{
			Short:     "c",
			Long:      "get_ctcss_tone",
			Cacheable: true,
		},
		{
			Short:              "D",
			Long:               "set_dcs_code",
			InvalidatesCommand: "get_dcs_code",
		},
		{
			Short:     "d",
			Long:      "get_dcs_code",
			Cacheable: true,
		},
		{
			Short:              "\x90",
			Long:               "set_ctcss_sql",
			InvalidatesCommand: "get_ctcss_sql",
		},
		{
			Short:     "\x91",
			Long:      "get_ctcss_sql",
			Cacheable: true,
		},
		{
			Short:              "\x92",
			Long:               "set_dcs_sql",
			InvalidatesCommand: "get_dcs_sql",
		},
		{
			Short:     "\x93",
			Long:      "get_dcs_sql",
			Cacheable: true,
		},
		{
			Short:              "I",
			Long:               "set_split_freq",
			InvalidatesCommand: "get_split_freq",
		},
		{
			Short:     "i",
			Long:      "get_split_freq",
			Cacheable: true,
		},
		{
			Short:              "X",
			Long:               "set_split_mode",
			InvalidatesCommand: "get_split_mode",
		},
		{
			Short:     "x",
			Long:      "get_split_mode",
			Cacheable: true,
		},
		{
			Short:              "K",
			Long:               "set_split_freq_mode",
			InvalidatesCommand: "get_split_freq_mode",
		},
		{
			Short:     "k",
			Long:      "get_split_freq_mode",
			Cacheable: true,
		},
		{
			Short:              "S",
			Long:               "set_split_vfo",
			InvalidatesCommand: "get_split_vfo",
		},
		{
			Short:     "s",
			Long:      "get_split_vfo",
			Cacheable: true,
		},
		{
			Short:              "N",
			Long:               "set_ts",
			InvalidatesCommand: "get_ts",
		},
		{
			Short:     "n",
			Long:      "get_ts",
			Cacheable: true,
		},
		{
			Short:              "U",
			Long:               "set_func",
			InvalidatesCommand: "get_func",
			HasSubCommand:      true,
		},
		{
			Short:         "u",
			Long:          "get_func",
			HasSubCommand: true,
			Cacheable:     true,
		},
		{
			Short:              "L",
			Long:               "set_level",
			InvalidatesCommand: "get_level",
			HasSubCommand:      true,
		},
		{
			Short:         "l",
			Long:          "get_level",
			HasSubCommand: true,
			Cacheable:     true,
		},
		{
			Short:              "P",
			Long:               "set_parm",
			InvalidatesCommand: "get_parm",
			HasSubCommand:      true,
		},
		{
			Short:         "p",
			Long:          "get_parm",
			HasSubCommand: true,
			Cacheable:     true,
		},
		{
			Short: "B",
			Long:  "set_bank",
		},
		{
			Short:              "E",
			Long:               "set_mem",
			InvalidatesCommand: "get_mem",
		},
		{
			Short:     "e",
			Long:      "get_mem",
			Cacheable: true,
		},
		{
			Short: "G",
			Long:  "vfo_op",
		},
		{
			Short: "g",
			Long:  "scan",
		},
		{
			Short:              "H",
			Long:               "set_channel",
			InvalidatesCommand: "get_channel",
		},
		{
			Short:     "h",
			Long:      "get_channel",
			Cacheable: true,
		},
		{
			Short:              "A",
			Long:               "set_trn",
			InvalidatesCommand: "get_trn",
		},
		{
			Short:     "a",
			Long:      "get_trn",
			Cacheable: true,
		},
		{
			Short:              "Y",
			Long:               "set_ant",
			InvalidatesCommand: "get_ant",
		},
		{
			Short:     "y",
			Long:      "get_ant",
			Cacheable: true,
		},
		{
			Short: "*",
			Long:  "reset",
		},
		{
			Short:              "\x87",
			Long:               "set_powerstat",
			InvalidatesCommand: "get_powerstat",
		},
		{
			Short:     "\x88",
			Long:      "get_powerstat",
			Cacheable: true,
		},
		{
			Short: "\x89",
			Long:  "send_dtmf",
		},
		{
			Short: "\x8a",
			Long:  "recv_dtmf",
		},
		{
			Short: "b",
			Long:  "send_morse",
		},
		{
			Short: "w",
			Long:  "send_cmd",
		},
		{
			Short:     "_",
			Long:      "get_info",
			Cacheable: true,
		},
		{
			Short:     "1",
			Long:      "dump_caps",
			Cacheable: true,
		},
		{
			Short: "2",
			Long:  "power2mW",
		},
		{
			Short: "3",
			Long:  "mW2power",
		},
		{
			Short: "\x8f",
			Long:  "dump_state",
		},
		{
			Short:     "\xf0",
			Long:      "chk_vfo",
			Cacheable: true,
		},
		{
			Short: "\xf1",
			Long:  "halt",
		},
		{
			Short: "\x8c",
			Long:  "pause",
		},
	}
)

func init() {
	for _, cmd := range Commands {
		ShortCommands[cmd.Short] = cmd
		LongCommands[cmd.Long] = cmd
	}
}

func ShortCommand(s string) Command {
	cmd, ok := ShortCommands[s]
	if !ok {
		panic("unknown command " + s)
	}
	return cmd
}

func LongCommand(s string) Command {
	cmd, ok := LongCommands[s]
	if !ok {
		panic("unknown command " + s)
	}
	return cmd
}

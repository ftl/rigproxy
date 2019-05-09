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
			Cacheable:            true,
			SupportsExtendedMode: true,
		},
		{
			Short:                "M",
			Long:                 "set_mode",
			InvalidatesCommand:   "get_mode",
			SupportsExtendedMode: true,
		},
		{
			Short:                "m",
			Long:                 "get_mode",
			Cacheable:            true,
			SupportsExtendedMode: true,
		},
		{
			Short:                "V",
			Long:                 "set_vfo",
			InvalidatesCommand:   "get_vfo",
			SupportsExtendedMode: true,
		},
		{
			Short:                "v",
			Long:                 "get_vfo",
			Cacheable:            true,
			SupportsExtendedMode: true,
		},
		{
			Short:                "J",
			Long:                 "set_rit",
			InvalidatesCommand:   "get_rit",
			SupportsExtendedMode: true,
		},
		{
			Short:                "j",
			Long:                 "get_rit",
			Cacheable:            true,
			SupportsExtendedMode: true,
		},
		{
			Short:                "Z",
			Long:                 "set_xit",
			InvalidatesCommand:   "get_xit",
			SupportsExtendedMode: true,
		},
		{
			Short:                "z",
			Long:                 "get_xit",
			Cacheable:            true,
			SupportsExtendedMode: true,
		},
		{
			Short:                "T",
			Long:                 "set_ptt",
			InvalidatesCommand:   "get_ptt",
			SupportsExtendedMode: true,
		},
		{
			Short:                "t",
			Long:                 "get_ptt",
			Cacheable:            true,
			SupportsExtendedMode: true,
		},
		{
			Short:                "\x8b",
			Long:                 "get_dcd",
			Cacheable:            true,
			SupportsExtendedMode: true,
		},
		{
			Short:                "R",
			Long:                 "set_rptr_shift",
			InvalidatesCommand:   "get_rptr_shift",
			SupportsExtendedMode: true,
		},
		{
			Short:                "r",
			Long:                 "get_rptr_shift",
			Cacheable:            true,
			SupportsExtendedMode: true,
		},
		{
			Short:                "O",
			Long:                 "set_rptr_offs",
			InvalidatesCommand:   "get_rptr_offs",
			SupportsExtendedMode: true,
		},
		{
			Short:                "o",
			Long:                 "get_rptr_offs",
			Cacheable:            true,
			SupportsExtendedMode: true,
		},
		{
			Short:                "C",
			Long:                 "set_ctcss_tone",
			InvalidatesCommand:   "get_ctcss_tone",
			SupportsExtendedMode: true,
		},
		{
			Short:                "c",
			Long:                 "get_ctcss_tone",
			Cacheable:            true,
			SupportsExtendedMode: true,
		},
		{
			Short:                "D",
			Long:                 "set_dcs_code",
			InvalidatesCommand:   "get_dcs_code",
			SupportsExtendedMode: true,
		},
		{
			Short:                "d",
			Long:                 "get_dcs_code",
			Cacheable:            true,
			SupportsExtendedMode: true,
		},
		{
			Short:                "\x90",
			Long:                 "set_ctcss_sql",
			InvalidatesCommand:   "get_ctcss_sql",
			SupportsExtendedMode: true,
		},
		{
			Short:                "\x91",
			Long:                 "get_ctcss_sql",
			Cacheable:            true,
			SupportsExtendedMode: true,
		},
		{
			Short:                "\x92",
			Long:                 "set_dcs_sql",
			InvalidatesCommand:   "get_dcs_sql",
			SupportsExtendedMode: true,
		},
		{
			Short:                "\x93",
			Long:                 "get_dcs_sql",
			Cacheable:            true,
			SupportsExtendedMode: true,
		},
		{
			Short:                "I",
			Long:                 "set_split_freq",
			InvalidatesCommand:   "get_split_freq",
			SupportsExtendedMode: true,
		},
		{
			Short:                "i",
			Long:                 "get_split_freq",
			Cacheable:            true,
			SupportsExtendedMode: true,
		},
		{
			Short:                "X",
			Long:                 "set_split_mode",
			InvalidatesCommand:   "get_split_mode",
			SupportsExtendedMode: true,
		},
		{
			Short:                "x",
			Long:                 "get_split_mode",
			Cacheable:            true,
			SupportsExtendedMode: true,
		},
		{
			Short:                "K",
			Long:                 "set_split_freq_mode",
			InvalidatesCommand:   "get_split_freq_mode",
			SupportsExtendedMode: true,
		},
		{
			Short:                "k",
			Long:                 "get_split_freq_mode",
			Cacheable:            true,
			SupportsExtendedMode: true,
		},
		{
			Short:                "S",
			Long:                 "set_split_vfo",
			InvalidatesCommand:   "get_split_vfo",
			SupportsExtendedMode: true,
		},
		{
			Short:                "s",
			Long:                 "get_split_vfo",
			Cacheable:            true,
			SupportsExtendedMode: true,
		},
		{
			Short:                "N",
			Long:                 "set_ts",
			InvalidatesCommand:   "get_ts",
			SupportsExtendedMode: true,
		},
		{
			Short:                "n",
			Long:                 "get_ts",
			Cacheable:            true,
			SupportsExtendedMode: true,
		},
		{
			Short:                "U",
			Long:                 "set_func",
			InvalidatesCommand:   "get_func",
			HasSubCommand:        true,
			SupportsExtendedMode: true,
		},
		{
			Short:                "u",
			Long:                 "get_func",
			HasSubCommand:        true,
			Cacheable:            true,
			SupportsExtendedMode: true,
		},
		{
			Short:                "L",
			Long:                 "set_level",
			InvalidatesCommand:   "get_level",
			HasSubCommand:        true,
			SupportsExtendedMode: true,
		},
		{
			Short:                "l",
			Long:                 "get_level",
			HasSubCommand:        true,
			Cacheable:            true,
			SupportsExtendedMode: true,
		},
		{
			Short:                "P",
			Long:                 "set_parm",
			InvalidatesCommand:   "get_parm",
			HasSubCommand:        true,
			SupportsExtendedMode: true,
		},
		{
			Short:                "p",
			Long:                 "get_parm",
			HasSubCommand:        true,
			Cacheable:            true,
			SupportsExtendedMode: true,
		},
		{
			Short:                "B",
			Long:                 "set_bank",
			SupportsExtendedMode: true,
		},
		{
			Short:                "E",
			Long:                 "set_mem",
			InvalidatesCommand:   "get_mem",
			SupportsExtendedMode: true,
		},
		{
			Short:                "e",
			Long:                 "get_mem",
			Cacheable:            true,
			SupportsExtendedMode: true,
		},
		{
			Short:                "G",
			Long:                 "vfo_op",
			SupportsExtendedMode: true,
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
			Short:                "A",
			Long:                 "set_trn",
			InvalidatesCommand:   "get_trn",
			SupportsExtendedMode: true,
		},
		{
			Short:                "a",
			Long:                 "get_trn",
			Cacheable:            true,
			SupportsExtendedMode: true,
		},
		{
			Short:                "Y",
			Long:                 "set_ant",
			InvalidatesCommand:   "get_ant",
			SupportsExtendedMode: true,
		},
		{
			Short:                "y",
			Long:                 "get_ant",
			Cacheable:            true,
			SupportsExtendedMode: true,
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
			Short:                "1",
			Long:                 "dump_caps",
			Cacheable:            true,
			SupportsExtendedMode: true,
		},
		{
			Short:                "3",
			Long:                 "dump_conf",
			Cacheable:            true,
			SupportsExtendedMode: true,
		},
		{
			Short:                "2",
			Long:                 "power2mW",
			SupportsExtendedMode: true,
		},
		{
			Short:                "4",
			Long:                 "mW2power",
			SupportsExtendedMode: true,
		},
		{
			Short:                "\x8f",
			Long:                 "dump_state",
			SupportsExtendedMode: true,
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
			Short:                "\x8c",
			Long:                 "pause",
			SupportsExtendedMode: true,
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

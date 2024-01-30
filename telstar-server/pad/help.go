package pad

// FIXME returning to pad should redisplay welcome message with cls.
// FIXME after several CR's entered the pad should redisplay welcome message with cls.

const (
	RETURN = ""
	HELP   = "\fTELSTAR VPAD HELP: \r\n\r\n" + // \f os HOME+CLS i.e. ctrl L
		//123456789012345678901234567890123456789
		" HOSTS        Displays the shortcode\r\n" +
		"              directory of hosts.\r\n" +
		" CALL <host>  Connect to the specified\r\n" +
		"              host.\r\n" +
		"              or <ipaddress>:<port>.\r\n" +
		" PARS         Displays parameters for\r\n" +
		"              the current profile.\r\n" +
		" PROFILE      Sets/displays the current\r\n" +
		"              profile.\r\n" +
		" PROFILES     Displays the available \r\n" +
		"              profiles.\r\n" +
		" HELP         Displays this file.\r\n" +
		" HELP <topic> Displays additional help.\r\n" +
		"              e.g. HELP CALL.\r\n" +
		RETURN

	HELP_PARS = "\fHELP PARAMETERS:\r\n\r\n" +
		" Parameters within the current profile\r\n" +
		" can be changed using the command:\r\n\r\n" +
		" e.g. to set echo on\r\n\r\n" +
		"    P2=1\r\n" +
		RETURN

	HELPPROFILE = "\fHELP PROFILE:\r\n\r\n" +
		" Change the current profile with\r\n\r\n" +
		" e.g. set the current profile to P8.\r\n\r\n" +
		"    PROFILE=P8 \r\n" +
		RETURN

	HELPCALL = "\fHELP CALL:\r\n\r\n" +
		" CALL is used to connect to a remote\r\n" +
		" service.\r\n\r\n" +
		" e.g\r\n\r\n" +
		"    CAll glasstty.com:6502\r\n\r\n" +
		" Alternatively CALL can be used to\r\n" +
		" connect to a remote service using a\r\n" +
		" directory entry (see HOSTS).\r\n\r\n" +
		" e.g.\r\n\r\n" +
		"    CALL Telstar\r\n" +
		RETURN
)

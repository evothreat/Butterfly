package main

type JobType int

const (
	UNKNOWN JobType = iota
	SHELL_CMD
	UPLOAD
	DOWNLOAD
	CHDIR
	SLEEP
	BOOST
	SCREENSHOT
	MSG
)

type JobTypeInfo struct {
	jtype JobType
	argsN int
	// maybe also add handler function
}

var jobTypesMap = map[string]*JobTypeInfo{
	"cmd": {
		jtype: SHELL_CMD,
		argsN: -1,
	},
	"upload": {
		jtype: UPLOAD,
		argsN: 1,
	},
	"download": {
		jtype: DOWNLOAD,
		argsN: 2,
	},
	"sleep": {
		jtype: SLEEP,
		argsN: 1,
	},
	"boost": {
		jtype: BOOST,
		argsN: 1,
	},
	"chdir": {
		jtype: CHDIR,
		argsN: 1,
	},
	"msg": {
		jtype: MSG,
		argsN: 2,
	},
	"shot": {
		jtype: SCREENSHOT,
		argsN: 0,
	},
}

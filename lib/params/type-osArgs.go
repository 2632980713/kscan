package params

import (
	"flag"
	"fmt"
	"kscan/lib/slog"
	"os"
	"regexp"
)

type OsArgs struct {
	help, debug, scanPing, check, spy, noColor        bool
	target, port, output, proxy, path, host, encoding string
	outputJson                                        string
	USAGE, HELP, LOGO, SYNTAX                         string
	top, threads, timeout, rarity                     int
	//hydra模块
	hydra, hydraUpdate             bool
	hydraUser, hydraPass, hydraMod string
	//fofa模块
	fofa, fofaField, fofaFixKeyword string
	fofaSize                        int
	fofaSyntax                      bool
	scan                            bool
	//参数校验正则
	intsReg, strsReg, proxyReg *regexp.Regexp
}

func (o OsArgs) NoColor() bool {
	return o.noColor
}

func (o OsArgs) Fofa() string {
	return o.fofa
}

func (o OsArgs) Scan() bool {
	return o.scan
}

func (o OsArgs) FofaField() string {
	return o.fofaField
}

func (o OsArgs) FofaFixKeyword() string {
	return o.fofaFixKeyword
}

func (o OsArgs) FofaSize() int {
	return o.fofaSize
}

func (o OsArgs) Target() string {
	return o.target
}

func (o OsArgs) Port() string {
	return o.port
}

func (o OsArgs) Output() string {
	return o.output
}

func (o OsArgs) OutputJson() string {
	return o.outputJson
}

func (o OsArgs) Proxy() string {
	return o.proxy
}

func (o OsArgs) Path() string {
	return o.path
}
func (o OsArgs) Host() string {
	return o.host
}

func (o OsArgs) Top() int {
	return o.top
}

func (o OsArgs) Threads() int {
	return o.threads
}

func (o OsArgs) Timeout() int {
	return o.timeout
}

func (o OsArgs) Rarity() int {
	return o.rarity
}

func (o OsArgs) ScanPing() bool {
	return o.scanPing
}

func (o OsArgs) Check() bool {
	return o.check
}

func (o OsArgs) Debug() bool {
	return o.debug
}

func (o OsArgs) Hydra() bool {
	return o.hydra
}

func (o OsArgs) HydraUpdate() bool {
	return o.hydraUpdate
}

func (o OsArgs) HydraMod() string {
	return o.hydraMod
}

func (o OsArgs) HydraUser() string {
	return o.hydraUser
}

func (o OsArgs) HydraPass() string {
	return o.hydraPass
}

func (o OsArgs) Encoding() string {
	return o.encoding
}

func (o OsArgs) Spy() bool {
	return o.spy
}

//初始化参数
func (o *OsArgs) LoadOsArgs() {
	//自定义Usage
	flag.Usage = func() {
		fmt.Print(o.LOGO)
	}
	flag.BoolVar(&o.help, "h", false, "")
	flag.BoolVar(&o.help, "help", false, "")
	flag.BoolVar(&o.debug, "debug", false, "")
	flag.BoolVar(&o.debug, "d", false, "")
	//spy模块
	flag.BoolVar(&o.spy, "spy", false, "")
	//hydra模块
	flag.BoolVar(&o.hydra, "hydra", false, "")
	flag.BoolVar(&o.hydraUpdate, "hydra-update", false, "")
	flag.StringVar(&o.hydraUser, "hydra-user", "", "")
	flag.StringVar(&o.hydraPass, "hydra-pass", "", "")
	flag.StringVar(&o.hydraMod, "hydra-mod", "", "")
	//fofa模块
	flag.StringVar(&o.fofa, "fofa", "", "")
	flag.StringVar(&o.fofa, "f", "", "")
	flag.StringVar(&o.fofaField, "fofa-field", "", "")
	flag.StringVar(&o.fofaFixKeyword, "fofa-fix-keyword", "", "")
	flag.IntVar(&o.fofaSize, "fofa-size", 100, "")
	flag.BoolVar(&o.fofaSyntax, "fofa-syntax", false, "")
	flag.BoolVar(&o.scan, "scan", false, "")
	//kscan模块
	flag.StringVar(&o.target, "target", "", "")
	flag.StringVar(&o.target, "t", "", "")
	flag.StringVar(&o.port, "p", "", "")
	flag.StringVar(&o.port, "port", "", "")
	flag.StringVar(&o.proxy, "proxy", "", "")
	flag.StringVar(&o.path, "path", "", "")
	flag.StringVar(&o.host, "host", "", "")
	flag.IntVar(&o.rarity, "rarity", 9, "")
	flag.IntVar(&o.top, "top", 400, "")
	flag.IntVar(&o.threads, "threads", 400, "")
	flag.IntVar(&o.timeout, "timeout", 3, "")
	flag.BoolVar(&o.scanPing, "Pn", false, "")
	flag.BoolVar(&o.check, "check", false, "")
	//输出模块
	flag.StringVar(&o.encoding, "encoding", "utf-8", "")
	flag.StringVar(&o.output, "o", "", "")
	flag.StringVar(&o.output, "output", "", "")
	flag.StringVar(&o.outputJson, "oJ", "", "")
	flag.BoolVar(&o.noColor, "Cn", false, "")
	flag.Parse()
}

//初始化函数
func (o *OsArgs) PrintBanner() {
	//不带参数则对应usage
	if len(os.Args) == 1 {
		slog.Data(o.LOGO)
		slog.Data(o.USAGE)
		os.Exit(0)
	}
	if o.help {
		slog.Data(o.LOGO)
		slog.Data(o.USAGE)
		slog.Data(o.HELP)
		os.Exit(0)
	}
	if o.fofaSyntax {
		slog.Data(o.LOGO)
		slog.Data(o.USAGE)
		slog.Data(o.SYNTAX)
		os.Exit(0)
	}
	//打印logo
	slog.Data(o.LOGO)
}

func (o *OsArgs) CheckArgs() {
	//判断必须的参数是否存在
	if o.spy == true {
		return
	}

	if o.target == "" && o.fofa == "" {
		slog.Error("至少有target、fofa两个参数中的一个")
	}

	//判断冲突参数
	if o.port != "" && o.top != 400 {
		slog.Error("port、top参数不能同时使用")
	}
	if o.target != "" && o.fofa != "" {
		slog.Error("target、fofa参数不能同时使用")
	}

	//判断内容
	if o.port != "" {
		if !o.intsReg.MatchString(o.port) {
			slog.Error("PORT参数输入错误,其格式应为80，8080，8081-8090")
		}
	}
	if o.top != 0 {
		if o.top > 1000 || o.top < 1 {
			slog.Error("TOP参数输入错误,TOP参数应为1-1000之间的整数。")
		}
	}
	if o.output != "" {
		//验证output参数
	}
	if o.outputJson != "" {
		//验证outputJson参数
	}
	if o.proxy != "" {
		if !o.proxyReg.MatchString(o.proxy) {
			slog.Error("PROXY参数输入错误，其格式应为：http://ip:port，支持socks5/4")
		}
	}
	if o.path != "" {
		if !o.strsReg.MatchString(o.path) {
			slog.Error("PATH参数输入错误，其格式应为：/asdfasdf，可使用逗号输入多个路径")
		}
	}
	if o.host != "" {
		//验证host参数
	}
	if o.threads != 0 {
		if o.threads > 2048 {
			slog.Error("Threads参数最大值为2048")
		}
		//验证threads参数
	}
	if o.timeout != 3 {
		//验证timeout参数
	}
}

func New(logo string, usage string, help string, syntax string) *OsArgs {
	return &OsArgs{
		LOGO:     logo,
		USAGE:    usage,
		HELP:     help,
		SYNTAX:   syntax,
		intsReg:  regexp.MustCompile("^((?:[0-9]+)(?:-[0-9]+)?)(?:,(?:[0-9]+)(?:-[0-9]+)?)*$"),
		strsReg:  regexp.MustCompile("^([\\.A-Za-z0-9/]+)(,[\\.A-Za-z0-9/])*$"),
		proxyReg: regexp.MustCompile("^(http|https|socks5|socks4)://[0-9.]+:[0-9]+$"),
	}
}

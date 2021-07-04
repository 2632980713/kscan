package app

import (
	"github.com/lcvvvv/urlparse"
	"kscan/lib/IP"
	"kscan/lib/chinese"
	"kscan/lib/misc"
	"kscan/lib/params"
	"kscan/lib/slog"
	"os"
	"regexp"
	"runtime"
	"strings"
	"time"
)

type Config struct {
	HostTarget, UrlTarget       []string
	Port                        []int
	Output                      *os.File
	Proxy, Host, Path, Encoding string
	OSEncoding, NewLine         string
	Threads                     int
	Timeout                     time.Duration
	HostTargetNum               int
	UrlTargetNum                int
	PortNum                     int
	ScanPing, Check, Spy        bool
	//FofaEmail, FofaKey    string
}

var Setting = New()

func (c *Config) WriteLine(s string) {
	if c.OSEncoding == "utf-8" {
		s = chinese.ToUTF8(s)
	} else {
		s = chinese.ToGBK(s)
	}
	s = s + c.NewLine
	_, _ = c.Output.WriteString(s)
}

func (c *Config) Load(p *params.OsArgs) {
	if p.Spy() {
		c.Spy = p.Spy()
		return
	}

	c.loadTarget(p.Target(), false)
	c.loadTargetNum()
	c.loadPort(p.Port())
	c.loadPort(p.Top())
	c.loadPortNum()
	c.loadOutput(p.Output())
	c.loadTimeout(p.Timeout())
	c.ScanPing = p.ScanPing()
	c.Check = p.Check()
	c.Path = p.Path()
	c.Proxy = p.Proxy()
	c.Host = p.Host()
	c.Threads = p.Threads()
	c.Encoding = p.Encoding()
}

func (c *Config) loadTarget(expr string, recursion bool) {
	//判断target字符串是否为文件
	if regexp.MustCompile("^file:.+$").MatchString(expr) {
		expr = strings.Replace(expr, "file:", "", 1)
		err := misc.ReadLine(expr, c.loadTarget)
		if err != nil {
			if recursion == true {
				slog.Debug(expr + err.Error())
			} else {
				slog.Error(expr + err.Error())
			}
		}
		c.HostTarget = misc.RemoveDuplicateElement(c.HostTarget)
		c.UrlTarget = misc.RemoveDuplicateElement(c.UrlTarget)
		return
	}
	//判断target字符串是否为类IP/MASK
	if ok := IP.FormatCheck(expr); ok {
		c.HostTarget = append(c.HostTarget, IP.ExprToList(expr)...)
		return
	}
	//判断target字符串是否为类URL
	if url, err := urlparse.Load(expr); err != nil {
		if recursion == true {
			slog.Debug(expr + err.Error())
		} else {
			slog.Error(expr + err.Error())
		}
	} else {
		if url.Scheme != "" {
			c.UrlTarget = append(c.UrlTarget, expr)
			c.HostTarget = append(c.HostTarget, url.Netloc)
			return
		} else {
			c.HostTarget = append(c.HostTarget, url.Netloc)
			return
		}
	}
}
func (c *Config) loadTimeout(i int) {
	c.Timeout = time.Duration(i) * time.Second
}

func (c *Config) loadPort(v interface{}) {
	switch v.(type) {
	case int:
		if v.(int) == 0 {
			return
		}
		c.Port = TOP_1000[:v.(int)]
	case string:
		if v.(string) == "" {
			return
		}
		c.Port = intParam2IntArr(v.(string))
	}
}

func (c *Config) loadOutput(expr string) {
	if expr == "" {
		return
	}
	f, err := os.OpenFile(expr, os.O_CREATE+os.O_RDWR, 0764)
	if err != nil {
		slog.Error(err.Error())
	} else {
		c.Output = f
	}
}

func (c *Config) loadTargetNum() {
	c.HostTargetNum = len(c.HostTarget)
	c.UrlTargetNum = len(c.UrlTarget)
}

func (c *Config) loadPortNum() {
	c.PortNum = len(c.Port)
}

func New() Config {
	return Config{
		HostTarget:    []string{},
		HostTargetNum: 0,
		UrlTarget:     []string{},
		UrlTargetNum:  0,
		Path:          "/",
		Port:          TOP_1000[:400],
		PortNum:       0,
		Output:        nil,
		Proxy:         "",
		Host:          "",
		Threads:       500,
		Timeout:       0,
		Encoding:      "utf-8",
		OSEncoding:    getOSEncoding(),
		NewLine:       getNewline(),
	}
}

func getNewline() string {
	if runtime.GOOS == "windows" {
		return "\r\n"
	} else {
		return "\n"
	}
}

func getOSEncoding() string {
	if runtime.GOOS == "windows" {
		return "gb2312"
	} else {
		return "utf-8"
	}
}

func intParam2IntArr(v string) []int {
	var res []int
	vArr := strings.Split(v, ",")
	for _, v := range vArr {
		var vvArr []int
		if strings.Contains(v, "-") {
			iArr := strings.Split(v, "-")
			if len(iArr) != 2 {
				slog.Error("参数输入错误！！！")
			} else {
				smallNum := misc.Str2Int(iArr[0])
				bigNum := misc.Str2Int(iArr[1])
				if smallNum >= bigNum {
					slog.Error("参数输入错误！！！")
				}
				vvArr = append(vvArr, misc.Xrange(smallNum, bigNum)...)
			}
		} else {
			vvArr = append(vvArr, misc.Str2Int(v))
		}
		res = append(res, vvArr...)
	}
	return res
}

var TOP_1000 = []int{21, 22, 23, 25, 53, 69, 80, 110, 135, 210, 211, 212, 213, 214, 215,
	512, 873, 888, 1080, 1158, 1433, 1521, 2100, 2181, 5120, 5121, 5122, 5123, 5124, 5125,
	5432, 5632, 5900, 5901, 5902, 5984, 6379, 7001, 7080, 54320, 54321, 54322, 54323, 54324, 54325,
	8081, 8082, 8083, 8084, 8085, 8086, 8087, 8088, 8089, 80810, 80811, 80812, 80813, 80814, 80815,
	9001, 9002, 9003, 9004, 9005, 9006, 9007, 9008, 9009, 90010, 90011, 90012, 90013, 90014, 90015,
	9200, 9300, 9418, 9999, 10000, 11211, 27017, 27018, 50000, 92000, 92001, 92002, 92003, 92004, 92005,
	84, 85, 86, 87, 88, 89, 7002, 7003, 7004, 840, 841, 842, 843, 844, 845,
	7070, 7071, 7072, 7073, 7074, 7075, 7076, 7077, 7078, 70700, 70701, 70702, 70703, 70704, 70705,
	8007, 8008, 8200, 90, 801, 8011, 8100, 8012, 8070, 80070, 80071, 80072, 80073, 80074, 80075,
	8181, 800, 18080, 8099, 8899, 8360, 8300, 8800, 8180, 81810, 81811, 81812, 81813, 81814, 81815,
	3000, 41516, 880, 8484, 6677, 8016, 7200, 9085, 5555, 30000, 30001, 30002, 30003, 30004, 30005,
	6080, 8880, 8020, 889, 8881, 38501, 1010, 93, 6666, 60800, 60801, 60802, 60803, 60804, 60805,
	3050, 8787, 2000, 10001, 8013, 6888, 8040, 10021, 2011, 30500, 30501, 30502, 30503, 30504, 30505,
	6060, 7788, 8066, 9898, 6001, 8801, 10040, 9998, 803, 60600, 60601, 60602, 60603, 60604, 60605,
	802, 10003, 8014, 2080, 7288, 8044, 9992, 8889, 5644, 8020, 8021, 8022, 8023, 8024, 8025,
	8021, 8700, 91, 9900, 9191, 3312, 8186, 8735, 8380, 80210, 80211, 80212, 80213, 80214, 80215,
	3333, 2046, 9061, 2375, 9011, 8061, 8093, 9876, 8030, 33330, 33331, 33332, 33333, 33334, 33335,
	70, 8383, 5155, 92, 8188, 2517, 8062, 11324, 2008, 700, 701, 702, 703, 704, 705,
	8987, 8038, 809, 2010, 8983, 7700, 3535, 7921, 9093, 89870, 89871, 89872, 89873, 89874, 89875,
	114, 2012, 701, 8810, 8400, 9099, 8098, 8808, 20000, 1140, 1141, 1142, 1143, 1144, 1145,
	1107, 28099, 12345, 2006, 9527, 51106, 688, 25006, 8045, 11070, 11071, 11072, 11073, 11074, 11075,
	2001, 8035, 10088, 20022, 4001, 2013, 20808, 8095, 106, 20010, 20011, 20012, 20013, 20014, 20015,
	7272, 3380, 3220, 7801, 5256, 5255, 10086, 1300, 5200, 72720, 72721, 72722, 72723, 72724, 72725,
	806, 5050, 8183, 8688, 1001, 58080, 1182, 9025, 8112, 8060, 8061, 8062, 8063, 8064, 8065,
	7081, 8877, 8480, 9182, 58000, 8026, 11001, 10089, 5888, 70810, 70811, 70812, 70813, 70814, 70815,
	5003, 8481, 6002, 9889, 9015, 8866, 8182, 8057, 8399, 50030, 50031, 50032, 50033, 50034, 50035,
	1039, 28080, 5678, 7500, 8051, 18801, 15018, 15888, 38443, 10390, 10391, 10392, 10393, 10394, 10395,
	8990, 3456, 2051, 9098, 444, 9131, 97, 7100, 7711, 89900, 89901, 89902, 89903, 89904, 89905,
	14007, 8184, 7012, 8079, 9888, 9301, 59999, 49705, 1979, 140070, 140071, 140072, 140073, 140074, 140075,
	206, 5156, 8813, 3030, 1790, 8802, 9012, 5544, 3721, 2060, 2061, 2062, 2063, 2064, 2065,
	8056, 7111, 1500, 7088, 5881, 9437, 5655, 8102, 6000, 80560, 80561, 80562, 80563, 80564, 80565,
	8666, 103, 8, 9666, 8999, 9111, 8071, 9092, 522, 86660, 86661, 86662, 86663, 86664, 86665,
	1700, 8036, 8032, 8033, 8111, 60022, 955, 3080, 8788, 17000, 17001, 17002, 17003, 17004, 17005,
	188, 8910, 9022, 10004, 866, 8582, 4300, 9101, 6879, 1880, 1881, 1882, 1883, 1884, 1885,
	8341, 30001, 6890, 8168, 8955, 16788, 8190, 18060, 7041, 83410, 83411, 83412, 83413, 83414, 83415,
	6010, 8898, 9910, 9190, 9082, 8260, 8445, 1680, 8890, 60100, 60101, 60102, 60103, 60104, 60105,
	9704, 5233, 8991, 11366, 7888, 8780, 7129, 6600, 9443, 97040, 97041, 97042, 97043, 97044, 97045,
	2585, 60, 9494, 31945, 2060, 8610, 8860, 58060, 6118, 25850, 25851, 25852, 25853, 25854, 25855,
	8064, 7101, 5081, 7380, 7942, 10016, 8027, 2093, 403, 80640, 80641, 80642, 80643, 80644, 80645,
	6443, 5966, 27000, 7017, 6680, 8401, 9036, 8988, 8806, 64430, 64431, 64432, 64433, 64434, 64435,
	812, 15004, 9110, 8213, 8868, 1213, 8193, 8956, 1108, 8120, 8121, 8122, 8123, 8124, 8125,
	8039, 8600, 50090, 1863, 8191, 65, 6587, 8136, 9507, 80390, 80391, 80392, 80393, 80394, 80395,
	8680, 7999, 7084, 18082, 3938, 18001, 9595, 442, 4433, 86800, 86801, 86802, 86803, 86804, 86805,
	2125, 6090, 10007, 7022, 1949, 6565, 65001, 1301, 19244, 21250, 21251, 21252, 21253, 21254, 21255,
	1005, 22343, 7086, 8601, 6259, 7102, 10333, 211, 10082, 10050, 10051, 10052, 10053, 10054, 10055,
	38086, 666, 6603, 1212, 65493, 96, 9053, 7031, 23454, 380860, 380861, 380862, 380863, 380864, 380865,
	48080, 9086, 10118, 40069, 28780, 20153, 20021, 20151, 58898, 480800, 480801, 480802, 480803, 480804, 480805,
	6546, 3880, 8902, 22222, 19045, 5561, 7979, 5203, 8879, 65460, 65461, 65462, 65463, 65464, 65465,
	9504, 8103, 8567, 1666, 8720, 8197, 3012, 8220, 9039, 95040, 95041, 95042, 95043, 95044, 95045,
	2808, 447, 3600, 3606, 9095, 45177, 19101, 171, 133, 28080, 28081, 28082, 28083, 28084, 28085,
	381, 1443, 15580, 23352, 3443, 1180, 268, 2382, 43651, 3810, 3811, 3812, 3813, 3814, 3815,
	2005, 18002, 2009, 59777, 591, 1933, 9013, 8477, 9696, 20050, 20051, 20052, 20053, 20054, 20055,
	5601, 2901, 2301, 5201, 302, 610, 8031, 5552, 8809, 56010, 56011, 56012, 56013, 56014, 56015,
	5280, 7909, 17003, 1088, 7117, 20052, 1900, 10038, 30551, 52800, 52801, 52802, 52803, 52804, 52805,
	7915, 8384, 9918, 9919, 55858, 7215, 77, 9845, 20140, 79150, 79151, 79152, 79153, 79154, 79155,
	208, 2886, 877, 6101, 5100, 804, 983, 5600, 8402, 2080, 2081, 2082, 2083, 2084, 2085,
	31188, 47583, 8710, 22580, 1042, 2020, 34440, 20, 7703, 311880, 311881, 311882, 311883, 311884, 311885,
	4040, 61081, 12001, 3588, 7123, 2490, 4389, 1313, 19080, 40400, 40401, 40402, 40403, 40404, 40405,
	7899, 30058, 7094, 6801, 321, 1356, 12333, 11362, 11372, 78990, 78991, 78992, 78993, 78994, 78995,
	9096, 8130, 7050, 7713, 40080, 8104, 13988, 18264, 8799, 90960, 90961, 90962, 90963, 90964, 90965,
	9512, 8905, 11660, 1025, 44445, 44401, 17173, 436, 560, 95120, 95121, 95122, 95123, 95124, 95125,
	8488, 8901, 8512, 10443, 9113, 9119, 6606, 22080, 5560, 84880, 84881, 84882, 84883, 84884, 84885,
	333, 73, 7547, 8054, 6372, 223, 3737, 9800, 9019, 3330, 3331, 3332, 3333, 3334, 3335,
	2086, 1002, 9188, 8094, 8201, 8202, 30030, 2663, 9105, 20860, 20861, 20862, 20863, 20864, 20865,
	3010, 7083, 5010, 5501, 309, 1389, 10070, 10069, 10056, 30100, 30101, 30102, 30103, 30104, 30105,
	4180, 10777, 270, 6365, 9801, 1046, 7140, 1004, 9198, 41800, 41801, 41802, 41803, 41804, 41805,
	50100, 8391, 34899, 7090, 6100, 8777, 8298, 8281, 7023, 3377,
}

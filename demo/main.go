package main

import (
	"github.com/PuerkitoBio/goquery"
	"fmt"
	"strings"
	"github.com/smartwalle/ngx"
	"github.com/tealeg/xlsx"
	"github.com/smartwalle/conv4go"
	"os"
	"io/ioutil"
	"flag"
	"time"
	"net/http"
	"math/rand"
	"net/url"
)

func main() {
	//var s = gin.Default()
	//s.GET("/test", func (c *gin.Context) {
	//	fmt.Println(c.ClientIP())
	//	for k, v := range c.Request.Header {
	//		fmt.Println(k,": ", v)
	//	}
	//})
	//s.Run(":9999")

	//var pURL, _ = url.Parse("http://125.94.0.250:8080")
	//var tr = &http.Transport{}
	//tr.Proxy = http.ProxyURL(pURL)
	//
	//var client = &http.Client{}
	//client.Transport = tr
	//
	//var r = ngx.NewRequest("GET", "http://httpbin.org/get")
	//r.AddHeader("user-agent", getUserAgent())
	//r.Client = client
	//fmt.Println("xxx")
	//fmt.Println(r.Exec().MustString())
	//
	//return

	var keyword string
	var city string
	var split bool

	flag.StringVar(&keyword, "k", "", "搜索关键字，例如：Vape Shop")
	flag.StringVar(&city, "c", "", "搜索的城市名称，如果没有传递，默认将从 cities.txt 里面读取所有的城市；例如：Los Angeles")
	flag.BoolVar(&split, "s", false, "如果有多个城市，数据是否存在为多个文件；默认为 false，即存储在一个文件中")
	flag.Parse()

	if strings.TrimSpace(keyword) == "" {
		fmt.Println("请输入需要搜索的关键字")
		return
	}

	var cityList []string
	city = strings.TrimSpace(city)

	if city == "" {
		cityList = loadCityList()
	} else {
		cityList = []string{city}
	}

	_, err := os.Stat("./"+keyword)
	if os.IsNotExist(err) {
		os.MkdirAll("./"+keyword, os.ModePerm)
	}

	if split {
		for _, city := range cityList {
			xFile := xlsx.NewFile()
			sheet, _ := xFile.AddSheet("sheet")

			var row = sheet.AddRow()
			row.AddCell().SetString("Location")
			row.AddCell().SetString("Name")
			row.AddCell().SetString("Phone")
			row.AddCell().SetString("Address")

			search(sheet, keyword, city, 1)
			xFile.Save(fmt.Sprintf("./%s/%s.xlsx", keyword, strings.Replace(city, " ", "_", 0)))
		}
	} else {
		xFile := xlsx.NewFile()
		sheet, _ := xFile.AddSheet("sheet")

		var row = sheet.AddRow()
		row.AddCell().SetString("Location")
		row.AddCell().SetString("Name")
		row.AddCell().SetString("Phone")
		row.AddCell().SetString("Address")

		for _, city := range cityList {
			search(sheet, keyword, city, 1)
		}
		xFile.Save(fmt.Sprintf("./%s/%s.xlsx", keyword, time.Now().Format("2006-01-02")))
	}
}

func loadCityList() []string {
	var cityFile, err = os.Open("./cities.txt")
	if err != nil {
		fmt.Println(err)
		return nil
	}
	defer cityFile.Close()

	data, err :=ioutil.ReadAll(cityFile)
	if err != nil {
		fmt.Println(err)
		return nil
	}
	return strings.Split(string(data), "\n")
}

var agentList = []string{
	"Mozilla/5.0 (Windows NT 6.1; WOW64) AppleWebKit/537.1 (KHTML, like Gecko) Chrome/22.0.1207.1 Safari/537.1",
	"Mozilla/5.0 (X11; CrOS i686 2268.111.0) AppleWebKit/536.11 (KHTML, like Gecko) Chrome/20.0.1132.57 Safari/536.11",
	"Mozilla/5.0 (Windows NT 6.1; WOW64) AppleWebKit/536.6 (KHTML, like Gecko) Chrome/20.0.1092.0 Safari/536.6",
	"Mozilla/5.0 (Windows NT 6.2) AppleWebKit/536.6 (KHTML, like Gecko) Chrome/20.0.1090.0 Safari/536.6",
	"Mozilla/5.0 (Windows NT 6.2; WOW64) AppleWebKit/537.1 (KHTML, like Gecko) Chrome/19.77.34.5 Safari/537.1",
	"Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/536.5 (KHTML, like Gecko) Chrome/19.0.1084.9 Safari/536.5",
	"Mozilla/5.0 (Windows NT 6.0) AppleWebKit/536.5 (KHTML, like Gecko) Chrome/19.0.1084.36 Safari/536.5",
	"Mozilla/5.0 (Windows NT 6.1; WOW64) AppleWebKit/536.3 (KHTML, like Gecko) Chrome/19.0.1063.0 Safari/536.3",
	"Mozilla/5.0 (Windows NT 5.1) AppleWebKit/536.3 (KHTML, like Gecko) Chrome/19.0.1063.0 Safari/536.3",
	"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_8_0) AppleWebKit/536.3 (KHTML, like Gecko) Chrome/19.0.1063.0 Safari/536.3",
	"Mozilla/5.0 (Windows NT 6.2) AppleWebKit/536.3 (KHTML, like Gecko) Chrome/19.0.1062.0 Safari/536.3",
	"Mozilla/5.0 (Windows NT 6.1; WOW64) AppleWebKit/536.3 (KHTML, like Gecko) Chrome/19.0.1062.0 Safari/536.3",
	"Mozilla/5.0 (Windows NT 6.2) AppleWebKit/536.3 (KHTML, like Gecko) Chrome/19.0.1061.1 Safari/536.3",
	"Mozilla/5.0 (Windows NT 6.1; WOW64) AppleWebKit/536.3 (KHTML, like Gecko) Chrome/19.0.1061.1 Safari/536.3",
	"Mozilla/5.0 (Windows NT 6.1) AppleWebKit/536.3 (KHTML, like Gecko) Chrome/19.0.1061.1 Safari/536.3",
	"Mozilla/5.0 (Windows NT 6.2) AppleWebKit/536.3 (KHTML, like Gecko) Chrome/19.0.1061.0 Safari/536.3",
	"Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/535.24 (KHTML, like Gecko) Chrome/19.0.1055.1 Safari/535.24",
	"Mozilla/5.0 (Windows NT 6.2; WOW64) AppleWebKit/535.24 (KHTML, like Gecko) Chrome/19.0.1055.1 Safari/535.24",
	"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_11_0) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/67.0.3396.99 Safari/537.36",
}

func getUserAgent() string {
	return agentList[rand.Intn(len(agentList)-1)]
}

var ipList = []string{
	"https://45.63.71.50:9999",
	"http://47.254.22.115:8080",
	"http://47.89.241.103:8080",
	"https://18.222.124.59:8080",
	"http://35.188.131.123:80",
	"https://67.205.148.246:8080",
	"http://208.67.183.240:80",
	"https://66.82.144.29:8080",
	"http://184.172.238.18:80",
	"http://47.252.1.152:80",
	"https://129.213.76.9:3128",
	"http://173.255.204.5:80",
	"http://52.87.190.26:80",
	"http://159.89.184.107:80",
	"http://184.172.231.86:80",
	"https://167.99.224.142:8080",
	"http://70.169.150.235:48678",
	"http://35.230.34.45:80",
	"http://167.99.236.182:8080",
	"http://47.75.160.100:80",
	"http://206.81.3.55:3128",
	"http://74.208.217.160:80",
	"http://18.222.185.202:1080",
	"http://24.214.133.12:8080",
	"http://50.233.137.36:80",
	"http://167.99.2.42:8080",
	"http://45.77.210.86:8118",
	"http://50.233.137.38:80",
	"http://151.213.212.168:8080",
	"http://34.201.72.254:80",
	"http://74.94.80.101:53281",
	"https://24.220.198.57:8080",
	"http://198.1.122.29:80",
	"http://34.225.249.169:80",
	"http://40.74.243.24:80",
	"https://104.236.51.165:8080",
	"http://45.76.235.248:3128",
	"http://35.234.129.7:80",
	"http://47.254.52.17:80",
	"http://52.8.41.246:3128",
	"http://74.209.243.116:3128",
	"http://165.138.225.250:8080",
	"http://76.79.173.179:3128",
	"http://149.28.111.39:8118",
	"http://129.213.24.42:3128",
	"http://129.146.70.253:80",
	"http://136.24.133.49:52225",
	"http://199.38.114.195:80",
	"http://52.42.234.31:80",
	"http://162.243.118.179:3128",
	"http://31.220.53.67:80",
	"http://108.35.4.138:80",
	"http://47.254.78.25:3128",
	"http://208.180.235.95:80",
	"http://52.204.132.171:8080",
	"http://45.57.204.78:3128",
	"http://67.229.141.187:3128",
	"http://45.35.55.78:8118",
	"http://50.198.134.65:3128",
	"http://64.56.10.240:8080",
	"http://207.201.249.2:8080",
	"http://198.58.123.138:8123",
	"http://104.155.103.87:80",
	"http://45.79.139.169:80",
	"http://67.205.154.3:80",
	"http://192.161.162.101:3128",
	"http://67.205.174.209:3128",
	"http://138.68.232.41:3128",
	"http://138.68.240.218:3128",
	"http://104.236.238.10:3128",
	"http://162.243.107.120:3128",
	"http://66.39.5.184:80",
	"http://138.197.58.55:3128",
	"http://45.55.27.161:3128",
	"http://45.55.16.10:3128",
	"http://40.114.14.173:80",
	"http://45.55.27.17:3128",
	"http://138.197.204.55:3128",
	"http://47.89.181.10:8118",
	"http://162.243.108.129:3128",
	"http://45.55.27.157:3128",
	"http://67.205.146.29:3128",
	"http://45.55.27.15:3128",
	"http://184.178.217.66:3128",
	"http://45.35.55.79:8118",
	"http://47.75.53.32:80",
	"https://45.63.77.42:8080",
	"http://192.169.154.178:80",
	"https://104.236.54.196:8080",
	"http://206.189.237.170:80",
	"http://35.225.208.4:80",
	"http://104.131.214.218:80",
	"https://104.207.141.219:8080",
	"https://192.241.170.125:8118",
	"http://173.212.49.74:8080",
	"http://65.122.77.210:80",
	"http://207.148.29.241:8080",
	"http://64.185.120.62:53281",
	"http://45.35.55.80:8118",
	"http://72.35.40.34:8080",
	"http://206.127.186.17:8080",
	"http://173.45.67.182:9090",
	"https://173.164.26.117:3128",
	"https://170.65.170.70:80",
	"http://162.249.248.218:53281",
	"http://192.30.85.39:8080",
	"http://204.44.117.22:3128",
	"http://198.46.133.240:3128",
	"http://208.100.18.210:3128",
	"http://23.254.110.96:3128",
	"http://23.250.31.74:3128",
	"http://206.214.82.63:3128",
	"http://108.62.124.252:3128",
	"http://50.2.15.71:3128",
	"http://50.2.15.56:3128",
	"http://46.102.101.181:3128",
	"http://46.102.101.75:3128",
	"http://46.102.101.65:3128",
	"http://45.72.69.88:3128",
	"http://45.72.69.92:3128",
	"http://206.214.82.239:3128",
	"http://208.100.18.140:3128",
	"http://216.158.209.133:3128",
	"http://216.158.209.119:3128",
	"http://216.158.209.90:3128",
	"http://216.158.209.88:3128",
	"http://216.158.209.56:3128",
	"http://216.158.209.237:3128",
	"http://216.158.209.214:3128",
	"http://216.158.209.136:3128",
	"http://50.118.223.80:3128",
	"http://50.31.106.49:8800",
	"http://198.55.109.146:3128",
	"http://198.55.109.117:3128",
	"http://198.154.83.138:3128",
	"http://198.154.83.171:3128",
	"http://185.170.216.105:3128",
	"http://185.170.216.99:3128",
	"http://185.170.216.85:3128",
	"http://185.170.216.83:3128",
	"http://185.170.216.48:3128",
	"http://198.55.109.198:3128",
	"http://198.154.83.15:3128",
	"http://198.55.109.55:3128",
	"http://50.118.223.218:3128",
	"http://50.118.223.193:3128",
	"http://50.118.223.240:3128",
	"http://50.118.223.236:3128",
	"http://196.19.166.220:3128",
	"http://198.46.133.161:3128",
	"http://198.46.133.195:3128",
	"http://198.46.133.200:3128",
	"http://198.46.133.213:3128",
	"http://196.19.166.219:3128",
	"http://69.162.162.78:3128",
	"http://67.202.113.83:3128",
	"http://69.147.248.64:3128",
	"http://89.32.71.177:3128",
	"http://185.170.216.182:3128",
	"http://69.147.248.171:3128",
	"http://69.147.248.98:3128",
	"http://69.147.248.70:3128",
	"http://89.35.106.60:3128",
	"http://192.161.162.98:3128",
	"http://192.161.162.61:3128",
	"http://196.19.166.200:3128",
	"http://196.19.166.170:3128",
	"http://185.189.44.104:3128",
	"http://185.189.44.15:3128",
	"http://185.170.216.185:3128",
	"http://50.118.223.234:3128",
	"http://198.154.83.231:3128",
	"http://208.100.18.239:3128",
	"http://206.214.82.204:3128",
	"http://206.214.82.196:3128",
	"http://206.214.82.190:3128",
	"http://206.214.82.179:3128",
	"http://206.214.82.174:3128",
	"http://206.214.82.171:3128",
	"http://208.100.18.59:3128",
	"http://204.44.117.190:3128",
	"http://206.214.82.19:3128",
	"http://198.154.83.223:3128",
	"http://198.154.83.215:3128",
	"http://206.214.82.128:3128",
	"http://204.44.117.165:3128",
	"http://204.44.117.169:3128",
	"http://198.154.83.179:3128",
	"http://34.239.103.154:3128",
	"http://35.194.93.201:8080",
	"http://174.78.222.139:8080",
	"http://208.184.72.106:8080",
	"http://24.154.112.20:8080",
	"http://216.198.170.70:8080",
	"http://174.32.123.230:87",
	"http://72.72.72.123:8090",
	"http://97.72.106.204:87",
	"http://202.5.17.36:8090",
	"http://35.229.125.215:312",
}

func getIp() string {
	return ipList[rand.Intn(len(ipList)-1)]
}

func getClient() *http.Client {
	var pURL, _ = url.Parse(getIp())
	var tr = &http.Transport{}
	tr.Proxy = http.ProxyURL(pURL)

	var client = &http.Client{}
	client.Transport = tr

	return client
}

func search(sheet *xlsx.Sheet, keyword, loc string, page int) {
	var url = ngx.MustURL("https://www.yelp.com/search")
	url.Add("find_desc", keyword)
	url.Add("find_loc", loc)
	url.Add("start", fmt.Sprintf("%d", (page-1)*10))

	var req = ngx.NewRequest("GET", url.String())
	req.AddHeader("user-agent", getUserAgent())
	req.Client = getClient()

	fmt.Println("开始请求：", url.String())

	var rsp = req.Exec()
	if rsp.StatusCode() != http.StatusOK {
		fmt.Println("请求失败")
		return
	}

	var doc, err = goquery.NewDocumentFromReader(rsp.Reader())
	if err != nil {
		fmt.Println(err)
		return
	}

	doc.Find("li.regular-search-result").Each(func(i int, s *goquery.Selection) {
		var name = s.Find("a.biz-name").Find("span").Text()
		var phone = s.Find("span.biz-phone").Text()
		var addr = s.Find("address").Text()

		var row = sheet.AddRow()
		row.AddCell().SetString(strings.TrimSpace(loc))
		row.AddCell().SetString(strings.TrimSpace(name))
		row.AddCell().SetString(strings.TrimSpace(phone))
		row.AddCell().SetString(strings.TrimSpace(addr))
	})

	var pageText = doc.Find("div.page-of-pages").Text()
	var rList = strings.Split(pageText, "of")
	if len(rList) > 1 {
		var totalPage = conv4go.Int(strings.TrimSpace(rList[1]))
		fmt.Printf("当前第 %d 页，共 %d 页: %s \n", page, totalPage, url.String())
		if page < totalPage {
			search(sheet, keyword, loc, page+1)
		}
	}
}

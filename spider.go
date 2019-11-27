package main

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"regexp"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/gocolly/colly"
)

func main() {
	fmt.Println("开始爬取")
	urls := readURL("./urls.txt")
	path := "./images"
	exist, err := PathExists(path)
	if err != nil {
		log.Fatal(err)
	}
	if exist == false {
		err = os.Mkdir(path, os.ModePerm)
		if err != nil {
			log.Fatal(err)
		}
	}
	for k, ul := range urls {
		time.Sleep(time.Second / 10)
		fmt.Printf("正在爬取第%d个URL", k)
		fmt.Println(ul)
		collect(ul)
	}

	fmt.Println("爬取结束")

}

func collect(urlStr string) {
	u, err := url.Parse(urlStr)
	if err != nil {
		log.Fatal(err)
	}

	c := colly.NewCollector(
	//colly.UserAgent("myUserAgent"),
	//colly.AllowedDomains("tmall.com", "taobao.com"),
	)
	// 超时设定
	c.SetRequestTimeout(100 * time.Second)
	// 指定Agent信息
	c.UserAgent = "Mozilla/5.0 (Windows NT 10.0; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/63.0.3239.108 Safari/537.36"
	c.OnRequest(func(r *colly.Request) {
		// Request头部设定
		r.Headers.Set("Host", u.Host)
		r.Headers.Set("Connection", "keep-alive")
		r.Headers.Set("Accept", "*/*")
		r.Headers.Set("Origin", u.Host)
		r.Headers.Set("Referer", urlStr)
		r.Headers.Set("Accept-Encoding", "gzip, deflate")
		r.Headers.Set("Accept-Language", "zh-CN, zh;q=0.9")

		//time.Sleep(2 * time.Second)
	})

	c.OnHTML("title", func(e *colly.HTMLElement) {
		// 1.获取页面id和userid

	})

	c.OnResponse(func(resp *colly.Response) {

		fmt.Println("response received", resp.StatusCode)

		//write("./images/1.html", resp.Body)

		//fmt.Println(string(resp.Body))
		// goquery直接读取resp.Body的内容
		htmlDoc, err := goquery.NewDocumentFromReader(bytes.NewReader(resp.Body))

		// 读取url再传给goquery，访问url读取内容，此处不建议使用
		// htmlDoc, err := goquery.NewDocument(resp.Request.URL.String())

		if err != nil {
			log.Fatal(err)
		}

		htmlDoc.Find(".tb-detail-hd>h1").Each(func(i int, s *goquery.Selection) {
			productName := strings.TrimSpace(s.Text())

			dir := "./images/" + productName
			exist, err := PathExists(dir)

			if err != nil {
				fmt.Printf("get dir error![%v]\n", err)
				return
			}

			if exist {
				fmt.Printf("has dir![%v]\n", dir)
			} else {
				fmt.Printf("no dir![%v]\n", dir)
				// 创建文件夹
				err := os.Mkdir(dir, os.ModePerm)
				if err != nil {
					fmt.Printf("mkdir failed![%v]\n", err)
				} else {
					fmt.Printf("mkdir success!\n")
				}
			}

			htmlDoc.Find("#J_UlThumb img").Each(func(i int, s *goquery.Selection) {
				//productName :=
				//images := make([]string, 100)
				s.Each(func(z int, s *goquery.Selection) {
					imageURL, exist := s.Attr("src")
					if exist {
						fmt.Println(i)
						imageURL = "https:" + strings.Replace(imageURL, "60x60", "430x430", 1)
						name := fmt.Sprintf("%d.jpg", i+1)
						path := dir + "/" + name
						writeImage(path, imageURL)
						//images = append(images, strings.TrimSpace(imageSrc))
						//fmt.Println(imageSrc)
					}
				})
				//fmt.Println(images)
			})

		})

		htmlDoc.Find("meta[name=microscope-data]").Each(func(i int, s *goquery.Selection) {
			//productName :=
			pageId := ""
			userId := ""
			content, b := s.Attr("content")
			if b != true {
				log.Fatal(err)
			}

			contents := strings.Split(strings.TrimSpace(content), ";")
			for _, val := range contents {
				spVal := strings.Split(val, "=")
				if len(spVal) > 1 {
					spVal[0] = strings.TrimSpace(spVal[0])
					spVal[1] = strings.TrimSpace(spVal[1])

					if spVal[0] == "pageId" {
						pageId = spVal[1]
					} else if spVal[0] == "userid" {
						userId = spVal[1]
					}

				}
			}
			if pageId != "" && userId != "" {
				//time.Sleep(1*time.Second)
				collectDetail(pageId, userId)
			}
		})

	})

	c.OnError(func(resp *colly.Response, errHttp error) {
		err = errHttp
	})

	err = c.Visit(urlStr)
}
func collectDetail(pageId, userId string) {
	return
	detailURL := "https://hdc1new.taobao.com/asyn.htm"
	detailURL = detailURL + "?pageId=" + pageId + "&userId=" + userId
	u, err := url.Parse(detailURL)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("开始获取详情图片")
	detailCollector := colly.NewCollector(
	//colly.UserAgent("myUserAgent"),
	//colly.AllowedDomains("tmall.com", "taobao.com"),
	)
	// 超时设定
	detailCollector.SetRequestTimeout(100 * time.Second)
	// 指定Agent信息
	detailCollector.UserAgent = "Mozilla/5.0 (Windows NT 10.0; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/63.0.3239.108 Safari/537.36"
	detailCollector.OnRequest(func(r *colly.Request) {

		// Request头部设定
		r.Headers.Set("Host", u.Host)
		r.Headers.Set("Connection", "keep-alive")
		r.Headers.Set("Accept", "*/*")
		r.Headers.Set("Origin", u.Host)
		r.Headers.Set("Referer", detailURL)
		r.Headers.Set("Accept-Encoding", "gzip, deflate")
		r.Headers.Set("Accept-Language", "zh-CN, zh;q=0.9")

		//time.Sleep(2 * time.Second)
	})

	detailCollector.OnHTML("title", func(e *colly.HTMLElement) {
		fmt.Println("OnHTML-获取详情响应")
	})
	detailCollector.OnResponse(func(resp *colly.Response) {

		fmt.Println("获取详情响应")
		write("./images/2.js", resp.Body)

		// goquery直接读取resp.Body的内容
		htmlDoc, err := goquery.NewDocumentFromReader(bytes.NewReader(resp.Body))
		if err != nil {
			log.Fatal(err)
		}
		htmlDoc.Find(".xx_inner").Each(func(i int, s *goquery.Selection) {
			style, exists := s.Attr("style")
			if exists {
				imageRep := regexp.MustCompile(`^//.+\)&`)
				image := imageRep.FindString(style)
				image = strings.Trim(image, ")")
				fmt.Println(image)
			}
		})
	})
	detailCollector.OnError(func(resp *colly.Response, errHttp error) {
		err = errHttp
	})

	err = detailCollector.Visit(detailURL)
}

func readURL(name string) []string {
	var urls []string

	file, err := os.Open(name)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
	br := bufio.NewReader(file)
	for {
		a, _, c := br.ReadLine()
		if c == io.EOF {
			break
		}
		urls = append(urls, string(a))
	}

	return urls
}

func write(name string, data []byte) {

	newFile, err := os.OpenFile(name, os.O_WRONLY|os.O_TRUNC|os.O_CREATE, 0666)

	if err != nil {
		log.Fatal(err)
	}
	defer newFile.Close()

	_, err = newFile.Write(data)
	if err != nil {
		log.Fatal(err)
	}
}

func writeImage(name, url string) {
	fmt.Println("开始下载图片")
	res, err := http.Get(url)
	if err != nil || res.StatusCode != 200 {
		fmt.Println("下载失败")
		return
	}
	data, err := ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Println("读取图片失败")
		return
	}
	write(name, data)
}

func PathExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

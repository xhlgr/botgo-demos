// 被动消息示例
// 客户端At机器人并传入 hello，  则收到机器人回复： Hello World
package main

import (
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
	"math/rand"
	"encoding/json"
    
	"github.com/tencent-connect/botgo"
	"github.com/tencent-connect/botgo/dto"
	"github.com/tencent-connect/botgo/dto/message"
	"github.com/tencent-connect/botgo/token"
	"github.com/tencent-connect/botgo/websocket"
	yaml "gopkg.in/yaml.v2"
)

var conf struct {
	AppID uint64 `yaml:"appid"`
	Token string `yaml:"token"`
}

var dict = map[string]string{}

func init() {
	content, err := ioutil.ReadFile("../config.yaml")
	if err != nil {
		log.Println("read conf failed")
		os.Exit(1)
	}
	if err := yaml.Unmarshal(content, &conf); err != nil {
		log.Println(err)
		os.Exit(1)
	}
	log.Println(conf)
	
	d, err := ioutil.ReadFile("./dict.json")//读取文件输入到map名为dist，我只是觉得独立文件比较好编辑，抑或是引入一个go文件好呢？
    if err != nil {
        fmt.Print(err)
    }
    err = json.Unmarshal([]byte(string(d)), &dict)
    if err != nil {
        panic(err)
    }
}

func main() {
	token := token.BotToken(conf.AppID, conf.Token)
	api := botgo.NewOpenAPI(token).WithTimeout(3 * time.Second)//NewOpenAPI//NewSandboxOpenAPI
	ctx := context.Background()
	ws, err := api.WS(ctx, nil, "")
	log.Printf("%+v, err:%v", ws, err)
	if err != nil {
		log.Printf("%+v, err:%v", ws, err)
	}
	rand.Seed(time.Now().UnixNano()) //以当前系统时间作为种子参数
	//rand.Intn(100)随机0-99即[0,99)

	var atMessage websocket.ATMessageEventHandler = func(event *dto.WSPayload, data *dto.WSATMessageData) error {
		// 打印一些值 供参考 无实际作用
		fmt.Println(event.Data)
		fmt.Println(data.GuildID, data.ChannelID, data.Content)
		fmt.Println(data.Author.ID, data.Author.Username)

		// 发被动消息到频道
    input := strings.ToLower(message.ETLInput(data.Content))//message.ETLInput去掉at机器人(去掉第一个at)并去掉两边的空格即得到关键词并转小写字母
		fmt.Println(input)
		// 根据词典中的输入，进行输出
		if v, ok := dict[input]; ok {//在词典中直接对应输出
			if _, err := api.PostMessage(context.Background(), data.ChannelID,
				&dto.MessageToCreate{
					Content: message.MentionUser(data.Author.ID) + v,
					MsgID:   data.ID, // 填充 MsgID 则为被动消息，不填充则为主动消息
				},
			); err != nil {
				log.Fatalln(err)
			}
    }else if strings.Index(input, "天气") != -1{//格式：天气+城市名(如北京)
			tmpname := data.Content[strings.Index(data.Content, "天气")+6:]//中文算3*2
			api.PostMessage(ctx, data.ChannelID, &dto.MessageToCreate{MsgID: data.ID, Content: message.MentionUser(data.Author.ID) + getWeather(tmpname)})
		}else{//不在dict中再处理
		    switch input { //下面仅是最新禁言功能示例，机器人需设置为管理员
          case "禁言套餐"://禁发消息的用户
              tmparr := getjinyan()//返回文本和时长秒string
              fmt.Println("禁言数据：",tmparr)
              api.PostMessage(ctx, data.ChannelID, &dto.MessageToCreate{MsgID: data.ID, Content: message.MentionUser(data.Author.ID) + tmparr.Text})
              tsecond,_ := strconv.ParseInt(tmparr.Tsecond, 10, 64)//文本转int64
              mute := &dto.UpdateGuildMute{//个人禁言
            MuteEndTimstamp: strconv.FormatInt(time.Now().Unix()+tsecond, 10),//设置结束时间戳秒,int64按十进制转成string
            //MuteSeconds: tmparr.Tsecond,//设置时长秒
          }
          err := api.MemberMute(ctx, data.GuildID, data.Author.ID, mute)
          if err != nil {
            log.Println(err)
          }
        case "冷静"://全频道禁言，不建议正式使用，否则人人能通过机器人禁言
            mute := &dto.UpdateGuildMute{
              MuteEndTimstamp: strconv.FormatInt(time.Now().Unix()+60, 10),//一分钟
            }
          err := api.GuildMute(ctx, data.GuildID, mute)
          if err != nil {
            log.Println(err)
          }
    		}
		}
		return nil
	}

	intent := websocket.RegisterHandlers(atMessage)     // 注册socket消息处理
	botgo.NewSessionManager().Start(ws, token, &intent) // 启动socket监听

}

//本地数据的结构构造
type JinyanArray struct {
	Text string//返回文本
	Tsecond string//时长秒
}
type JinyanResult struct {
	Jinyan  []JinyanArray
}
//本地禁言套餐数据
func getjinyan() JinyanArray{
    d, err := ioutil.ReadFile("./multidict.json")//读取multidict.json
    if err != nil {
        fmt.Print(err)
    }
    var JinyanRes JinyanResult
    if err = json.Unmarshal(d, &JinyanRes); err != nil {
    	fmt.Printf("Unmarshal err, %v\n", err)
    	return JinyanArray{"禁言数据解析错误，默认禁1分钟","60"}
    }
    Jinyanarray := JinyanRes.Jinyan
    //fmt.Println("JinyanArray: ", Jinyanarray)
    return Jinyanarray[rand.Intn(len(Jinyanarray))]//随机选一个返回
}

// ==========以下，获取地区代码的天气信息==========
type WeatherDate struct {
	Days         string // "例如：20211128"
	Week         string //"例如：星期日"
	Citynm       string //"城市名"
	Weather      string //"例如：多云转阴"
	Temperature  string //"例如：10℃/2℃"
	Weather_curr string //"例如：霾"
	Wind         string //"例如：西南风"
	Winp         string //"例如：1级"
}
type WeatherResult struct {
	Success string
	Result  WeatherDate
}

func getWeather(cityNm string) string {
	//文档https://www.nowapi.com/api/weather.today 文档说示例中sign会不定期调整
	resp, err := http.Get("http://api.k780.com/?app=weather.today&cityNm="+cityNm+"&appkey=10003&sign=b59bc3ef6191eb9f747dd4e83c99f2a4&format=json")
	if err != nil {
		log.Fatalln("天气预报接口请求异常")
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln("天气预报接口数据异常")
	}

	var weatherRes WeatherResult
	if err = json.Unmarshal(body, &weatherRes); err != nil {
		fmt.Printf("Unmarshal err, %v\n", err)
		return ""
	}

	fmt.Println("")
	fmt.Println("body", string(body))
	fmt.Println("")
	fmt.Println("weatherRes: ", weatherRes)

	var weather = weatherRes.Result
	var res = "【天气预报】\r\n" + weather.Days + " " + weather.Week + "\r\n" + weather.Citynm + "天气：" + weather.Weather + ", " + weather.Wind + " " + weather.Winp
	fmt.Println("")
	fmt.Println("res: ", res)

	return res
	// fmt.Println("")
	// fmt.Println(string(body))
}



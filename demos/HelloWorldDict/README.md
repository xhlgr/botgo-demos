本示例演示被动消息使用方法

本人未系统学习Go，代码裁缝，如有不合理也正常。本示例根据[botgosdk收发示例](https://github.com/tencent-connect/botgo/tree/master/examples/receive-and-send)及[botgo-demos](https://github.com/tencent-connect/botgo-demos)拼接而成。

## demo运行步骤

<s>配置机器人的语料（现在默认不校验语料则可不配置，但为了了解内容是否符合规定可以配置审核看看。）</s>
- 1.配置json文件：
  - dict.json：一问一答，key为问题，值为回复
  - multi.json：一问不同答，按需设置，随后在文件构建结构


- 2.执行
```sh
$ go mod tidy
$ go run .
```

- 3.频道内添加对应的机器人，并在子频道内at机器人，发送`hello`, 会收到机器人返回的`Hello World`

# QQ频道Go SDK使用示例

## demos目录说明

- ActiveMessage 主动消息推送
- HelloWorld 被动消息推送
- HelloWorldDict 被动消息推送
  - 拼接了上述被动消息推送和主动消息推送中的`getWeather()`，以及[botgo收发示例](https://github.com/tencent-connect/botgo/tree/master/examples/receive-and-send)

## 运行说明
所有的demo都需要前置下面的操作：

- 申请[QQ机器人](https://bot.q.qq.com/#/home)，获取 `BotAppID`、`Bot Token`，修改`demos/config.example.yaml`为`demos/config.yaml`，并填入内容，格式参考：

```yaml
appid: 123456789
token: "O2AUl44m1mPu4jKVjAwpNtXXXXXXXXXX"
```

## 环境配置参考
- 下载Go1.17.5(编辑时最新版本)。（官方文档表示1.13 及以上版本即可）
- 根据[安装文档](https://go.dev/doc/install)安装。作为小白，这里记录或者是翻译一下linux安装：
  - 下载压缩包：`wget https://go.dev/dl/go1.17.5.linux-amd64.tar.gz`
  - 移除原有go（若有）及解压go到`/usr/local`目录：`rm -rf /usr/local/go && tar -C /usr/local -xzf go1.17.5.linux-amd64.tar.gz`
  - 在`/etc/profile`文件中添加go安装路径及工作路径(可选)，如：
  ```
  export PATH=$PATH:/usr/local/go/bin
  export GOPATH=/path/to/workspace/go
  ```
  - 重新登录SSH或运行：`source /etc/profile`刷新环境变量
  - 运行`go version`查看到版本号即安装完成
- `cd`打开上述demos的目录运行即可。


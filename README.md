## tunnel 

一个简单的将 UDP 报文伪装为 TCP/HTTP 报文 UDP 隧道,在 Linux 下基于原始套接字实现, 在 Windows/Macos 下基于 pcap 实现.

## Features

* 将 UDP 流量通过隧道的形式伪装为 TCP 流量  
* 支持 aes/chacha20 加密,启用加密时设置 method 字段和 password 字段,method 字段为空则不启用加密  
* 支持将流量伪装为 HTTP 流量,默认启用,可用 host 字段设置伪装的 Host,也可以设置 nohttp 字段禁用这项功能  

## Build 
```
go get -u -v github.com/ccsexyz/tunnel
```

## Basic Usage 
  
```
// test.json
{
    "type": "server",
    "localaddr": "127.0.0.1:7676",
    "remoteaddr": ":6666",
    "method": "aes-128-cfb",
    "password": "123"
}
tunnel test.json
```
具体注意事项参考[kcpraw](https://github.com/ccsexyz/kcptun)

## Parameters  
* type: 服务器类型,local 为本地客户端,server 为远程服务器  
* localaddr: 本地监听地址  
* remoteaddr: 对本地客户端指远程服务器的地址,对远程服务器指隧道的目的地址  
* nohttp: 设置为 true 时不进行 HTTP 握手  
* host: 设置 HTTP 伪装的 Host  
* method: 加密方法,可选 aes-128/192/256-cfb/ctr, chacha20, chacha20-ietf, rc4-md5  
* password: 加密使用的密码  
* ignrst: 忽略对端发送的 RST 报文,不推荐使用,仅当你无法过滤 RST 报文时考虑使用  
* expires: 设置每个 session 的过期时间,单位为秒,默认为 30    
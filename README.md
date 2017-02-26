## tunnel 

一个简单的将 UDP 报文伪装为 TCP/HTTP 报文 UDP 隧道,在 Linux 下基于原始套接字实现, 在 Windows/Macos 下基于 pcap 实现.

## 开发目的

解决在使用某些基于 UDP 协议的软件时, 可能会遇到的被运营商 QOS 甚至断流的问题. 

## 基本用法

下面以 kcptun 为例演示如何使用 tunnel 的客户端与服务端:

假定原来的 kcptun 服务端与客户端的参数分别为:
``` 
kcptun_server -l :8080 -t :9999 
kcptun_client -r <server ip>:8080 -l :8080
```
编辑用于服务端和客户端的配置文件

tserver.json
```json
[{
"localaddr": ":8888",
"targetaddr": ":8080"
}]
```

tclient.json
```json
[{
"localaddr": ":10000",
"remoteaddr": "<server ip>:8888"
}]
```

然后分别在服务端与本机执行下面的命令即可建立一条 UDP 隧道
``` 
# 服务器
tserver tserver.json
# 客户端
tclient tclient.json
```

然后更改 kcptun 客户端的启动参数,即可通过建立的隧道连接 kcptun 的服务端
``` 
kcptun_client -r ":10000" -l :8080
```

如果需要两个或者更多的隧道直接按如下方式配置,而不需要启动多个可执行程序,比如
```json 
[
    {
    "localaddr": ":<port1>",
    "remoteaddr": "<ip1>:<remote_port1>"
    }, 
    {
    "localaddr": ":<port2>",
    "remoteaddr": "<ip2>:<remote_port2>"
    }
]
```

## 额外参数  

nohttp: 在 TCP 三次握手后不进行 HTTP 握手,即不伪装为 HTTP 流量,服务端和客户端在这一项上必须保持一致   
host: 设置 HTTP 伪装时所使用的 Host 字段,如设置了 nohttp 这个字段将会失效  
ignrst: 忽略对端发送的 RST 报文,不推荐使用,仅当你无法过滤 RST 报文时考虑使用  
expires: 服务端每个 session 的过期时间,默认为5分钟  
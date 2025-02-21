# 项目简述

本项目的目的是为了合并"域名解析地址，并查询解析归属地"这两个动作，同时想获取同一个域名在不同地区的解析结果；这是个命令行工具。

## 能够帮助你

1. 想要临时使用DOT，DOH获得的解析记录，自行判断是否被污染
2. 直接获取到域名解析结果的IP归属地
3. 模拟客户端所在地址解析结果，也许能帮助你简单的排查不同地区的域名解析问题（通常是CDN的CNAME记录）

## 使用

### DOH

```bash
./dnsregion dns --dns-server https://doh.pub/dns-query www.google.com
```

得到结果

```
DNS SERVER: https://dns.alidns.com/resolve

* www.google.com
|- 31.13.106.4  瑞典|unknown|unknown|脸书
```

### DOT

```bash
./dnsregion dns --dns-server tls://223.5.5.5 www.google.com
```

得到结果

```
DNS SERVER: 223.5.5.5:853

* www.google.com
|- 199.59.148.20        新加坡|unknown|unknown|unknown
```

### EDNS

> 内置client subnet 使用子命令`list`查看，支持自定义subnet资源，但文件必须符合json格式

```bash
./dnsregion edns -d cloudflare-dns.com 贵州_联通 广东_移动
```

得到结果

```
Dns Server:     223.5.5.5:53
Query:  A cloudflare-dns.com

* 贵州_联通
|- * 贵州_贵阳_联通 (edns_client_subnet=221.13.30.242)
|  |- 美国|unknown|unknown(104.16.249.249)
|  |- 美国|unknown|unknown(104.16.248.249)

* 广东_移动
|- * 广东_珠海_移动 (edns_client_subnet=211.136.192.6)
|  |- 美国|unknown|unknown(104.16.249.249)
|  |- 美国|unknown|unknown(104.16.248.249)
|- * 广东_肇庆_移动 (edns_client_subnet=120.196.165.24)
|  |- 美国|unknown|unknown(104.16.249.249)
|  |- 美国|unknown|unknown(104.16.248.249)
|- * 广东_深圳_移动 (edns_client_subnet=120.196.165.24)
|  |- 美国|unknown|unknown(104.16.249.249)
|  |- 美国|unknown|unknown(104.16.248.249)
|- * 广东_东莞_移动 (edns_client_subnet=120.196.165.24)
|  |- 美国|unknown|unknown(104.16.248.249)
|  |- 美国|unknown|unknown(104.16.249.249)
|- * 广东_河源_移动 (edns_client_subnet=120.196.165.24)
|  |- 美国|unknown|unknown(104.16.249.249)
|  |- 美国|unknown|unknown(104.16.248.249)
|- * 广东_佛山_移动 (edns_client_subnet=211.136.192.6)
|  |- 美国|unknown|unknown(104.16.248.249)
|  |- 美国|unknown|unknown(104.16.249.249)
|- * 广东_惠州_移动 (edns_client_subnet=120.196.165.24)
|  |- 美国|unknown|unknown(104.16.248.249)
|  |- 美国|unknown|unknown(104.16.249.249)
|- * 广东_中山_移动 (edns_client_subnet=120.196.165.24)
|  |- 美国|unknown|unknown(104.16.249.249)
|  |- 美国|unknown|unknown(104.16.248.249)
|- * 广东_汕头_移动 (edns_client_subnet=120.196.165.24)
|  |- 美国|unknown|unknown(104.16.249.249)
|  |- 美国|unknown|unknown(104.16.248.249)
|- * 广东_广州_移动 (edns_client_subnet=211.136.192.6)
|  |- 美国|unknown|unknown(104.16.248.249)
|  |- 美国|unknown|unknown(104.16.249.249)
```

## 注意

* 大多数电信运营商的dns服务器只对内提供，只使用公共DNS服务器的话是没有办法100%的获取到不同地区的真实解析情况；
* 解析IP地址使用了[ip2region](https://github.com/lionsoul2014/ip2region)这个项目的内置数据，得到的ip地域信息不一定是准确的
* 本项目内置的edns_client_subnet地址使用的是收集来的全国各地区电信运营商的DNS服务器地址（模拟不同地区），所以有很多重复的条目

* 如果有客户端请求地址，得到的结果可能更准确
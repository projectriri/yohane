# 网关核心插件

## 命令解析

识别特定前缀的消息，并生产命令（[前 commander](https://github.com/projectriri/bot-gateway/tree/2bc3356bf7733435756d8e064cbe2d8bd2f0a470/adapters/commander)）

## 命令过滤

对于解析过的命令，根据消费方频道的协议要求（正则）来帮助应用过滤（[#10](https://github.com/projectriri/bot-gateway/issues/10)）

## 命令别名

`yohane::alias`，`yohane::unalias`

在特定聊天下进行命令映射。（[#11](https://github.com/projectriri/bot-gateway/issues/11)）

先对原命令进行解析，解析后如果命令/参数匹配别名表中的值，则替换。

## 帮助菜单

`help`，`man`

从每个消费者频道注册的帮助信息中读取帮助内容。（[#9](https://github.com/projectriri/bot-gateway/issues/9)）

如果有命令别名，在特定聊天或网页（带参数）上追加命令别名一栏。

命令别名的帮助信息，也是要有的。

## 状态查询

`ping`，`stat`

响应 ping，发出 ping，插件信息和状态一览。

从路由 API 读取网关状态信息。

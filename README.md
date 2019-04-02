# 命令解析器（网关插件）

文档见 [https://projectriri.github.io/bot-gateway/docs/Plugins.html#yohane](https://projectriri.github.io/bot-gateway/docs/Plugins.html#yohane)

## 命令解析

识别特定前缀的消息，并生产命令。

## 命令过滤

对于解析过的命令，根据消费方频道的协议要求来帮助应用过滤（[#10](https://github.com/projectriri/bot-gateway/issues/10)）

## 命令别名

`alias`，`unalias`

在特定聊天下进行命令映射。（[#11](https://github.com/projectriri/bot-gateway/issues/11)）

先对原命令进行解析，解析后如果命令/参数匹配别名表中的值，则替换。

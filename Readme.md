# MongoDB管理工具集 -- mongodb manage command tools

1. 该项目实现了mongodb admin用户的若干管理接口，项目思路来源于[mongodb/mongo-tools](https://github.com/mongodb/mongo-tools); 
2. 该项目的实现是基于命令行的，特修改其代码，将能用到的功能改为支持函数调用，只保留了`mongostat`和`mongotop`两个命令;
3. mongodb的驱动采用[go-mgo/mgo](gopkg.in/mgo.v2)，**该驱动作者已不维护**
4. 仅支持mongodb 3.2+ 和 3.4+版本，3.6+版本未测试

## 支持的命令包括

1. 查询集合索引 （collindexes）
2. 查询集合名称 （collnames）
3. 连接池统计 （connPoolStats）
4. 创建用户 （createuser）
5. DB用量统计 （dbstat）
6. 主机信息 （hostinfo）
7. 副本集配置 （hostinfo）
8. 初始化副本集 （replinit）
9. 副本集状态检测 （replstatus）
10. 服务器状态检测 （serverstatus）
11. DB信息 （showdbs）
12. mongostat
13. mongotop
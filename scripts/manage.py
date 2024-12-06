#!/usr/bin/env python3
import shutil
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
        'en': 'Input the path to install SafeLine WAF',
        'zh': '请输入雷池 WAF 的安装目录'
    },
    'input-mgt-port': {
        'en': 'Input the mgt port',
        'zh': '请输入雷池 WAF 的管理端口'
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
        'en': 'SafeLine WAF does not support %s OS yet',
        'zh': '雷池 WAF 暂时不支持 %s 操作系统'
    },
    'unsupported-arch': {
        'en': 'SafeLine WAF does not support %s processor yet',
        'zh': '雷池 WAF 暂时不支持 %s 处理器'
    },
    'prepare-to-install': {
        'en': 'Will be going to installing SafeLine WAF for you.',
        'zh': '即将为您安装雷池 WAF'
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
        'en': 'The environment does not meet the installation conditions of SafeLine WAF',
        'zh': '当前环境不符合雷池 WAF 的安装条件'
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
        'en': 'Running SafeLine WAF requires Docker, but Docker is not installed',
        'zh': '运行雷池 WAF 依赖 Docker, 但是 Docker 没安装'
    },
    'docker-compose-not-installed': {
        'en': 'Running SafeLine WAF requires Docker Compose, but Docker Compose is not installed',
        'zh': '运行雷池 WAF 依赖 Docker Compose, 但是 Docker Compose 没安装'
    },
    'docker-version-too-low': {
        'en': 'Docker version is too low, it does not match SafeLine WAF',
        'zh': 'Docker 版本太低, 不满足雷池 WAF 的安装需求'
    },
    'if-install-docker': {
        'en': 'Do you want the latest version of Docker to be automatically installed for you?',
        'zh': '是否需要为你自动安装 Docker 的最新版本'
    },
    'if-restart-docker': {
        'en': 'Do you want to restart SafeLine WAF Docker container',
        'zh': '是否需要重启雷池 WAF 的容器'
    },
    'if-update-docker': {
        'en': 'Do you want to update your Docker version?',
        'zh': '是否需要为你自动更新 Docker 版本'
    },
    'install-docker-failed': {
        'en': 'Failed to install Docker. Please try to install Docker manually before installing SafeLine WAF',
        'zh': '安装 Docker 失败, 请尝试手动安装 Docker 后再来安装雷池 WAF'
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
        'en': 'Insufficient disk capacity of "%s", at least 5 GB is required to install SafeLine WAF',
        'zh': '"%s" 的磁盘容量不足，安装雷池 WAF 至少需要 5 GB'
    },
    'pg-pass-contains-invalid-char': {
        'en': 'The POSTGRES_PASSWORD variable contains special characters and cannot be started normally',
        'zh': 'POSTGRES_PASSWORD 变量包含特殊字符, 无法正常启动'
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
        'en': 'Downloading the compose.yaml file',
        'zh': '正在下载 compose.yaml 文件'
    },
    'download-reset-tengine': {
        'en': 'Downloading the reset_tengine script',
        'zh': '正在下载 reset_tengine.sh 文件'
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
        'en': 'SafeLine WAF installation completed',
        'zh': '雷池 WAF 安装完成'
    },
    'upgrade-finish': {
        'en': 'SafeLine WAF upgrade completed',
        'zh': '雷池 WAF 升级完成'
    },
    'go-to-panel': {
        'en': 'SafeLine WAF management panel: https://%s:%s/',
        'zh': '雷池 WAF 管理面板: https://%s:%s/'
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
        'en': 'Installing',
        'zh': '安装通道'
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
        'en': 'Failed to remove safeline installation directory',
        'zh': '删除雷池安装目录失败'
    },
    'uninstall-finish': {
        'en': 'SafeLine WAF uninstall completed',
        'zh': '雷池 WAF 卸载完成'
    },
    'docker-down': {
        'en': 'Stopping SafeLine WAF container',
        'zh': '正在停止雷池 WAF 容器'
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
        'en': 'Do you want to uninstall SafeLine WAF, this operation will delete all data in the directory',
        'zh': '是否确认卸载雷池，该操作会删除目录下所有数据'
    },
    'restart-docker-finish': {
        'en': 'Restart SafeLine WAF docker container completed',
        'zh': '重启雷池 WAF 容器完成'
    },
    'restart': {
        'en': 'RESTART',
        'zh': '重启'
    }
}


lang = ''

def text(label, var=()):
    t = texts.get(label, {
        'en': 'Unknown "%s" (%s)' % (label, var),
        'zh': '未知变量 "%s" (%s)'  % (label, var)
    })
    return t[lang if lang in t else 'en'] % var

BOLD    = 1
DIM     = 1
BLINK   = 5
REVERSE = 7
RED     = 31
GREEN   = 32
YELLOW  = 33
BLUE    = 34
CYAN    = 36

DEBUG = False
LTS = False
IMAGE_CLEAN = False
EN = False
INSTALL = False
DOMAIN = 'waf-ce.chaitin.cn'

def color(t, attrs=[], end=True):
    t = '\x1B[%sm%s' % (';'.join([str(i) for i in attrs]), t)
    if end:
        t = t + '\x1B[m'
    return t


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

def get_url(url):
    try:
        response = urlopen(url)
        content = response.read()
        return content.decode('utf-8')
    except Exception as e:
        log.error(e)

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
                compose_command = 'docker compose'
                version_output = proc[1].strip()
            else:
                log.debug('docker compose can not find detach argument')
        else:
            compose_proc = exec_command('docker-compose', 'version')
            if compose_proc[0] == 0:
                help_proc = exec_command('docker-compose', 'up', '--help')
                if help_proc[0] == 0 and '--detach' in help_proc[1]:
                    compose_command = 'docker-compose'
                    version_output = compose_proc[1].strip()
                else:
                    log.debug('docker-compose can not find detach argument')
            else:
                log.warning(text('docker-compose-not-installed'))

        if version_output != '':
            t = re.findall(r'^Docker Compose version v?(\d+)\.', version_output)
            if len(t) == 0:
                log.warning(text('docker-compose-not-installed'))
            elif int(t[0]) < 2:
                log.warning(text('docker-version-too-low'))
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
                log.warning(text('install-docker-failed'))
                return False


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
                log.warning(text('docker-not-installed'))
            elif int(t[0]) < 20:
                log.warning(text('docker-version-too-low'))
            else:
                break
        else:
            log.warning(text('docker-not-installed'))
            
        action = ui_choice(text('if-install-docker'), [
            ('y', text('yes')),
            ('n', text('no')),
        ])
        if action.lower() == 'n':
            return False
        elif action.lower() == 'y':
            if not install_docker():
                log.warning(text('install-docker-failed'))
                return False

    if not precheck_docker_compose():
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

def docker_down(cwd):
    log.info(text('docker-down'))
    try:
        subprocess.check_call(compose_command+' down', cwd=cwd, shell=True)
        return True
    except Exception:
        return False

def get_url_time(url):
    now = datetime.datetime.now()
    try:
        urlopen(url)
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

    if config['RELEASE'] == '' and LTS:
        config['RELEASE'] = '-lts'
        config['CHANNEL'] = '-lts'

    default_try = False
    if config['MGT_PORT'] == '9443':
        default_try = True
    while not config['MGT_PORT'].isnumeric() or int(config['MGT_PORT']) >= 65536 or int(config['MGT_PORT']) <= 0 or not check_port(config['MGT_PORT']):
        if not default_try:
            config['MGT_PORT'] = '9443'
            default_try = True
        else:
            config['MGT_PORT'] = ui_read(text('input-mgt-port'),None)

    if config['REGION'] == '' and EN:
        config['REGION'] = '-g'

    if not config['POSTGRES_PASSWORD'].isalnum():
        log.info(text('pg-pass-contains-invalid-char'))
        raise Exception(text('pg-pass-contains-invalid-char'))

    if config['IMAGE_PREFIX'] == '' or config['IMAGE_PREFIX'] in pull_failed_prefix:
        config['IMAGE_PREFIX'] = image_source()
        if config['IMAGE_PREFIX'] == '':
            raise Exception(text('fail-to-connect-image-source'))

    config['IMAGE_TAG'] = 'latest'

    with open(env_path, 'w') as f:
        for k in config:
            f.write('%s=%s\n' % (k, config[k]))
    return config

def show_address(mgt_port):
    s = socket.socket(socket.AF_INET, socket.SOCK_DGRAM)
    s.connect(("8.8.8.8", 80))
    local_ip = s.getsockname()[0]
    log.info(text('go-to-panel', (local_ip, mgt_port)))
    log.info(text('go-to-panel', ('0.0.0.0', mgt_port)))

def reset_admin():
    log.info(text('reset-admin'))
    while True:
        p = exec_command('docker', 'inspect','--format=\'{{.State.Health.Status}}\'', 'safeline-mgt')
        if p[0] == 0 and p[1].strip().replace("'",'') == 'healthy':
            break
        elif p[0] != 0:
            log.debug("get safeline-mgt status error: "+str(p[2]))
        log.info("wait safeline-mgt healthy, sleep 5s")
        time.sleep(5)
    proc = exec_command('docker exec safeline-mgt /app/mgt-cli reset-admin --once',shell=True)
    if proc[0] != 0:
        log.warning(proc[2])
    elif proc[1].strip() != '':
        log.info('\n'+proc[1].strip())

def install():
    global INSTALL
    INSTALL = True
    log.info(text('prepare-to-install'))

    if not precheck():
        log.error(text('precheck-failed'))
        return
    log.info(text('precheck-passed'))

    while True:
        safeline_path = ui_read(text('input-target-path'), '/data/safeline')
        if not safeline_path.startswith('/'):
            log.warning(text('invalid-path', safeline_path))
            continue
        if os.path.exists(safeline_path):
            log.warning(text('path-exists', safeline_path))
            continue
        if free_space(safeline_path) < 5 * 1024 * 1024 * 1024:
            log.warning(text('insufficient-disk-capacity'))
            continue
        break

    try:
        os.makedirs(safeline_path)
    except Exception as e:
        log.error(text('fail-to-create-dir', safeline_path) + ' ' + str(e))
        return

    log.info(text('remain-disk-capacity', (safeline_path, humen_size(free_space(safeline_path)))))

    log.info(text('download-compose'))
    if not save_file_from_url('https://'+DOMAIN+'/release/latest/compose.yaml',os.path.join(safeline_path, 'docker-compose.yaml')):
        log.error(text('fail-to-download-compose'))
        return
    if os.path.exists(os.path.join(safeline_path, 'compose.yaml')):
        os.rename(os.path.join(safeline_path, 'compose.yaml'),os.path.join(safeline_path, 'compose.yaml.bak'))

    while True:
        config = generate_config(safeline_path)
        if docker_pull(safeline_path):
            break

        pull_failed_prefix.append(config['IMAGE_PREFIX'])
        log.info(text('try-another-image-source'))

    if not docker_up(safeline_path):
        log.error(text('fail-to-up'))
        return
    
    log.info(text('install-finish'))
    reset_admin()
    show_address(config['MGT_PORT'])

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
        return ui_read(text('input-target-path'),None)

    return safeline_path

def save_file_from_url(url, path):
    log.debug('saving '+url+' to '+path)
    data = get_url(url)
    if data is None:
        return False
    with open(path, 'w') as f:
        f.write(data)
    return True

def upgrade():
    safeline_path = get_installed_dir()

    if not precheck_docker_compose():
        log.error(text('precheck-failed'))
        return

    log.info(text('download-compose'))
    if not save_file_from_url('https://'+DOMAIN+'/release/latest/compose.yaml', os.path.join(safeline_path, 'docker-compose.yaml')):
        log.error(text('fail-to-download-compose'))
        return

    if os.path.exists(os.path.join(safeline_path, 'compose.yaml')):
        os.rename(os.path.join(safeline_path, 'compose.yaml'),os.path.join(safeline_path, 'compose.yaml.bak'))

    while True:
        config = generate_config(safeline_path)
        if docker_pull(safeline_path):
            break

        pull_failed_prefix.append(config['IMAGE_PREFIX'])
        log.info(text('try-another-image-source'))

    if not docker_up(safeline_path):
        log.error(text('fail-to-up'))
        return

    if IMAGE_CLEAN:
        image_clean()

    log.info(text('upgrade-finish'))
    reset_admin()
    show_address(config['MGT_PORT'])
    pass

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

    if not precheck_docker_compose():
        log.error(text('precheck-failed'))
        return

    env_file = os.path.join(safeline_path, '.env')
    if not os.path.exists(env_file):
        log.error(text('fail-to-find-env'))
        return

    config = {}
    read_config(env_file, config)
    if config['POSTGRES_PASSWORD'] == '':
        log.error(text('fail-to-find-postgres-password'))
        return

    if not docker_exec('safeline-pg','psql -U safeline-ce -c "ALTER USER \\"safeline-ce\\" WITH PASSWORD \''+config['POSTGRES_PASSWORD']+'\';"'):
        log.error(text('fail-to-reset-postgres-password'))
        return

    action = ui_choice(text('if-restart-docker'), [
        ('y', text('yes')),
        ('n', text('no')),
    ])

    if action.lower() == 'y':
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
        log.error(text('precheck-failed'))
        return

    if not docker_restart_all(safeline_path):
        return

    log.info(text('restart-docker-finish'))

def backup():
    pass

def uninstall():
    safeline_path = get_installed_dir()

    action = ui_choice(text('if-remove-waf')+": "+safeline_path,[
        ('y', text('yes')),
        ('n', text('no')),
    ])

    if action == 'n':
        return

    if not precheck_docker_compose():
        log.error(text('precheck-failed'))
        return

    if not docker_down(safeline_path):
        log.error(text('fail-to-docker-down'))
        return

    try:
        shutil.rmtree(safeline_path)
    except Exception as e:
        log.debug("remove dir failed: "+str(e))
        log.error(text('fail-to-remove-dir'))

    log.info(text('uninstall-finish'))

def init_global_config():
    global lang, DEBUG, LTS, IMAGE_CLEAN, EN, DOMAIN
    lang = 'zh'
    if '--debug' in sys.argv:
        DEBUG = True

    if '--lts' in sys.argv:
        LTS = True

    if '--image-clean' in sys.argv:
        IMAGE_CLEAN = True

    if '--en' in sys.argv:
        EN = True
        lang = 'en'
        DOMAIN = 'waf.chaitin.com'

def main():
    init_global_config()
    banner()

    log.info(text('hello1'))
    log.info(text('hello2'))
    print()

    if LTS:
        log.info(text('install-channel')+": "+text('lts-release'))

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
        log.error(text('not-linux', platform.system()))
        return

    if platform.machine() not in ('aarch64', 'x86_64', 'AMD64'):
        log.error(text('unsupported-arch', platform.machine()))
        return


    action = ui_choice(text('choice-action'), [
        ('1', text('install')),
        ('2', text('upgrade')),
        ('3', text('uninstall')),
        ('4', text('repair')),
        ('5', text('restart')),
        # ('4', text('backup'))
    ])

    if action == '1':
        install()
    elif action == '2':
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
        print(color(text('talking-group') + '\n', [GREEN]))


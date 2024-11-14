#!/usr/bin/env python3

import sys
import datetime
import platform
import os
from urllib.request import urlopen
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
    'update-config': {
        'en': 'Updating .env configuration files',
        'zh': '正在更新 .env 配置文件'
    },
    'download-compose': {
        'en': 'Downloading the compose.yaml file',
        'zh': '正在下载 compose.yaml 文件'
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
    'install-finish': {
        'en': 'SafeLine WAF installation completed',
        'zh': '雷池 WAF 安装完成'
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
    'upgrade': {
        'en': 'UPDRADE',
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
    }
}


lang = os.getenv('lang')

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
        print('\033[0;%dm[%-5s %s]: %s\033[0m' % (c, l, t, s))

    @staticmethod
    def debug(s):
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
    sys.stdout.write('%s  %s: ' % (
        color(question, [GREEN]),
        color('(%s %s)' % (text('default-value'), default), [YELLOW]),
        ))
    r = input().strip()
    if len(r) == 0:
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
    t = filter(lambda x: 'MemFree' in x, open('/proc/meminfo', 'r').readlines())
    return int(next(t).split()[1]) * 1024

def exec_command(*args):
    try:
        proc = subprocess.run(args, check=False, capture_output=True, universal_newlines=True)
        return proc.returncode, proc.stdout, proc.stderr
    except Exception as e:
        return -1, b'', b''

def exec_command_with_loading(*args):
    try:
        proc = subprocess.Popen(args, stdout=subprocess.PIPE, stderr=subprocess.PIPE, universal_newlines=True)
        loading = ["⣾", "⣽", "⣻", "⢿", "⡿", "⣟", "⣯", "⣷"]
        iloading = 0
        while proc.poll() is None:
            sys.stderr.write('\r' + loading[iloading])
            sys.stderr.flush()
            iloading = (iloading + 1) % len(loading)
            time.sleep(0.1)
        sys.stderr.write('\r')
        sys.stderr.flush()
        return proc.returncode, proc.stdout.read(), proc.stderr.read()
    except Exception as e:
        return -1, b'', str(e)

def install_docker():
    log.info(text('install-docker'))
    r = exec_command_with_loading('bash', '-c', 'curl -ssLk https://get.docker.com | bash')
    return r[0] == 0

def precheck():
    if platform.machine() in ('x86_64', 'AMD64') and 'ssse3' not in open('/proc/cpuinfo', 'r').read().lower():
        log.warning(text('ssse3-not-support'))
        return False

    if free_memory() < 1024 * 1024 * 1024:
        log.warning(text('insufficient-memory'))
        return False

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
            r = install_docker()
            if r == False:
                log.warning(text('install-docker-failed'))
                return False

    while True:
        proc = exec_command('docker', 'compose', 'version')
        if proc[0] == 0:
            t = re.findall(r'^Docker Compose version v(\d+)\.', proc[1])
            if len(t) == 0:
                log.warning(text('docker-compose-not-installed'))
            elif int(t[0]) < 2:
                log.warning(text('docker-version-too-low'))
            else:
                break
        else:c

        action = ui_choice(text('if-update-docker'), [
            ('y', text('yes')),
            ('n', text('no')),
        ])
        if action.lower() == 'n':
            return False
        elif action.lower() == 'y':
            r = install_docker()
            if r == False:
                log.warning(text('install-docker-failed'))
                return False
            
    return True

def docker_pull(cwd):
    log.info(text('docker-pull'))
    try:
        subprocess.check_call('docker compose pull', cwd=cwd, shell=True)
        return True
    except Exception as e:
        return False
    

def docker_up(cwd):
    log.info(text('docker-up'))
    try:
        subprocess.check_call('docker compose up -d', cwd=cwd, shell=True)
        return True
    except Exception as e:
        return False

def get_config(path, _config=None):
    log.info(text('update-config'))
    config = {
        'SAFELINE_DIR': '',
        'POSTGRES_PASSWORD': '',
        'MGT_PORT': '',
        'ARCH': '',
        'CHANNEL': '',
        'REGION': '',
        'IMAGE_PREFIX': '',
        'IMAGE_TAG': '',
        'SUBNET_PREFIX': ''
    }

    if _config is None:
        with open(path + '/.env', 'r') as f:
            for line in f.readlines():
                s = line.index('=')
                if s > 0:
                    k = line[:s].strip()
                    v = line[s + 1:].strip()
                    config[k] = v
    else:
        for k in _config:
            config[k] = _config[k]

    if config['ARCH'] not in ('arm64', 'x86_64'):
        config['ARCH'] = 'arm64' if platform.machine() == 'aarch64' else 'x86_64'


    if config['CHANNEL'] not in ('preview', 'lts'):
        config['CHANNEL'] = 'preview'

    if not config['MGT_PORT'].isnumeric() or int(config['MGT_PORT']) < 65536:
        config['MGT_PORT'] = '9443'

    if config['REGION'] not in ('chinese', 'international'):
        config['REGION'] = 'chinese' if lang == 'zh' else 'international'

    if not config['POSTGRES_PASSWORD'].isalnum():
        log.info(text('pg-pass-contains-invalid-char'))
        raise Exception(text('pg-pass-contains-invalid-char'))

    if config['IMAGE_PREFIX'] not in ('swr.cn-east-3.myhuaweicloud.com/chaitin-safeline', 'chaitin'):
        config['IMAGE_PREFIX'] = 'swr.cn-east-3.myhuaweicloud.com/chaitin-safeline' if lang == 'zh' else 'chaitin'

    config['IMAGE_TAG'] = 'latest'

    with open(path + '/.env', 'w') as f:
        for k in config:
            f.write('%s=%s\n' % (k, config[k]))

    return config

def show_address(config):
    s = socket.socket(socket.AF_INET, socket.SOCK_DGRAM)
    s.connect(("8.8.8.8", 80))
    local_ip = s.getsockname()[0]
    log.info(text('go-to-panel', (local_ip, config['MGT_PORT'])))
    log.info(text('go-to-panel', ('0.0.0.0', config['MGT_PORT'])))

def install():
    log.info(text('prepare-to-install'))

    if not precheck():
        log.error(text('precheck-failed'))
        return
    log.info(text('precheck-passed'))

    config = {}

    while True:
        config['SAFELINE_DIR'] = ui_read(text('input-target-path'), '/data/safeline')
        if not config['SAFELINE_DIR'].startswith('/'):
            log.warning(text('invalid-path', config['SAFELINE_DIR']))
            continue
        if os.path.exists(config['SAFELINE_DIR']):
            log.warning(text('path-exists', config['SAFELINE_DIR']))
            continue
        if free_space(config['SAFELINE_DIR']) < 5 * 1024 * 1024 * 1024:
            log.warning(text('insufficient-disk-capacity'))
            continue
        break

    try:
        os.makedirs(config['SAFELINE_DIR'])
    except Exception as e:
        log.error(text('fail-to-create-dir', config['SAFELINE_DIR']) + ' ' + str(e))
        return

    log.info(text('remain-disk-capacity', (config['SAFELINE_DIR'], humen_size(free_space(config['SAFELINE_DIR'])))))

    config['POSTGRES_PASSWORD'] = ''.join([random.choice(string.ascii_letters + string.digits) for i in range(20)])
    config['SUBNET_PREFIX'] = rand_subnet()

    config = get_config(config['SAFELINE_DIR'], config)

    log.info(text('download-compose'))
    compose_data = get_url('https://waf-ce.chaitin.cn/release/latest/compose.yaml')
    if compose_data is None:
        log.error(text('fail-to-download-compose'))
        return
    
    with open(config['SAFELINE_DIR'] + '/compose.yaml', 'w') as f:
        f.write(compose_data)

    if docker_pull(config['SAFELINE_DIR']) == False:
        log.error(text('fail-to-pull-image'))
        return

    if docker_up(config['SAFELINE_DIR']) == False:
        log.error(text('fail-to-up'))
        return
    
    log.info(text('install-finish'))
    show_address(config)

def main():
    banner()

    log.info(text('hello1'))
    log.info(text('hello2'))
    print()

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
        ('3', text('repair')),
        ('4', text('backup'))
    ])

    if action == '1':
        install()
    elif action == '2':
        upgrade()
    elif action == '3':
        repair()
    elif action == '4':
        backup()

if __name__ == '__main__':
    try:
        main()
    finally:
        print(color(text('talking-group') + '\n', [GREEN]))


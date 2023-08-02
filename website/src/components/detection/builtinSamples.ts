const BuiltinAttackSamples = [
  `
POST /doUpload.action HTTP/1.1
Host: localhost:8080
Content-Length: 1000000000
Content-Type: multipart/form-data; boundary=----WebKitFormBoundaryXd004BVJN9pBYBL2

------WebKitFormBoundaryXd004BVJN9pBYBL2
Content-Disposition: form-data; name="upload"; filename="ok"


<%eval request("sb")%>
------WebKitFormBoundaryXd004BVJN9pBYBL2--
`,
  `GET /fe/Channel/25545911?tabNum=cat%20%2Fetc%2Fhosts HTTP/1.1
Host: monster
User-Agent: Chrome`,
  `GET /scripts/%2e%2e/%2e/Windows/System32/cmd.exe?/c+dir+c HTTP/1.1
Host: a.cn
User-Agent: Chrome
Cookie: 2333;`,
  `GET /?s=/../../../etc/login\0.img HTTP/1.1
Host: monster
User-Agent: Chrome`,
  `POST / HTTP/1.1
Host: monster
User-Agent: Chrome
Content-Type: multipart/form-data; boundary=----WebKitFormBoundary4JGjXRl94NnI4Og7

------WebKitFormBoundary4JGjXRl94NnI4Og7
Content-Disposition: form-data; name="file"; filename="hack.asp%00.png"
Content-Type: application/octet-stream

<%@codepage=65000%><%response.codepage=65001:eval(request("key"))%>
------WebKitFormBoundary4JGjXRl94NnI4Og7--"`,
  `GET /?s=%ac%ed%00%05%73%72%00%1a%63%6f%6d%2e%63%74%2e%61%72%61%6c%65%69%69%2e%74%65%73%74%2e%50%65%72%73%6f%6e%00%00%00%00%00%00%00%01%02%00%08%49%00%03%61%61%61%43%00%03%63%63%63%42%00%03%64%64%64%5a%00%03%65%65%65%4a%00%03%66%66%66%46%00%03%67%67%67%44%00%03%68%68%68%4c%00%03%62%62%62%74%00%12%4c%6a%61%76%61%2f%6c%61%6e%67%2f%53%74%72%69%6e%67%3b%78%70%00%00%00%01%00%62%65%01%00%00%00%00%00%00%00%01%3f%80%00%00%40%00%00%00%00%00%00%00%74%00%03%61%61%61 HTTP/1.1
Host: monster
User-Agent: Chrome`,
  `GET / HTTP/1.1
Host: a.cn
Content-Type: a
Origin: b
User-Agent: c
X-Forwarded-For: d
Cookie: FOOID=1;
A: AAAAAAAAA
Referer: %{(#nike='multipart/form-data').(#dm=@ognl.OgnlContext@DEFAULT_MEMBER_ACCESS).(#_memberAccess?(#_memberAccess=#dm):((#container=#context['com.opensymphony.xwork2.ActionContext.container']).(#ognlUtil=#container.getInstance(@com.opensymphony.xwork2.ognl.OgnlUtil@class)).(#ognlUtil.getExcludedPackageNames().clear()).(#ognlUtil.getExcludedClasses().clear()).(#context.setMemberAccess(#dm)))).(#cmd='ifconfig').(#iswin=(@java.lang.System@getProperty('os.name').toLowerCase().contains('win'))).(#cmds=(#iswin?{'cmd.exe','/c',#cmd}:{'/bin/bash','-c',#cmd})).(#p=new java.lang.ProcessBuilder(#cmds)).(#p.redirectErrorStream(true)).(#process=#p.start()).(#ros=(@org.apache.struts2.ServletActionContext@getResponse().getOutputStream())).(@org.apache.commons.io.IOUtils@copy(#process.getInputStream(),#ros)).(#ros.flush())}`,
  `GET / HTTP/1.1
Host: monster
User-Agent: Chrome
Content-Type: application/x-www-form-urlencoded
Content-Length: 55

<?php echo $a; ?>`,
  `GET / HTTP/1.1
Host: monster
User-Agent: Chrome
Content-Type: application/x-www-form-urlencoded
Content-Length: 55

a=O:9:"FileClass":1:{s:8:"filename";s:10:"config.php";}`,
  `GET / HTTP/1.1
Host: monster
User-Agent: () { :;}; echo; echo $(/bin/ls -al)
Referer: Chrome
Cookie: 2333;`,
  `GET /somepath?q=1'%20or%201=1 HTTP/1.1`,
  `GET /hello?content=hello&user=ftp://a.com HTTP/1.1
Host: www.baidu.com
User-Agent: Chrome
Referer: http://www.qq.com/abc.html`,
  `GET /?flag=%7B%7Bconfig.__class__.__init__.__globals__%5B'os'%5D.popen('id').read()%7D%7D HTTP/1.1
Host: 1.2.3.4
User-Agent: Mozilla/5.0 (Macintosh; Intel Mac OS X 10.15; rv:97.0) Gecko/20100101 Firefox/97.0
Accept: text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,*/*;q=0.8
Accept-Language: zh-CN,zh;q=0.8,zh-TW;q=0.7,zh-HK;q=0.5,en-US;q=0.3,en;q=0.2
Accept-Encoding: gzip, deflate
Connection: close`,
  `POST / HTTP/1.1

<?xml version="1.0" encoding="utf-8"?>
<xsl:stylesheet version="1.0" xmlns:xsl="http://www.w3.org/1999/XSL/Transform">
  <xsl:template match="/fruits">
	<xsl:value-of select="system-property('xsl:vendor')"/>
  </xsl:template>
</xsl:stylesheet>`,
  `GET /somepath?q=<script>alert('XSS')</script> HTTP/1.1`,
  `POST / HTTP/1.1

<?xml version="1.0"?>
<!DOCTYPE data [
<!ELEMENT data (#ANY)>
<!ENTITY file SYSTEM "file:///etc/passwd">
]>
<data>&file;</data>`,
];

const BuiltinNonAttackSamples = [
  `GET / HTTP/1.1
    Host: 1.2.3.4`,

  `POST / HTTP/1.1

name=haha&submit=submit`,
];

export { BuiltinAttackSamples, BuiltinNonAttackSamples };

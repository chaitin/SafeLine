require('./main.less')
//
// var url = window.location.toString()
// var cookie = document.cookie.toString() /* no http-only cookies */
// var referrer =  document.referrer
// var userAgent = window.navigator.userAgent
var timestamp = (function () {
  var d = new Date()
  return d.getFullYear() + '-' + padding(d.getMonth() + 1) + '-' + padding(d.getDate()) + ' ' + padding(d.getHours()) + ':' + padding(d.getMinutes())

  function padding(d) {
    return d < 10 ? '0' + d : d
  }
})()

// $('#info .url').text('被拦截URL: ' + url)
// $('#info .cookie').text('Cookie: ' + cookie)
// $('#info .referrer').text('Referrer: ' + referrer)
// $('#info .user-agent').text('客户端特征: ' + userAgent)
document.getElementById('time').innerHTML = ('拦截时间: ' + timestamp)
// $('#time').text('拦截时间: ' + timestamp)
// $('#info .client-ip').text('用户IP: ' + clientIp)

window.onload = function () {
  var nodes = document.getElementsByTagName("body")[0].childNodes; 
  var fit2inserts = null; 
  for (var i = 0; i < nodes.length; i++) {
    if (nodes[i].nodeType == 8 && nodes[i].data.trimLeft().startsWith("event_id")) { 
      fit2inserts = nodes[i] 
    } 
  }
  /** 
   * 以后引擎按此约定插入新参数：
   * <!-- event_id: ****** type: A anymore: **** -->
  */
  try { 
    var inserts = (document.getElementsByTagName("html")[0].nextSibling || fit2inserts);
    var insertsData = inserts && (inserts.data || "");
    var incertDataList = insertsData.split(" ")
    var getVal = function (key) {
      return incertDataList[incertDataList.indexOf(key + ":") + 1]
    }
    var event_id = getVal('event_id')
    var type = getVal('TYPE')
    // var anymore = getVal('anymore') // 新增参数示例
  } catch (e) { 
    console.log(e) 
  } 
  if (event_id) { 
    document.getElementById("EventID").innerText = "EventID: " + event_id 
  }
  if (type) { 
    document.getElementById("TYPE").innerText = "TYPE: " + type 
  }
};

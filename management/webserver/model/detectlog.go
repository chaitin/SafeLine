package model

import (
	"fmt"
	"math/rand"
	"strings"
	"time"

	"chaitin.cn/dev/go/errors"

	"chaitin.cn/patronus/safeline-2/management/webserver/pkg/constants"
	"chaitin.cn/patronus/safeline-2/management/webserver/pkg/database"
	"chaitin.cn/patronus/safeline-2/management/webserver/utils"
)

// DetectLog is designed to be used in response, not a good naming.
type DetectLog struct {
	DetectLogBasic
	DetectLogDetail
	EventId    string `json:"event_id"` // to eliminate ambiguous
	Website    string `json:"website"`
	AttackType string `json:"attack_type"`
	Module     string `json:"module"`
	Reason     string `json:"reason"`
}

func getRuleModule(ruleId string, attackType int) string {
	var module string
	if strings.HasPrefix(ruleId, "m_rule") { // `m_rule/65543`
		module = "m_rule"
	} else if strings.HasPrefix(ruleId, "/") {
		if attackType == -2 {
			module = "whitelist"
		} else { //  if attackType == -3
			module = "blacklist"
		}
	} else {
		module = ruleId
	}
	return constants.RuleModule[module]
}

func getRuleReason(ruleId string, attackType int) string {
	var reason string
	if strings.HasPrefix(ruleId, "m_rule") { // `m_rule/65543`
		reason = constants.RuleReason[strings.ReplaceAll(ruleId, "m_rule/", "")]
	} else if strings.HasPrefix(ruleId, "/") {
		reason = ruleId[1:]
	} else {
		reason = fmt.Sprintf("检测到 %s 攻击", getAttackType(attackType))
	}
	return reason
}

func getAttackType(at int) string {
	atStr, ok := constants.AttackType[at]
	if !ok {
		atStr = constants.AttackType[62] // unknown
	}
	return atStr
}

func getCountry(country string) string {
	countryStr, ok := constants.CountryCode[country]
	if !ok {
		countryStr = ""
	}
	return countryStr
}

func TransformDetectLog(basic *DetectLogBasic, detail *DetectLogDetail) (*DetectLog, error) {
	if basic == nil {
		return nil, errors.New("basic *DetectLogBasic cannot be nil")
	}
	if detail != nil && basic.EventId != detail.EventId {
		return nil, errors.New("EventId field should be the same for basic and detail")
	}

	basic.Country = getCountry(basic.Country)
	dLog := DetectLog{
		DetectLogBasic: *basic,
		EventId:        basic.EventId,
		Website:        utils.BuildUrl(basic.Protocol, basic.Host, basic.DstPort, basic.UrlPath),
		AttackType:     getAttackType(basic.AttackType),
		Module:         getRuleModule(basic.RuleId, basic.AttackType),
		Reason:         getRuleReason(basic.RuleId, basic.AttackType),
	}
	if detail != nil {
		dLog.DetectLogDetail = *detail
	}

	return &dLog, nil
}

type DetectLogBasic struct {
	ID      uint   `json:"id"               gorm:"primarykey"`
	EventId string `json:"event_id"         gorm:"uniqueIndex;not null"`

	SiteUUID string `json:"site_uuid"       gorm:"column:site_uuid"`
	SrcIp    string `json:"src_ip"          gorm:"index"`
	SocketIp string `json:"socket_ip"       gorm:"index"`

	Protocol int    `json:"protocol"`
	Host     string `json:"host"`
	UrlPath  string `json:"url_path"`
	DstPort  uint   `json:"dst_port"`

	Country  string `json:"country"`
	Province string `json:"province"`
	City     string `json:"city"`

	AttackType int    `json:"attack_type"   gorm:"index"`
	RiskLevel  int    `json:"risk_level"    gorm:"index"`
	Action     int    `json:"action"        gorm:"index"`
	RuleId     string `json:"rule_id"`

	Timestamp int64 `json:"timestamp"       gorm:"index"`
}

type DetectLogDetail struct {
	ID      uint   `                        gorm:"primarykey"`
	EventId string `                        gorm:"uniqueIndex;not null"`

	SrcPort uint   `json:"src_port"`
	DstIp   string `json:"dst_ip"`

	Method      string `json:"method"`
	QueryString string `json:"query_string"`
	StatusCode  uint   `json:"status_code"`

	ReqHeader string `json:"req_header"`
	ReqBody   string `json:"req_body"`
	RspHeader string `json:"rsp_header"`
	RspBody   string `json:"rsp_body"`

	Payload    string `json:"payload"`
	Location   string `json:"location"`
	DecodePath string `json:"decode_path"`
}

func InitDetectLogSamples() {
	db := database.GetDB()
	var detectLogBasicList []DetectLogBasic
	var detectLogDetailList []DetectLogDetail

	timestamp := time.Now().Unix()
	detail := DetectLogDetail{
		SrcPort:     58694,
		DstIp:       "10.2.35.143",
		Method:      "GET",
		QueryString: "",
		ReqHeader:   "GET /webshell.php HTTP/1.1\nUser-Agent: curl/7.77.0\nAccept: */*\n\n\"",
		ReqBody:     "",
		RspHeader:   "",
		RspBody:     "",
		Payload:     "",
		Location:    "urlpath",
		DecodePath:  "",
	}

	protocolList := []int{constants.ProtocolHTTP, constants.ProtocolHTTPS}
	portList := []uint{80, 443}
	provinceList := []string{}
	cityList := []string{}
	ipList := []string{}
	ruleIdList := []string{}

	for i := 0; i < 100; i++ {
		randInt := rand.Intn(1000)
		eventId := utils.RandStr(32)

		basic := DetectLogBasic{}
		basic.EventId = eventId
		basic.SrcIp = ipList[randInt%len(ipList)]
		basic.SocketIp = ipList[randInt%len(ipList)]
		basic.Protocol = protocolList[randInt%len(protocolList)]
		basic.DstPort = portList[randInt%len(portList)]
		basic.Host = fmt.Sprintf("%s.com", utils.RandStr(5))
		basic.UrlPath = fmt.Sprintf("/%s", utils.RandStr(10))
		basic.Country = "CN"
		basic.Province = provinceList[randInt%len(provinceList)]
		basic.City = cityList[randInt%len(cityList)]
		basic.AttackType = randInt % 32
		basic.RiskLevel = randInt % 4
		basic.Action = randInt % 2
		basic.RuleId = ruleIdList[randInt%len(ruleIdList)]
		basic.Timestamp = timestamp - int64(randInt*100)

		detailCopy := detail
		detailCopy.EventId = eventId
		detectLogBasicList = append(detectLogBasicList, basic)
		detectLogDetailList = append(detectLogDetailList, detailCopy)
	}
	db.CreateInBatches(detectLogBasicList, 100)
	db.CreateInBatches(detectLogDetailList, 100)
}

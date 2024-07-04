package fvm

import (
	"encoding/json"
	"fmt"
	"strings"

	"gorm.io/gorm"

	"chaitin.cn/patronus/safeline-2/management/webserver/pkg/database"

	"chaitin.cn/patronus/safeline-2/management/webserver/pkg/config"

	"chaitin.cn/dev/go/errors"
	"chaitin.cn/dev/go/log"
	"chaitin.cn/patronus/safeline-2/management/webserver/model"
	"chaitin.cn/patronus/safeline-2/management/webserver/pkg/fvm/fsl"
)

const (
	Main              = "main"
	Preprocess        = "preprocess"
	MassPackage       = "mass_package"
	PolicyRule        = "policy_rule"
	PolicyRuleGlobal  = "policy_rule_global"
	PolicyGroup       = "policy_group"
	PolicyGroupDetect = "policy_group_detect"
	WeblogGenerate    = "weblog_generate"
	ReqWeblogGenerate = "req_weblog_generate"
	RspWeblogGenerate = "rsp_weblog_generate"
	CommonLog         = "common_log"

	STR  = "STRING"
	INT  = "INTEGER"
	INET = "INET"
	BOOL = "BOOLEAN"
)

var logger = log.GetLogger("fvm")

func InitFVMBytecode() {
	go func() {
		fullFSL, err := GenerateFullFSL(database.GetDB().DB)
		if err != nil {
			logger.Fatalln("Failed to generate fsl", err)
		}
		if err = CompileAndSave(fullFSL); err != nil {
			logger.Fatalln("Failed to compile and save fsl", err)
		}
	}()
}

func PushFSL(tx *gorm.DB) error {
	fullFSL, err := GenerateFullFSL(tx)
	if err != nil {
		logger.Debugf("Full fsl: %s", fullFSL)
		logger.Error(err)
		return err
	}
	if err = CompileAndPush(fullFSL, config.GlobalConfig.Detector.Addr); err != nil {
		logger.Error(err)
		return err
	}
	return nil
}

func GenerateFullFSL(tx *gorm.DB) (fullFSL string, err error) {
	logger.Info("Generate FSL")

	pFSL, err := PreprocessTable(tx)
	if err != nil {
		return "", err
	}
	fullFSL += pFSL

	fullFSL += MassPackageTable()

	prFSL, err := PolicyRuleTable(tx)
	if err != nil {
		return "", err
	}
	fullFSL += prFSL

	pgFSL, err := PolicyGroupTable(tx)
	if err != nil {
		return "", err
	}
	fullFSL += pgFSL

	fullFSL += WeblogGenerateTable()

	fullFSL += MainTable()

	return fullFSL, err
}

func PreprocessTable(db *gorm.DB) (pFSL string, err error) {
	pFSL += fsl.CreateTable(Preprocess, nil)
	pFSL += fsl.AppendInto(Preprocess, "http_parse", "", "http_parse()")
	pFSL += fsl.AppendInto(Preprocess, "init_variables", "", fsl.Actions(
		fsl.Set("@site_uuid", STR, "get_http_all_headers('SL-CE-SUID')"),
		fsl.Set("@src_ip", INET, "socket.src_ip"),
		fsl.Set("@skip_remaining", BOOL, "false"),
		fsl.Set("@rule_id", STR, fsl.Quote("")),
		fsl.Set("@attack_type", INT, "-1"),
		fsl.Set("@final_action", INT, "0"),
		fsl.Set("@risk_level", INT, "0"),
		fsl.Set("@host", STR, "http.decayed_host"),
		fsl.Set("@start_time", INT, "now.timestamp_us"),
		fsl.Set("@transport_scheme", STR, "transport.scheme"),
		fsl.Set("@forward_log", INT, "0"),
		fsl.Set("@convert_to_alog", INT, "1"),
		fsl.Set("@alog_save_to_db", INT, "0"),
		fsl.Set("@convert_to_blog", INT, "0"),
		fsl.Set("@blog_save_to_db", INT, "0"),
		fsl.Set("@convert_to_dlog", INT, "0"),
		fsl.Set("@dlog_save_to_db", INT, "0"),
	))
	pFSL += fsl.AppendInto(Preprocess, "set_host", "@host::STRING = ''", fsl.Actions(fsl.Set("@host", STR, "inet_to_string(socket.dst_ip)")))
	pFSL += fsl.AppendInto(Preprocess, "adapt_http2", "@transport_scheme::STRING = 'http2'", fsl.Actions(fsl.Set("@transport_scheme", STR, fsl.Quote("https"))))

	// get src_ip
	sc, err := model.GetSrcIPConfig(db)
	if err != nil {
		return "", err
	}
	if sc.Source == model.HTTPHeader {
		key := fmt.Sprintf("string_to_inet(get_http_all_headers(%s))", fsl.Quote(sc.Value))
		pFSL += fsl.AppendInto(Preprocess, "set_src_ip", fmt.Sprintf("%s != '::/0'::INET", key), fsl.Actions(fsl.Set("@src_ip", INET, key)))
	}
	return pFSL, nil
}

func MassPackageTable() (mpFSL string) {
	mpFSL += fsl.CreateTable(MassPackage, nil)
	mpFSL += fsl.AppendInto(MassPackage, "mass_package_detect", "string.length(http.raw_req_body) > 1048576", fsl.Actions(
		fsl.Set("@skip_remaining", BOOL, "true"),
		fsl.Set("@rule_id", STR, fsl.Quote("mass_package")),
		fsl.Set("@final_action", INT, "0"),
		fsl.Set("@convert_to_dlog", INT, "1"),
		fsl.Set("@dlog_save_to_db", INT, "1"),
		fsl.Set("@attack_type", INT, "-4"),
	))
	return mpFSL
}

func PolicyRuleTable(db *gorm.DB) (prFSL string, err error) {
	prFSL += fsl.CreateTable(PolicyRuleGlobal, map[fsl.State]fsl.State{fsl.S_ABORT: fsl.S_RETURN})

	var policyRuleList []model.PolicyRule
	res := db.Model(&model.PolicyRule{}).Where(&model.PolicyRule{IsEnabled: true}).Order("id desc").Find(&policyRuleList)
	if res.Error != nil {
		err = errors.New(fmt.Sprintf("Error occurred when fetch PolicyRule: %s", res.Error))
	}

	for _, policyRule := range policyRuleList {

		var patternBytes []byte
		patternBytes, err = policyRule.Pattern.MarshalJSON()
		if err != nil {
			return "", err
		}
		var patterns []model.PolicyRulePattern
		err = json.Unmarshal(patternBytes, &patterns)
		if err != nil {
			return "", err
		}

		var wheres []string
		for _, pattern := range patterns {
			// transform key to string
			var key string
			switch pattern.K {
			case model.KeySrcIP:
				key = "inet_to_string(@src_ip::INET)"
			case model.KeyURI:
				key = "@uri_decoded::STRING"
			case model.KeyHost:
				key = "http.host"
			default:
				return "", errors.New("wrong Key")
			}

			switch pattern.Op {
			case model.OpEq:
				wheres = append(wheres, fmt.Sprintf("string.equals_case(%s, %s)", key, fsl.Quote(pattern.V)))
			case model.OpMatch:
				// . -> \.
				// * -> .*
				valueQuote := strings.ReplaceAll(fsl.Quote(strings.ReplaceAll(pattern.V, ".", "\\.")), "*", ".*")
				wheres = append(wheres, fmt.Sprintf("pcre_match(%s, %s)", valueQuote, key))
			case model.OpCIDR:
				// use INET for CIDR compare
				key = "@src_ip::INET"
				wheres = append(wheres, fmt.Sprintf("%s IN CIDR(%s)", key, fsl.Quote(pattern.V)))
			case model.OpHas:
				wheres = append(wheres, fmt.Sprintf("string.contains_case(%s, %s)", key, fsl.Quote(pattern.V)))
			case model.OpPrefix:
				l := len(pattern.V)
				wheres = append(wheres, fmt.Sprintf("string.equals_case(string.substr(%s, 0, %d), %s)", key, l, fsl.Quote(pattern.V)))
			case model.OpRe:
				wheres = append(wheres, fmt.Sprintf("pcre_match(%s, %s)", fsl.Quote(pattern.V), key))
			default:
				return "", errors.New("wrong Operator")
			}
		}

		// save log only when deny and dry_run
		logOption := "1"
		attackType := "-3"
		// DROP or ACCEPT
		appendix := "DROP"
		if policyRule.Action == 0 {
			logOption = "0"
			attackType = "-2"
			appendix = "set_ctx_allow(), ACCEPT" // this selector seems no use
		}
		prFSL += fsl.AppendInto(PolicyRuleGlobal, fmt.Sprintf("%s_%d", PolicyRuleGlobal, policyRule.ID), fsl.Wheres(wheres...), fsl.Actions(
			fsl.Set("@skip_remaining", BOOL, "true"),
			fsl.Set("@rule_id", STR, fsl.Quote(fmt.Sprintf("/%s", policyRule.Comment))),
			fsl.Set("@final_action", INT, fmt.Sprintf("%d", policyRule.Action)),
			fsl.Set("@convert_to_dlog", INT, logOption),
			fsl.Set("@dlog_save_to_db", INT, logOption),
			fsl.Set("@attack_type", INT, attackType),
			appendix,
			// no use in v1.1
			//fsl.Set("@risk_level", INT, ""),
		))
	}

	prFSL += fsl.CreateTable(PolicyRule, nil)
	prFSL += fsl.AppendInto(PolicyRule, "preprocess", "", fsl.Actions(fsl.Set("@uri_decoded", STR, "url_decode(http.uri)")))
	prFSL += fsl.AppendInto(PolicyRule, "policy_rule_global", "", fsl.Goto(PolicyRuleGlobal))
	prFSL += fsl.AppendInto(PolicyRule, "set_status_code", "@final_action::INTEGER = 1", "set_status_code(403)")

	return prFSL, err
}

func PolicyGroupTable(db *gorm.DB) (pgFSL string, err error) {
	skynetConfigStr, err := model.GetSkynetConfig(db)
	if err != nil {
		return "", err
	}

	pgFSL += fsl.CreateTarget("skynet_config", fsl.Quote(skynetConfigStr))

	pgFSL += fsl.CreateTable(PolicyGroupDetect, map[fsl.State]fsl.State{fsl.S_ABORT: fsl.S_RETURN})
	pgFSL += fsl.AppendInto(PolicyGroupDetect, "do_detect", "", "target 'skynet_config' type skynet do detect()")

	pgFSL += fsl.CreateTable(PolicyGroup, nil)
	pgFSL += fsl.AppendInto(PolicyGroup, "do_detect", "", fsl.Goto(PolicyGroupDetect))
	pgFSL += fsl.AppendInto(PolicyGroup, "dlog_persistence", "skynet.enable_log = true", fsl.Actions(
		fsl.Set("@convert_to_dlog", INT, "1"),
		fsl.Set("@dlog_save_to_db", INT, "1"),
	))
	pgFSL += fsl.AppendInto(PolicyGroup, "set_status_code", "skynet.action = 1", "set_status_code(403)")
	pgFSL += fsl.AppendInto(PolicyGroup, "hit_policy_group", "skynet.rule_id != ''", fsl.Actions(
		fsl.Set("@skip_remaining", BOOL, "true"),
		fsl.Set("@rule_id", STR, "skynet.rule_id"),
		fsl.Set("@attack_type", INT, "skynet.main_attack"),
		fsl.Set("@final_action", INT, "skynet.action"),
		fsl.Set("@risk_level", INT, "skynet.risk_level"),
		fsl.Set("@convert_to_dlog", INT, "1"),
	))
	return pgFSL, nil
}

func CommonLogTable() (clFSL string) {
	clFSL += fsl.CreateTable(CommonLog, nil)
	clFSL += fsl.AppendInto(CommonLog, "batch_set_payload_string", "", "weblog_payload_set_string_batch({'event_id', 'scheme', 'dst_ip', 'socket_ip', 'src_ip'}::ARRAY(STRING), {extra.uuid, transport.scheme, inet_to_string(socket.dst_ip), inet_to_string(socket.src_ip), inet_to_string(@src_ip::INET)}::ARRAY(STRING))")
	clFSL += fsl.AppendInto(CommonLog, "batch_set_payload_int", "", "weblog_payload_set_int_batch({'src_port', 'dst_port'}::ARRAY(STRING), {socket.src_port, socket.dst_port}::ARRAY(INTEGER))")
	clFSL += fsl.AppendInto(CommonLog, "batch_set_extra_bool", "", "weblog_extra_set_bool_batch({'alog_save_to_db', 'convert_to_blog', 'forward_log', 'convert_to_alog', 'convert_to_dlog', 'dlog_save_to_db', 'blog_save_to_db'}::ARRAY(STRING), {@alog_save_to_db::INTEGER, @convert_to_blog::INTEGER, @forward_log::INTEGER, @convert_to_alog::INTEGER, @convert_to_dlog::INTEGER, @dlog_save_to_db::INTEGER, @blog_save_to_db::INTEGER}::ARRAY(INTEGER))")
	clFSL += fsl.AppendInto(CommonLog, "generate", "", "weblog_generate()")
	return clFSL
}

func ReqWeblogGenerateTable() (reqWGFSL string) {
	reqWGFSL += fsl.CreateTable(ReqWeblogGenerate, nil)
	reqWGFSL += fsl.AppendInto(ReqWeblogGenerate, "init_http_req_body", "", fsl.Actions(fsl.Set("@raw_req_body", STR, "http.raw_req_body")))
	reqWGFSL += fsl.AppendInto(ReqWeblogGenerate, "cut_http_req_body", "string.length(@raw_req_body::STRING) > 1048576", fsl.Actions(fsl.Set("@raw_req_body", STR, "string.substr(http.raw_req_body, 0, 1048576)")))
	reqWGFSL += fsl.AppendInto(ReqWeblogGenerate, "batch_set_payload_binary", "", "weblog_payload_set_binary_batch({'req_payload', 'req_body'}::ARRAY(STRING), {skynet.payload, @raw_req_body::STRING}::ARRAY(STRING))")
	reqWGFSL += fsl.AppendInto(ReqWeblogGenerate, "batch_set_payload_string", "", "weblog_payload_set_string_batch({'site_uuid', 'req_header', 'req_detector_name', 'req_proxy_name', 'req_rule_id', 'req_location', 'req_decode_path', 'user_agent', 'query_string', 'method', 'host', 'referer', 'url_path', 'cookie'}::ARRAY(STRING), {@site_uuid::STRING, http.raw_req_header, fusion.node_id, extra.proxy_name, @rule_id::STRING, skynet.location, skynet.decode_path, http.user_agent, http.decoded_query, http.raw_method, @host::STRING, http.referer, http.path, http.raw_cookie}::ARRAY(STRING))")
	reqWGFSL += fsl.AppendInto(ReqWeblogGenerate, "batch_set_payload_int", "", "weblog_payload_set_int_batch({'req_start_time', 'req_detect_time', 'req_block_reason', 'req_attack_type', 'req_risk_level', 'req_end_time', 'req_action'}::ARRAY(STRING), {extra.req_start_time, now.timestamp_us - @start_time::INTEGER, fusion.block_reason, @attack_type::INTEGER, @risk_level::INTEGER, extra.req_end_time, @final_action::INTEGER}::ARRAY(INTEGER))")
	reqWGFSL += fsl.AppendInto(ReqWeblogGenerate, "common_log", "", fsl.Goto(CommonLog))
	return reqWGFSL
}

func RspWeblogGenerateTable() (rspWGFSL string) {
	rspWGFSL += fsl.CreateTable(RspWeblogGenerate, nil)
	rspWGFSL += fsl.AppendInto(RspWeblogGenerate, "init_http_rsp_body", "", fsl.Actions(fsl.Set("@raw_rsp_body", STR, "http.raw_rsp_body")))
	rspWGFSL += fsl.AppendInto(RspWeblogGenerate, "cut_http_rsp_body", "string.length(@raw_rsp_body::STRING) > 1048576", fsl.Actions(fsl.Set("@raw_rsp_body", STR, "string.substr(http.raw_rsp_body, 0, 1048576)")))
	rspWGFSL += fsl.AppendInto(RspWeblogGenerate, "batch_set_payload_binary", "", "weblog_payload_set_binary_batch({'rsp_payload', 'rsp_body'}::ARRAY(STRING), {skynet.payload, @raw_rsp_body::STRING}::ARRAY(STRING))")
	rspWGFSL += fsl.AppendInto(RspWeblogGenerate, "batch_set_payload_string", "", "weblog_payload_set_string_batch({'rsp_rule_id', 'rsp_location', 'rsp_decode_path', 'rsp_header', 'rsp_detector_name', 'rsp_proxy_name'}::ARRAY(STRING), {@rule_id::STRING, skynet.location, skynet.decode_path, http.raw_rsp_header, fusion.node_id, extra.proxy_name}::ARRAY(STRING))")
	rspWGFSL += fsl.AppendInto(RspWeblogGenerate, "batch_set_payload_int", "", "weblog_payload_set_int_batch({'rsp_detect_time', 'status_code', 'rsp_start_time', 'rsp_end_time', 'rsp_action', 'rsp_attack_type', 'rsp_risk_level', 'rsp_block_reason'}::ARRAY(STRING), {now.timestamp_us - @start_time::INTEGER, http.status_code, extra.rsp_start_time, extra.rsp_end_time, @final_action::INTEGER, @attack_type::INTEGER, @risk_level::INTEGER, fusion.block_reason}::ARRAY(INTEGER))")
	rspWGFSL += fsl.AppendInto(RspWeblogGenerate, "common_log", "", fsl.Goto(CommonLog))
	return rspWGFSL
}

func WeblogGenerateTable() (wgFSL string) {
	wgFSL += CommonLogTable()
	wgFSL += ReqWeblogGenerateTable()
	wgFSL += RspWeblogGenerateTable()
	wgFSL += fsl.CreateTable(WeblogGenerate, nil)
	wgFSL += fsl.AppendInto(WeblogGenerate, "req", "fusion.stage = 0", fsl.Goto(ReqWeblogGenerate))
	wgFSL += fsl.AppendInto(WeblogGenerate, "rsp", "fusion.stage = 1", fsl.Goto(RspWeblogGenerate))
	return wgFSL
}

func MainTable() (mainFSL string) {
	mainFSL += fsl.CreateTable(Main, nil)
	mainFSL += fsl.AppendInto(Main, "m_preprocess", "", fsl.Goto(Preprocess))
	mainFSL += fsl.AppendInto(Main, "m_mass_package", "", fsl.Goto(MassPackage))
	mainFSL += fsl.AppendInto(Main, "m_policy_rule", "@skip_remaining::BOOLEAN = false", fsl.Goto(PolicyRule))
	mainFSL += fsl.AppendInto(Main, "m_policy_group", "@skip_remaining::BOOLEAN = false", fsl.Goto(PolicyGroup))
	mainFSL += fsl.AppendInto(Main, "m_weblog_generate", "", fsl.Goto(WeblogGenerate))
	mainFSL += "ENTRYPOINT TABLE main;"
	return mainFSL
}

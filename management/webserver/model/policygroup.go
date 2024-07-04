package model

import (
	"encoding/json"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"

	"chaitin.cn/dev/go/errors"
	"chaitin.cn/patronus/safeline-2/management/webserver/pkg/constants"
	"chaitin.cn/patronus/safeline-2/management/webserver/pkg/database"
)

type PolicyGroup map[string]string

const (
	StrictMode  = "strict"
	DefaultMode = "default"
	DisableMode = "disable"

	SocketIP   = "socket_ip"
	HTTPHeader = "http_header"
)

type SrcIPConfig struct {
	Source string `json:"source"` // socket_ip, http_header
	Value  string `json:"value"`
}

type SkynetConfigModule struct {
	DetectConfig        interface{} `json:"detect_config,omitempty"`
	HighRiskAction      string      `json:"high_risk_action"`
	MediumRiskAction    string      `json:"medium_risk_action"`
	LowRiskAction       string      `json:"low_risk_action"`
	HighRiskEnableLog   int         `json:"high_risk_enable_log"`
	MediumRiskEnableLog int         `json:"medium_risk_enable_log"`
	LowRiskEnableLog    int         `json:"low_risk_enable_log"`
	State               string      `json:"state"`
}

type SkynetConfig struct {
	DetectConfig string      `json:"detect_config"`
	DeepDetect   bool        `json:"deep_detect"`
	Timeout      interface{} `json:"timeout"`

	Modules map[string]*SkynetConfigModule `json:"modules"`
}

func GetSrcIPConfig(db *gorm.DB) (*SrcIPConfig, error) {
	var scOption Options
	db.Where(&Options{Key: constants.SrcIPConfig}).First(&scOption)

	var sc SrcIPConfig
	err := json.Unmarshal([]byte(scOption.Value), &sc)
	if err != nil {
		return nil, err
	}

	return &sc, nil
}

func initSrcIPConfig() error {
	db := database.GetDB()

	// policyGroupGlobal config, three mode: strict/default/disable
	var srcIPConfig = SrcIPConfig{
		Source: SocketIP,
		Value:  "",
	}

	scStr, err := json.Marshal(srcIPConfig)
	if err != nil {
		return err
	}

	pgg := Options{Key: constants.SrcIPConfig, Value: string(scStr)}
	db.Clauses(clause.OnConflict{DoNothing: true}).Create(&pgg)

	return nil
}

func GetSkynetConfig(db *gorm.DB) (string, error) {
	// skynetConfigStr skynet config in 5.3.9
	const skynetConfigStr = "{\"decode_config\":{\"decode_methods\":[\"url decode\",\"JSON\",\"base64\",\"hex\",\"eval\",\"XML\",\"PHP deserialize\",\"utf7\"]},\"deep_detect\":false,\"modules\":{\"m_asp_code_injection\":{\"high_risk_action\":\"deny\",\"high_risk_enable_log\":1,\"low_risk_action\":\"continue\",\"low_risk_enable_log\":1,\"medium_risk_action\":\"continue\",\"medium_risk_enable_log\":1,\"state\":\"enabled\"},\"m_cmd_injection\":{\"high_risk_action\":\"deny\",\"high_risk_enable_log\":1,\"low_risk_action\":\"continue\",\"low_risk_enable_log\":1,\"medium_risk_action\":\"continue\",\"medium_risk_enable_log\":1,\"state\":\"enabled\"},\"m_csrf\":{\"high_risk_action\":\"deny\",\"high_risk_enable_log\":1,\"low_risk_action\":\"continue\",\"low_risk_enable_log\":1,\"medium_risk_action\":\"continue\",\"medium_risk_enable_log\":1,\"state\":\"enabled\"},\"m_file_include\":{\"detect_config\":{\"detect_suspicious_schema\":true},\"high_risk_action\":\"deny\",\"high_risk_enable_log\":1,\"low_risk_action\":\"continue\",\"low_risk_enable_log\":1,\"medium_risk_action\":\"continue\",\"medium_risk_enable_log\":1,\"state\":\"enabled\"},\"m_file_upload\":{\"detect_config\":{\"action_of_handling_custome_file\":[\"deny transmission\"],\"detect_file_content\":true,\"detect_real_file_content\":false,\"enable_module_in_detecting_file_content\":[\"java\",\"asp_code_injection\",\"php_code_injection\"],\"file_type_of_custome_action\":[],\"forbidden_extra_token\":false,\"forbidden_multiple_filename\":true,\"forbidden_suspicious_filename\":true},\"high_risk_action\":\"deny\",\"high_risk_enable_log\":1,\"low_risk_action\":\"continue\",\"low_risk_enable_log\":1,\"medium_risk_action\":\"continue\",\"medium_risk_enable_log\":1,\"state\":\"enabled\"},\"m_http\":{\"detect_config\":{\"warning_http_parse_failed\":false,\"warning_suspicious_http_version\":false},\"high_risk_action\":\"deny\",\"high_risk_enable_log\":1,\"low_risk_action\":\"continue\",\"low_risk_enable_log\":1,\"medium_risk_action\":\"continue\",\"medium_risk_enable_log\":1,\"state\":\"disabled\"},\"m_java\":{\"detect_config\":{\"detect_lookup\":false},\"high_risk_action\":\"deny\",\"high_risk_enable_log\":1,\"low_risk_action\":\"continue\",\"low_risk_enable_log\":1,\"medium_risk_action\":\"continue\",\"medium_risk_enable_log\":1,\"state\":\"enabled\"},\"m_java_unserialize\":{\"high_risk_action\":\"deny\",\"high_risk_enable_log\":1,\"low_risk_action\":\"continue\",\"low_risk_enable_log\":1,\"medium_risk_action\":\"continue\",\"medium_risk_enable_log\":1,\"state\":\"enabled\"},\"m_php_code_injection\":{\"high_risk_action\":\"deny\",\"high_risk_enable_log\":1,\"low_risk_action\":\"continue\",\"low_risk_enable_log\":1,\"medium_risk_action\":\"continue\",\"medium_risk_enable_log\":1,\"state\":\"enabled\"},\"m_php_unserialize\":{\"high_risk_action\":\"deny\",\"high_risk_enable_log\":1,\"low_risk_action\":\"continue\",\"low_risk_enable_log\":1,\"medium_risk_action\":\"continue\",\"medium_risk_enable_log\":1,\"state\":\"enabled\"},\"m_response\":{\"detect_config\":{\"detection_configure\":[\"detect_jsp_code_leak\",\"detect_php_code_leak\",\"detect_webshell\"],\"error_types\":[\"directory indexing\",\"SQL execution error\",\"server exception\"]},\"high_risk_action\":\"deny\",\"high_risk_enable_log\":1,\"low_risk_action\":\"continue\",\"low_risk_enable_log\":1,\"medium_risk_action\":\"continue\",\"medium_risk_enable_log\":1,\"state\":\"disabled\"},\"m_rule\":{\"detect_config\":{\"check_info_leak_by_rsp\":false,\"compatibility_mode\":false,\"detection_configure\":[\"detect_phpinfo\",\"detect_admin_page\",\"detect_backdoor\",\"detect_php_vuln\",\"detect_nginx\",\"detect_dedecms\",\"detect_apache\",\"detect_xml\",\"detect_struts2\",\"detect_java_vuln\",\"detect_php168\",\"detect_wordpress\",\"detect_directory_traversal\",\"detect_iis\",\"detect_gogs_gitea\",\"detect_thinkphp\",\"detect_jenkins\",\"detect_ecshop\",\"detect_nexus\",\"detect_drupal\",\"detect_ghostscript\",\"detect_atlassian\",\"detect_weblogic\",\"detect_coremail\",\"detect_phpcms\",\"detect_spring\",\"detect_southidc\",\"detect_fastjson\",\"detect_tbk_dvr\",\"detect_joomla\",\"detect_ecology_oa\",\"detect_jackson\",\"detect_xstream\",\"detect_activemq\",\"detect_solr\",\"detect_csii\",\"detect_big_ip\",\"detect_apisix\",\"detect_druid\",\"detect_log4j\"],\"info_leak_types\":[\"test file\",\"backup file\",\"code repository\",\"server sensitive file\"],\"rules_config\":{\"disable_ruleid\":[],\"enable_ruleid\":[],\"rules_status\":[]}},\"high_risk_action\":\"deny\",\"high_risk_enable_log\":1,\"low_risk_action\":\"continue\",\"low_risk_enable_log\":1,\"medium_risk_action\":\"continue\",\"medium_risk_enable_log\":1,\"state\":\"enabled\"},\"m_scanner\":{\"detect_config\":{\"language_types\":[\"python_scanner\",\"go_scanner\"],\"other_types\":[\"normal_scanner\"],\"spider_types\":[]},\"high_risk_action\":\"deny\",\"high_risk_enable_log\":1,\"low_risk_action\":\"continue\",\"low_risk_enable_log\":1,\"medium_risk_action\":\"continue\",\"medium_risk_enable_log\":1,\"state\":\"disabled\"},\"m_sqli\":{\"detect_config\":{\"detect_non_injection_sql\":true},\"high_risk_action\":\"deny\",\"high_risk_enable_log\":1,\"low_risk_action\":\"continue\",\"low_risk_enable_log\":1,\"medium_risk_action\":\"continue\",\"medium_risk_enable_log\":1,\"state\":\"enabled\"},\"m_ssrf\":{\"high_risk_action\":\"deny\",\"high_risk_enable_log\":1,\"low_risk_action\":\"continue\",\"low_risk_enable_log\":1,\"medium_risk_action\":\"continue\",\"medium_risk_enable_log\":1,\"state\":\"enabled\"},\"m_ssti\":{\"high_risk_action\":\"deny\",\"high_risk_enable_log\":1,\"low_risk_action\":\"continue\",\"low_risk_enable_log\":1,\"medium_risk_action\":\"continue\",\"medium_risk_enable_log\":1,\"state\":\"enabled\"},\"m_xss\":{\"detect_config\":{\"detect_complete_html\":true},\"high_risk_action\":\"deny\",\"high_risk_enable_log\":1,\"low_risk_action\":\"continue\",\"low_risk_enable_log\":1,\"medium_risk_action\":\"continue\",\"medium_risk_enable_log\":1,\"state\":\"enabled\"}},\"timeout\":{\"threshold\":1000,\"log\":\"enabled\"}}"
	var sc SkynetConfig
	err := json.Unmarshal([]byte(skynetConfigStr), &sc)
	if err != nil {
		return "", err
	}

	var pggOption Options
	db.Where(&Options{Key: constants.PolicyGroupGlobal}).First(&pggOption)

	var pgg PolicyGroup
	err = json.Unmarshal([]byte(pggOption.Value), &pgg)
	if err != nil {
		return "", err
	}

	for module, mode := range pgg {
		switch mode {
		case StrictMode:
			scm := sc.Modules[module]
			scm.HighRiskAction = "deny"
			scm.MediumRiskAction = "deny"
			scm.LowRiskAction = "deny"
			scm.State = "enabled"
			sc.Modules[module] = scm
		case DefaultMode:
			scm := sc.Modules[module]
			scm.HighRiskAction = "deny"
			scm.MediumRiskAction = "continue"
			scm.LowRiskAction = "continue"
			scm.State = "enabled"
			sc.Modules[module] = scm
		case DisableMode:
			scm := sc.Modules[module]
			scm.HighRiskAction = "continue"
			scm.MediumRiskAction = "continue"
			scm.LowRiskAction = "continue"
			scm.State = "disabled"
			sc.Modules[module] = scm
		default:
			return "", errors.New("no such mode")
		}
	}

	scStr, err := json.Marshal(sc)
	if err != nil {
		return "", err
	}

	return string(scStr), nil
}

func initPolicyGroupGlobal() error {
	db := database.GetDB()

	// policyGroupGlobal config, three mode: strict/default/disable
	var policyGroupGlobal = PolicyGroup{
		"m_asp_code_injection": DefaultMode,
		"m_cmd_injection":      DefaultMode,
		"m_csrf":               DefaultMode,
		"m_file_include":       DefaultMode,
		"m_file_upload":        DefaultMode,
		"m_http":               DefaultMode,
		"m_java":               DefaultMode,
		"m_java_unserialize":   DefaultMode,
		"m_php_code_injection": DefaultMode,
		"m_php_unserialize":    DefaultMode,
		"m_response":           DefaultMode,
		"m_rule":               DefaultMode,
		"m_scanner":            DefaultMode,
		"m_sqli":               DefaultMode,
		"m_ssrf":               DefaultMode,
		"m_ssti":               DefaultMode,
		"m_xss":                DefaultMode,
	}

	pggStr, err := json.Marshal(policyGroupGlobal)
	if err != nil {
		return err
	}

	pgg := Options{Key: constants.PolicyGroupGlobal, Value: string(pggStr)}
	db.Clauses(clause.OnConflict{DoNothing: true}).Create(&pgg)

	return nil
}

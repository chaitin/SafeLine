package api

type PolicyRuleAction int

const (
	PolicyRuleActionAllow PolicyRuleAction = iota
	PolicyRuleActionDeny
	PolicyRuleActionMax
)

type Key = string

const (
	KeySrcIP      Key = "src_ip"
	KeyURI        Key = "uri"
	KeyURINoQuery Key = "uri_no_query"
	KeyHost       Key = "host"
	KeyMethod     Key = "method"
	KeyReqHeader  Key = "req_header"
	KeyReqBody    Key = "req_body"
	KeyGetParam   Key = "get_param"
	KeyPostParam  Key = "post_param"
)

type Op = string

const (
	OpEq       Op = "eq"         // equal
	OpNotEq    Op = "not_eq"     // not equal
	OpMatch    Op = "match"      // match
	OpCIDR     Op = "cidr"       // cidr
	OpHas      Op = "has"        // has
	OpNotHas   Op = "not_has"    // not has
	OpPrefix   Op = "prefix"     // prefix
	OpRe       Op = "re"         // regex
	OpIn       Op = "in"         // in
	OpNotIn    Op = "not_in"     // not in
	OpNotCIDR  Op = "not_cidr"   // not cidr
	OpExist    Op = "exist"      // exist
	OpNotExist Op = "not_exist"  // not exist
	OpGeoEq    Op = "geo_eq"     // geo equal
	OpGeoNotEq Op = "geo_not_eq" // geo not equal
)

type Pattern struct {
	K    Key      `json:"k"`
	Op   Op       `json:"op"`
	V    []string `json:"v"`
	SubK string   `json:"sub_k"`
}

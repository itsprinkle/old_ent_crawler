package creditd

import (
	"reflect"

	"gsxt/credit"
)

// 基本信息
type BaseInfo struct {
	RegNo       string `json:"reg_no" struct:"注册号"`                 // 注册号
	CreditNo    string `json:"credit_no" struct:"信用代码"`             // 统一社会信用代码
	Name        string `json:"name" struct:"名称"`                    // 名称
	Type        string `json:"type" struct:"类型"`                    // 类型
	Formation   string `json:"formation" struct:"组成形式"`             // 组成形式
	State       string `json:"state" struct:"状态"`                   // 登记状态
	LegRep      string `json:"leg_rep" struct:"合伙人,经营者,负责人,代表,投资人"` // 法人
	Scope       string `json:"scope" struct:"范围"`                   // 经验范围
	Address     string `json:"address" struct:"住所,场所"`              // 住所
	RegCapi     string `json:"reg_capi" struct:"注册资本,出资总额,注册资金"`    // 注册资金
	RegOrg      string `json:"reg_org" struct:"登记机关"`               // 登记机关
	TermStartAt string `json:"term_start_at" struct:"期限自"`          // 经营期限自
	TermEndAt   string `json:"term_end_at" struct:"期限至"`            // 经营期限至
	CheckAt     string `json:"check_at" struct:"核准日期"`              // 核准日期
	StartAt     string `json:"start_at" struct:"成立日期,注册日期"`         // 成立日期
	EndAt       string `json:"end_at" struct:"注销日期"`                // 吊销日期
	RevokedAt   string `json:"revoked_at" struct:"吊销日期"`
}

// 股东信息
// 投资人信息
// 发起人
// 主管部门（出资人）信息
type StockholderInfo struct {
	Type     string `json:"type" struct:"股东（发起人）类型,股东类型,出资人类型,发起人类型,出资方式,投资人类型"` // 股东类型
	Name     string `json:"name" struct:"股东,出资人,姓名,发起人,投资人"`                     // 股东
	CertType string `json:"cert_type" struct:"证件类型,证照类型"`                        // 证件类型
	CertNo   string `json:"cert_no" struct:"证件号码,证照号码"`                          // 证件号码
}

// 变更信息
type ChangeInfo struct {
	Item   string `json:"item" struct:"变更事项,事项"`         // 变更事项
	Before string `json:"before" struct:"变更前,前"`         // 变更前内容
	After  string `json:"after" struct:"变更后,后"`          // 变更后内容
	Date   string `json:"change_at" struct:"变更日期,时间,日期"` // 变更日期
}

// 出资信息
type InvestmentInfo struct {
	Name           string `json:"name" struct:"股东,发起人,投资人名称"`                    // 股东
	TotalSubAmount string `json:"total_sub_amount,omitempty" struct:"认缴额,累计认缴额"` // 认缴额
	TotalActAmount string `json:"total_act_amount,omitempty" struct:"实缴额,累计实缴额"` // 实缴额
	SubType        string `json:"sub_type" struct:"认缴出资方式"`                      // 认缴出资方式
	SubAmount      string `json:"sub_amount" struct:"认缴出资额"`                     // 认缴出资额
	SubAt          string `json:"sub_at" struct:"认缴出资日期,认缴出资时间"`                 // 认缴出资日期
	ActType        string `json:"act_type" struct:"实缴出资方式"`                      // 实缴出资方式
	ActAmount      string `json:"act_amount" struct:"实缴出资额"`                     // 实缴出资额
	ActAt          string `json:"act_at" struct:"实缴出资日期,实缴出资时间"`                 // 实缴出资日期
}

// 成员信息
// 参加经营的家庭成员姓名
type MemberInfo struct {
	Name     string `json:"name" struct:"姓名"`     // 姓名
	Position string `json:"position" struct:"职务"` // 职务
}

// 分支机构
type BranchInfo struct {
	RegNo  string `json:"reg_no" struct:"注册号,信用代码"` // 注册号/统一社会信用代码
	Name   string `json:"name" struct:"名称"`         // 名称
	RegOrg string `json:"reg_org" struct:"登记机关"`    // 登记机关
}

// 股权出质登记信息
type EquityPledgeInfo struct {
	Name         string `json:"name" struct:"名称"`
	No           string `json:"equity_no" struct:"登记编号"` // 登记编号
	Pledgor      string `json:"pledgor" struct:"出质人"`
	PledgorNo    string `json:"pledgor_no" struct:"证照/证件号码"`    // 证照/证件号码
	EquityAmount string `json:"pledgor_strand" struct:"出质股权数额"` // 出质股权数额
	Pledgee      string `json:"pledgee" struct:"质权人"`           // 质权人
	PledgeeNo    string `json:"pledgee_no" struct:"证照/证件号码"`    // 证照/证件号码
	State        string `json:"state" struct:"状态"`              // 状态
	Date         string `json:"pledge_at" struct:"登记日期"`        // 股权出质设立登记日期
}

// 抵押权人概况
type MortgageeInfo struct {
	Name     string `json:"name" struct:"抵押权人名称"`       // 抵押权人名称
	CertType string `json:"cert_type" struct:"证件类型,类型"` // 证件类型
	CertNo   string `json:"cert_no" struct:"证件号码,号码"`   // 证件号码
}

// 被担保债权概况
type DebtSecuredInfo struct {
	Kind     string `json:"kind" struct:"种类"`      // 种类
	Amount   string `json:"amount" struct:"数额"`    // 数额
	Scope    string `json:"scope" struct:"担保的范围"`  // 担保的范围
	DebtTerm string `json:"debt_term" struct:"期限"` // 债务人履行债务的期限
	Remark   string `json:"remark" struct:"备注"`    // 备注
}

// 抵押物概况
type CollateralInfo struct {
	Name   string `json:"name" struct:"名称"`                   // 名称
	Owner  string `json:"ownership" struct:"所有权归属,所有权或使用权归属"` // 所有权归属
	Status string `json:"status" struct:"数量,质量,状况"`           // 数量、质量、状况、所在地等情况
	Remark string `json:"remark" struct:"备注"`                 // 备注
}

// 动产抵押登记信息
type MortgageRegInfo struct {
	No          string `json:"no" struct:"登记编号"`         // 登记编号
	RegOrg      string `json:"reg_org" struct:"登记机关"`    // 登记机关
	RegDate     string `json:"reg_at" struct:"登记日期"`     // 登记日期
	DebtType    string `json:"debut_type" struct:"种类"`   // 被担保债权种类
	DebtAmount  string `json:"debt_amount" struct:"数额"`  // 被担保债权数额
	DebtTerm    string `json:"debt_term" struct:"期限"`    // 债务人履行债务的期限
	SecureScope string `json:"secure_scope" struct:"范围"` // 担保的范围
	State       string `json:"state" struct:"状态"`        // 状态
	Remark      string `json:"remark" struct:"备注"`       // 备注
}

// 动产抵押信息
type MortgageInfo struct {
	MortgageRegInfo MortgageRegInfo  `json:"mortgage_reg"` // 动产抵押登记信息
	MortgageeInfos  []MortgageeInfo  `json:"mortgagee"`    // 抵押权人概况
	DebtSecuredInfo DebtSecuredInfo  `json:"debt_secured"` // 被担保债权概况
	CollateralInfos []CollateralInfo `json:"collateral"`   // 抵押物概况
}

// 行政处罚信息
type PunishInfo struct {
	No      string `json:"paper_no" struct:"书文号,行政处罚决定书文号"`         // 行政处罚决定书文号
	Name    string `json:"name" struct:"名称" match:"full"`           // 名称
	RegNo   string `json:"reg_no" struct:"注册号,信用代码"`                // 注册号
	LegRep  string `json:"leg_rep" struct:"姓名,负责人,代表人"`             // 法定代表人（负责人）姓名
	Type    string `json:"type" struct:"违法行为类型"`                    // 违法行为类型
	Content string `json:"content" struct:"处罚种类,行政处罚内容"`            // 行政处罚内容
	OrgName string `json:"org_name" struct:"处罚机关,决定机关,出行政处罚决定机关名称"` // 作出行政处罚决定机关名称
	Date    string `json:"dec_at" struct:"处罚决定书签发日期,处罚决定日期"`        // 作出行政处罚决定日期
	Detail  string `json:"detail" struct:"决定书地址"`                   // 行政处罚决定书
	Remark  string `json:"remark" struct:"备注"`                      // 备注
}

// 经营异常
type AbnormalInfo struct {
	AddCause    string `json:"add_cause" struct:"列入经营异常名录原因,标记经营异常状态原因"`    // 列入经营异常名录原因
	AddDate     string `json:"add_at" struct:"列入日期,标记日期"`                   // 列入日期
	RemoveCause string `json:"remove_cause" struct:"移出经营异常名录原因,恢复正常记载状态原因"` // 移出经营异常名录原因
	RemoveDate  string `json:"remove_at" struct:"移出日期,恢复日期"`                // 移出日期
	DecOrg      string `json:"dec_org" struct:"作出决定机关,机关,作出决定机关(列入)"`       // 作出决定机关
}

// 抽查检查信息
type SpotCheckInfo struct {
	CheckOrg string `json:"check_org" struct:"检查实施机关"` // 检查实施机关
	Type     string `json:"type" struct:"类型"`          // 类型
	Date     string `json:"check_at" struct:"日期"`      // 日期
	Result   string `json:"result" struct:"结果"`        // 结果
}

// 工商公示信息
type BusinessInfo struct {
	BaseInfo          BaseInfo           `json:"base"`           // 基本信息
	ChangeInfos       []ChangeInfo       `json:"changes"`        // 变更信息
	StockholderInfos  []StockholderInfo  `json:"stockholders"`   // 股东信息
	InvestmentInfos   []InvestmentInfo   `json:"investment"`     // 股东及出资信息
	MemberInfos       []MemberInfo       `json:"members"`        // 成员信息
	BranchInfos       []BranchInfo       `json:"branchs"`        // 分支机构
	EquityPledgeInfos []EquityPledgeInfo `json:"equity_pledges"` // 股权出质信息
	MortgageInfos     []MortgageInfo     `json:"mortgages"`      // 动产抵押信息
	PunishInfos       []PunishInfo       `json:"punishs"`        // 行政处罚信息
	AbnormalInfos     []AbnormalInfo     `json:"abnormals"`      // 经营异常
	SpotCheckInfos    []SpotCheckInfo    `json:"spot_checks"`    // 抽查检查信息
}

// 行政许可信息
type LicenseInfo struct {
	No      string `json:"no,omitempty" struct:"编号"`         // 许可文件编号
	Name    string `json:"name" struct:"名称"`                 // 许可文件名称
	From    string `json:"start_at,omitempty" struct:"有效期自"` // 有效期自
	To      string `json:"end_at" struct:"有效期至"`             // 有效期至
	Org     string `json:"org,omitempty" struct:"机关"`        // 许可机关
	Content string `json:"content,omitempty" struct:"内容"`    // 许可内容
	State   string `json:"state,omitempty" struct:"状态"`      // 状态
}

// 知识产权出质登记信息
type IntellectualInfo struct {
	No      string `json:"no" struct:"号"`        // 注册号/统一社会信用代码
	Name    string `json:"name" struct:"名称"`     // 名称
	Kind    string `json:"kind" struct:"种类,类别"`  // 种类
	Pledgor string `json:"pledgor" struct:"出质人"` // 出质人
	Pledgee string `json:"pledgee" struct:"质权人"` // 质权人
	Term    string `json:"term" struct:"期限"`     // 质权登记期限
	State   string `json:"state" struct:"状态"`    // 状态
}

// 股权变更信息
type StockChangeInfo struct {
	Stockholder string `json:"stockholder" struct:"股东,发起人"` // 股东
	Before      string `json:"before" struct:"变更前"`         // 变更前股权比例
	After       string `json:"after" struct:"变更后"`          // 变更后股权比例
	Date        string `json:"change_at" struct:"变更日期"`     // 股权变更日期
}

// 企业基本信息
type EntBaseInfo struct {
	RegNo      string `json:"reg_no" struct:"注册号"`         // 注册号
	CreditCode string `json:"credit_no" struct:"信用代码"`     // 统一社会信用代码
	Name       string `json:"name" struct:"名称"`            // 名称
	Telphone   string `json:"telphone" struct:"电话"`        // 联系电话
	Postcode   string `json:"postcode" struct:"邮政编码"`      // 邮编
	Email      string `json:"email" struct:"电子邮箱"`         // 邮箱
	Address    string `json:"address" struct:"住所,场所,地址"`   // 住所
	State      string `json:"state" struct:"状态"`           // 登记状态
	EmployNum  string `json:"employ_num" struct:"从业人数,人数"` // 雇员数量
	IsStock    string `json:"is_stock" struct:"股权转让"`      // 是否发生股权转让
	IsWebsite  string `json:"is_website" struct:"网店,网站"`   // 是否有网站
	IsInvest   string `json:"is_invest" struct:"购买,投资"`    // 是否有投资成立企业
	IsGuar     string `json:"is_guarantee" struct:"对外担保"`
	LegRep     string `json:"leg_rep" struct:"合伙人,经营者,负责人,代表,投资人"` // 法人
	Capi       string `json:"amount" struct:"资金数额"`                // 资金数额
	Type       string `json:"type" struct:"主体类型"`                  // 主体类型
	Relation   string `json:"relationship" struct:"隶属关系"`          // 隶属关系
}

// 生产经营
type OperationInfo struct {
	TotalAsset    string `json:"total_asset" struct:"资产总额"`
	TotalTax      string `json:"total_tax" struct:"纳税金额,纳税总额"`
	TotalDebt     string `json:"total_debt" struct:"负债总额"`
	TotalTurnover string `json:"total_turnover" struct:"营业总收入,营业收入,总收入"`
	MainIncome    string `json:"main_income" struct:"销售额,销售总额,营业额,主营,主营业务收入"`
	TotalProfit   string `json:"profit" struct:"利润总额,盈余总额"`
	NetProfit     string `json:"net_profit" struct:"净利润"`
	TotalEquity   string `json:"total_equity" struct:"权益合计"`        // 所有者权益合计
	FinancialLoan string `json:"financial_loan" struct:"金融贷款"`      // 金融贷款
	FundSubsidy   string `json:"fund_subsidy" struct:"获得政府扶持资金,补助"` // 获得政府扶持资金、补助
}

// 网站或者网店信息
type WebsiteInfo struct {
	Name string `json:"name" struct:"名称"` // 名称
	Type string `json:"type" struct:"类型"` // 类型
	Url  string `json:"url" struct:"网址"`  // 网址
}

// 对外投资信息
type InvEntInfo struct {
	Name  string `json:"name" struct:"名称"`         // 投资设立企业或购买股权企业名称
	RegNo string `json:"reg_no" struct:"注册号,信用代码"` // 注册号/统一社会信用代码
}

// 保证担保
type GuaranteeInfo struct {
	Creditor   string `json:"creditor" struct:"债权人"`       // 债权人
	Debtor     string `json:"debtor" struct:"债务人"`         // 债务人
	DebtKind   string `json:"debt_kind" struct:"主债权种类"`    // 主债权种类
	DebtAmount string `json:"debt_amount" struct:"主债权数额"`  // 主债权数额
	DebtTerm   string `json:"debt_term" struct:"履行债务的期限"`  // 履行债务的期限
	GuarTerm   string `json:"guar_term" struct:"保证的期间"`    // 保证的期间
	GuarType   string `json:"guar_type" struct:"方式,保证的方式"` // 保证方式
	GuarRange  string `json:"guar_range" struct:"范围"`      // 保证范围
}

// 年报信息
type ReportInfo struct {
	Year             string            `json:"year" struct:"年度"`      // 送报年度
	ReportAt         string            `json:"report_at" struct:"日期"` // 发布日期
	EntBaseInfo      EntBaseInfo       `json:"ent_base"`              // 企业基本信息
	OperationInfo    OperationInfo     `json:"operation"`             // 生产经营
	WebsiteInfos     []WebsiteInfo     `json:"websites"`              // 网站
	LicenseInfos     []LicenseInfo     `json:"licenses"`              // 行政许可
	BranchInfos      []BranchInfo      `json:"branchs"`               // 分支机构
	InvEntInfos      []InvEntInfo      `json:"inv_ents"`              // 对外投资企业
	GuaranteeInfos   []GuaranteeInfo   `json:"guarantees"`            // 保证担保
	InvestmentInfos  []InvestmentInfo  `json:"investment"`            // 股东及出资信息
	StockholderInfos []StockholderInfo `json:"stockholders"`          // 股东信息
	StockChangeInfos []StockChangeInfo `json:"stock_changes"`         // 股权变更
	ChangeInfos      []ChangeInfo      `json:"changes"`               // 变更事项
}

// 企业公示信息
type EnterpriseInfo struct {
	ReportInfos       []ReportInfo       `json:"reports"`       // 企业年报
	InvestmentInfos   []InvestmentInfo   `json:"investment"`    // 股东出资信息
	ChangeInfos       []ChangeInfo       `json:"changes"`       // 变更事项
	StockChangeInfos  []StockChangeInfo  `json:"stock_changes"` // 股权变更
	LicenseInfos      []LicenseInfo      `json:"licenses"`      // 行政许可
	IntellectualInfos []IntellectualInfo `json:"intellectuals"` // 知识产权出质
	PunishInfos       []PunishInfo       `json:"punishs"`       // 行政处罚
}

type InfoV1 struct {
	Business   BusinessInfo   `json:"bus"`
	Enterprise EnterpriseInfo `json:"ent"`
}

func s2s(in interface{}, out interface{}, data map[string]string) error {
	inValue := reflect.ValueOf(in).Elem()
	inType := inValue.Type()
	outValue := reflect.ValueOf(out).Elem()
	outType := outValue.Type()
	if inType.Kind() != reflect.Struct || outType.Kind() != reflect.Struct {
		panic("expected pointer to struct")
	}

	for i := 0; i < inType.NumField(); i++ {
		inFieldType := inType.Field(i)
		name := inFieldType.Name
		if len(data) > 0 {
			if val, ok := data[name]; ok {
				name = val
			}
		}
		if _, ok := outType.FieldByName(name); ok {
			outValue.FieldByName(name).Set(inValue.Field(i))
		} else {
			//fmt.Printf("%+v\n", in)
			//println("not found name", name)
		}
	}
	return nil
}

func ToV1(infoV2 credit.InfoV2) (infoV1 InfoV1) {
	// 基本信息
	s2s(&infoV2.Business.Base, &infoV1.Business.BaseInfo, map[string]string{
		"CreditCode":   "CreditNo",
		"OpFrom":       "TermStartAt",
		"OpTo":         "TermEndAt",
		"DateReg":      "StartAt",
		"DateApproved": "CheckAt",
		"DateCanceled": "EndAt",
		"DateRevoked":  "RevokedAt",
	})
	// 投资信息
	for _, inv2 := range infoV2.Business.Investors {
		var inv InvestmentInfo
		var stock StockholderInfo
		s2s(&inv2, &stock, nil)
		inv.Name = inv2.Name
		inv.TotalActAmount = inv2.ActCapi
		inv.TotalSubAmount = inv2.SubCapi
		for _, sub := range inv2.Subs {
			inv.SubAt = sub.Date
			inv.SubType = sub.Type
			inv.SubAmount = sub.Capi
		}
		for _, act := range inv2.Acts {
			inv.ActAt = act.Date
			inv.ActType = act.Type
			inv.ActAmount = act.Capi
		}
		infoV1.Business.StockholderInfos = append(infoV1.Business.StockholderInfos, stock)
		infoV1.Business.InvestmentInfos = append(infoV1.Business.InvestmentInfos, inv)
	}
	for _, change2 := range infoV2.Business.Changes {
		var change ChangeInfo
		s2s(&change2, &change, nil)
		infoV1.Business.ChangeInfos = append(infoV1.Business.ChangeInfos, change)
	}
	for _, mem2 := range infoV2.Business.Members {
		var mem MemberInfo
		s2s(&mem2, &mem, nil)
		infoV1.Business.MemberInfos = append(infoV1.Business.MemberInfos, mem)
	}
	for _, branch2 := range infoV2.Business.Branchs {
		var branch BranchInfo
		s2s(&branch2, &branch, nil)
		infoV1.Business.BranchInfos = append(infoV1.Business.BranchInfos, branch)
	}
	for _, ple2 := range infoV2.Business.Pledges {
		var ple EquityPledgeInfo
		s2s(&ple2, &ple, nil)
		infoV1.Business.EquityPledgeInfos = append(infoV1.Business.EquityPledgeInfos, ple)
	}
	for _, pun2 := range infoV2.Business.Punishs {
		var pun PunishInfo
		s2s(&pun2, &pun, map[string]string{
			"DecOrg": "OrgName",
		})
		infoV1.Business.PunishInfos = append(infoV1.Business.PunishInfos, pun)
	}
	for _, abn2 := range infoV2.Business.Abnormals {
		var abn AbnormalInfo
		s2s(&abn2, &abn, nil)
		infoV1.Business.AbnormalInfos = append(infoV1.Business.AbnormalInfos, abn)
	}
	for _, spot2 := range infoV2.Business.SpotChecks {
		var spot SpotCheckInfo
		s2s(&spot2, &spot, nil)
		infoV1.Business.SpotCheckInfos = append(infoV1.Business.SpotCheckInfos, spot)
	}
	for _, mort2 := range infoV2.Business.Mortgages {
		var mort MortgageInfo
		s2s(&mort2, &mort.MortgageRegInfo, nil)
		s2s(&mort2.Obligee, &mort.DebtSecuredInfo, nil)
		for _, pawn := range mort2.Pawns {
			var coll CollateralInfo
			s2s(&pawn, &coll, nil)
			mort.CollateralInfos = append(mort.CollateralInfos, coll)
		}
		for _, m2 := range mort2.Mortgagers {
			var m MortgageeInfo
			s2s(&m2, &m, nil)
			mort.MortgageeInfos = append(mort.MortgageeInfos, m)
		}
		infoV1.Business.MortgageInfos = append(infoV1.Business.MortgageInfos, mort)
	}

	entV2 := &infoV2.Enterprise
	entV1 := &infoV1.Enterprise
	for _, report2 := range entV2.Reports {
		var report ReportInfo
		report.Year = report2.Year
		report.ReportAt = report2.Date
		s2s(&report2.General, &report.EntBaseInfo, nil)
		s2s(&report2.Operation, &report.OperationInfo, nil)
		for _, w2 := range report2.Websites {
			var w WebsiteInfo
			s2s(&w2, &w, nil)
			report.WebsiteInfos = append(report.WebsiteInfos, w)
		}
		for _, lic2 := range report2.Licenses {
			var lic LicenseInfo
			s2s(&lic2, &lic, map[string]string{
				"StartAt": "From",
				"EndAt":   "To",
			})
			report.LicenseInfos = append(report.LicenseInfos, lic)
		}
		for _, branch2 := range report2.Branchs {
			var branch BranchInfo
			s2s(&branch2, &branch, nil)
			report.BranchInfos = append(report.BranchInfos, branch)
		}
		for _, inv2 := range report2.InvEnts {
			var inv InvEntInfo
			s2s(&inv2, &inv, nil)
			report.InvEntInfos = append(report.InvEntInfos, inv)
		}
		for _, gua2 := range report2.Guarantees {
			var gua GuaranteeInfo
			s2s(&gua2, &gua, nil)
			report.GuaranteeInfos = append(report.GuaranteeInfos, gua)
		}
		for _, inv2 := range report2.Investors {
			var inv InvestmentInfo
			//var stock StockholderInfo
			//s2s(&inv2, &stock, nil)
			inv.Name = inv2.Name
			inv.TotalActAmount = inv2.ActCapi
			inv.TotalSubAmount = inv2.SubCapi
			for _, sub := range inv2.Subs {
				inv.SubAt = sub.Date
				inv.SubType = sub.Type
				inv.SubAmount = sub.Capi
			}
			for _, act := range inv2.Acts {
				inv.ActAt = act.Date
				inv.ActType = act.Type
				inv.ActAmount = act.Capi
			}
			//report.StockholderInfos = append(report.StockholderInfos, stock)
			report.InvestmentInfos = append(report.InvestmentInfos, inv)
		}
		for _, stock2 := range report2.StockChanges {
			var stock StockChangeInfo
			s2s(&stock2, &stock, nil)
			report.StockChangeInfos = append(report.StockChangeInfos, stock)
		}
		for _, change2 := range report2.Changes {
			var change ChangeInfo
			s2s(&change2, &change, nil)
			report.ChangeInfos = append(report.ChangeInfos, change)
		}
		entV1.ReportInfos = append(entV1.ReportInfos, report)
	}

	for _, inv2 := range entV2.Investors {
		var inv InvestmentInfo
		inv.Name = inv2.Name
		inv.TotalActAmount = inv2.ActCapi
		inv.TotalSubAmount = inv2.SubCapi
		for _, sub := range inv2.Subs {
			inv.SubAt = sub.Date
			inv.SubType = sub.Type
			inv.SubAmount = sub.Capi
		}
		for _, act := range inv2.Acts {
			inv.ActAt = act.Date
			inv.ActType = act.Type
			inv.ActAmount = act.Capi
		}
		entV1.InvestmentInfos = append(entV1.InvestmentInfos, inv)
	}
	for _, change2 := range entV2.Changes {
		var change ChangeInfo
		s2s(&change2, &change, nil)
		entV1.ChangeInfos = append(entV1.ChangeInfos, change)
	}
	for _, stock2 := range entV2.StockChanges {
		var stock StockChangeInfo
		s2s(&stock2, &stock, nil)
		entV1.StockChangeInfos = append(entV1.StockChangeInfos, stock)
	}
	for _, lic2 := range entV2.Licenses {
		var lic LicenseInfo
		s2s(&lic2, &lic, map[string]string{
			"StartAt": "From",
			"EndAt":   "To",
		})
		entV1.LicenseInfos = append(entV1.LicenseInfos, lic)
	}
	for _, intell2 := range entV2.Intells {
		var intell IntellectualInfo
		s2s(&intell2, &intell, nil)
		entV1.IntellectualInfos = append(entV1.IntellectualInfos, intell)
	}
	for _, pun2 := range entV2.Punishs {
		var pun PunishInfo
		s2s(&pun2, &pun, map[string]string{
			"DecOrg": "OrgName",
		})
		entV1.PunishInfos = append(entV1.PunishInfos, pun)
	}
	return
}

func ToV2(infoV1 InfoV1) (infoV2 credit.InfoV2) {
	// 基本信息
	s2s(&infoV1.Business.BaseInfo, &infoV2.Business.Base, map[string]string{
		"CreditNo":    "CreditCode",
		"TermStartAt": "OpFrom",
		"TermEndAt":   "OpTo",
		"StartAt":     "DateReg",
		"CheckAt":     "DateApproved",
		"EndAt":       "DateCanceled",
		"RevokedAt":   "DateRevoked",
	})
	// 投资信息			 infoV1.Business.InvestmentInfos
	for _, inv1 := range infoV1.Business.StockholderInfos {
		var inv credit.InvestorInfo
		s2s(&inv1, &inv, nil)
		for _, im := range infoV1.Business.InvestmentInfos {
			if inv.Name == im.Name {
				inv.ActCapi = im.TotalActAmount
				inv.SubCapi = im.TotalSubAmount
				var act credit.CapiInfo
				act.Capi = im.ActAmount
				act.Date = im.ActAt
				act.Type = im.ActType
				var sub credit.CapiInfo
				sub.Capi = im.SubAmount
				sub.Date = im.SubAt
				sub.Type = im.SubType
				inv.Acts = append(inv.Acts, act)
				inv.Subs = append(inv.Subs, sub)
			}
		}
		infoV2.Business.Investors = append(infoV2.Business.Investors, inv)
	}
	for _, change1 := range infoV1.Business.ChangeInfos {
		var change credit.ChangeInfo
		s2s(&change1, &change, nil)
		infoV2.Business.Changes = append(infoV2.Business.Changes, change)
	}

	for _, mem1 := range infoV1.Business.MemberInfos {
		var mem credit.MemberInfo
		s2s(&mem1, &mem, nil)
		infoV2.Business.Members = append(infoV2.Business.Members, mem)
	}
	for _, branch1 := range infoV1.Business.BranchInfos {
		var branch credit.BranchInfo
		s2s(&branch1, &branch, nil)
		infoV2.Business.Branchs = append(infoV2.Business.Branchs, branch)

	}
	for _, ple1 := range infoV1.Business.EquityPledgeInfos {
		var ple credit.PledgeInfo
		s2s(&ple1, &ple, nil)
		infoV2.Business.Pledges = append(infoV2.Business.Pledges, ple)
	}
	for _, pun1 := range infoV1.Business.PunishInfos {
		var pun credit.PunishInfo
		s2s(&pun1, &pun, map[string]string{
			"OrgName": "DecOrg",
		})
		infoV2.Business.Punishs = append(infoV2.Business.Punishs, pun)
	}
	for _, abn1 := range infoV1.Business.AbnormalInfos {
		var abn credit.AbnormalInfo
		s2s(&abn1, &abn, nil)
		infoV2.Business.Abnormals = append(infoV2.Business.Abnormals, abn)
	}
	for _, spot1 := range infoV1.Business.SpotCheckInfos {
		var spot credit.SpotCheckInfo
		s2s(&spot1, &spot, nil)
		infoV2.Business.SpotChecks = append(infoV2.Business.SpotChecks, spot)
	}
	for _, mort1 := range infoV1.Business.MortgageInfos {
		var mort credit.MortgageInfo
		s2s(&mort1.MortgageRegInfo, &mort, nil)
		s2s(&mort1.DebtSecuredInfo, &mort.Obligee, nil)
		for _, pawn := range mort1.CollateralInfos {
			var coll credit.PawnInfo
			s2s(&pawn, &coll, nil)
			mort.Pawns = append(mort.Pawns, coll)
		}
		for _, m1 := range mort1.MortgageeInfos {
			var m credit.MortgagerInfo
			s2s(&m1, &m, nil)
			mort.Mortgagers = append(mort.Mortgagers, m)
		}
		infoV2.Business.Mortgages = append(infoV2.Business.Mortgages, mort)
	}
	entV2 := &infoV2.Enterprise
	entV1 := &infoV1.Enterprise
	for _, report1 := range entV1.ReportInfos {
		var report credit.ReportInfo

		report.Year = report1.Year
		report.Date = report1.ReportAt

		s2s(&report1.EntBaseInfo, &report.General, nil)
		s2s(&report1.OperationInfo, &report.Operation, nil)
		for _, w1 := range report1.WebsiteInfos {
			var w credit.WebsiteInfo
			s2s(&w1, &w, nil)
			report.Websites = append(report.Websites, w)
		}
		for _, lic1 := range report1.LicenseInfos {
			var lic credit.LicenseInfo
			s2s(&lic1, &lic, map[string]string{
				"From": "StartAt",
				"To":   "EndAt",
			})
			report.Licenses = append(report.Licenses, lic)
		}
		for _, branch1 := range report1.BranchInfos {
			var branch credit.BranchInfo
			s2s(&branch1, &branch, nil)
			report.Branchs = append(report.Branchs, branch)
		}
		for _, inv1 := range report1.InvEntInfos {
			var inv credit.InvEntInfo
			s2s(&inv1, &inv, nil)
			report.InvEnts = append(report.InvEnts, inv)
		}
		for _, gua1 := range report1.GuaranteeInfos {
			var gua credit.GuaranteeInfo
			s2s(&gua1, &gua, nil)
			report.Guarantees = append(report.Guarantees, gua)
		}
		for _, inv1 := range report1.InvestmentInfos {
			var inv credit.InvestorInfo

			inv.Name = inv1.Name
			inv.ActCapi = inv1.TotalActAmount
			inv.SubCapi = inv1.TotalSubAmount
			var act credit.CapiInfo
			var sub credit.CapiInfo
			act.Capi = inv1.ActAmount
			act.Date = inv1.ActAt
			act.Type = inv1.ActType
			sub.Capi = inv1.SubAmount
			sub.Date = inv1.SubAt
			sub.Type = inv1.SubType
			inv.Acts = append(inv.Acts, act)
			inv.Subs = append(inv.Subs, act)

			report.Investors = append(report.Investors, inv)
		}
		for _, stock1 := range report1.StockChangeInfos {
			var stock credit.StockChangeInfo
			s2s(&stock1, &stock, nil)
			report.StockChanges = append(report.StockChanges, stock)
		}
		for _, change1 := range report1.ChangeInfos {
			var change credit.ChangeInfo
			s2s(&change1, &change, nil)
			report.Changes = append(report.Changes, change)
		}
		entV2.Reports = append(entV2.Reports, report)
	}
	for _, inv1 := range entV1.InvestmentInfos {
		var inv credit.InvestorInfo
		inv.Name = inv1.Name
		inv.ActCapi = inv1.TotalActAmount
		inv.SubCapi = inv1.TotalSubAmount
		var act credit.CapiInfo
		var sub credit.CapiInfo
		act.Capi = inv1.ActAmount
		act.Date = inv1.ActAt
		act.Type = inv1.ActType
		sub.Capi = inv1.SubAmount
		sub.Date = inv1.SubAt
		sub.Type = inv1.SubType
		inv.Acts = append(inv.Acts, act)
		inv.Subs = append(inv.Subs, sub)
		entV2.Investors = append(entV2.Investors, inv)
	}
	for _, change1 := range entV1.ChangeInfos {
		var change credit.ChangeInfo
		s2s(&change1, &change, nil)
		entV2.Changes = append(entV2.Changes, change)
	}
	for _, stock1 := range entV1.StockChangeInfos {
		var stock credit.StockChangeInfo
		s2s(&stock1, &stock, nil)
		entV2.StockChanges = append(entV2.StockChanges, stock)
	}
	for _, lic1 := range entV1.LicenseInfos {
		var lic credit.LicenseInfo
		s2s(&lic1, &lic, map[string]string{
			"From": "StartAt",
			"To":   "EndAt",
		})
		entV2.Licenses = append(entV2.Licenses, lic)
	}
	for _, intell1 := range entV1.IntellectualInfos {
		var intell credit.IntellInfo
		s2s(&intell1, &intell, nil)
		entV2.Intells = append(entV2.Intells, intell)
	}
	for _, pun1 := range entV1.PunishInfos {
		var pun credit.PunishInfo
		s2s(&pun1, &pun, map[string]string{
			"OrgName": "DecOrg",
		})
		entV2.Punishs = append(entV2.Punishs, pun)
	}
	return
}

package models

const (
	RuleTypeSubscribe = iota + 1
	RuleTypeEventRelationship
	RuleLogicIntersection = iota - 1
	RuleLogicConcatenation
	RuleOpEqual = iota - 3
	RuleOpGreatThan
	RuleOpLessThan
	RuleOpInclude
	RuleOpRegex
)

// 规则：（关联：规则、接收者、维护项）
type MRule struct {
	BaseModel
	Name     string  `json:"name" gorm:"column:col_name;not null;comment:规则名称"`
	// RType    int     `json:"rtype" gorm:"column:col_rtype;not null;default:1;comment:定义规则类型：告警指派规则或事件关联性规则等等"`
	StartAt  int64   `json:"start_at" gorm:"column:col_start_at;not null;comment:规则生效时间段开始秒级时间戳"`
	EndAt    int64   `json:"end_at" gorm:"column:col_end_at;not null;comment:规则生效时间段结束秒级时间戳"`
	NextID   uint    `json:"next_id" gorm:"column:col_next_id;comment:规则链条下一节点id"`
	Nexts    []MRule `json:"nexts" gorm:"foreignKey:NextID"`
	Logic    int     `json:"logic" gorm:"column:col_logic;not null;default:1;comment:规则内容中每项之间的逻辑关系（且、或）"`
	Operator int     `json:"operator" gorm:"column:col_operator;not null;default:1;comment:规则内容中每一项大比对符号"`
	Key01    string  `json:"key01" gorm:"column:col_key01;comment:预留字段"`
	Key02    string  `json:"key02" gorm:"column:col_key02;comment:预留字段"`
	Key03    string  `json:"key03" gorm:"column:col_key03;comment:预留字段"`
	Key04    string  `json:"key04" gorm:"column:col_key04;comment:预留字段"`
	Key05    string  `json:"key05" gorm:"column:col_key05;comment:预留字段"`
	Key06    string  `json:"key06" gorm:"column:col_key06;comment:预留字段"`
	Key07    string  `json:"key07" gorm:"column:col_key07;comment:预留字段"`
	Key08    string  `json:"key08" gorm:"column:col_key08;comment:预留字段"`
	Key09    string  `json:"key09" gorm:"column:col_key09;comment:预留字段"`
	Key10    string  `json:"key10" gorm:"column:col_key10;comment:预留字段"`
	Key11    string  `json:"key11" gorm:"column:col_key11;comment:预留字段"`
	Key12    string  `json:"key12" gorm:"column:col_key12;comment:预留字段"`
	Key13    string  `json:"key13" gorm:"column:col_key13;comment:预留字段"`
	Key14    string  `json:"key14" gorm:"column:col_key14;comment:预留字段"`
	Key15    string  `json:"key15" gorm:"column:col_key15;comment:预留字段"`
	Key16    string  `json:"key16" gorm:"column:col_key16;comment:预留字段"`
	Key17    string  `json:"key17" gorm:"column:col_key17;comment:预留字段"`
	Key18    string  `json:"key18" gorm:"column:col_key18;comment:预留字段"`
	Key19    string  `json:"key19" gorm:"column:col_key19;comment:预留字段"`
	Key20    string  `json:"key20" gorm:"column:col_key20;comment:预留字段"`
	Key21    string  `json:"key21" gorm:"column:col_key21;comment:预留字段"`
	Key22    string  `json:"key22" gorm:"column:col_key22;comment:预留字段"`
	Key23    string  `json:"key23" gorm:"column:col_key23;comment:预留字段"`
	Key24    string  `json:"key24" gorm:"column:col_key24;comment:预留字段"`
	Key25    string  `json:"key25" gorm:"column:col_key25;comment:预留字段"`
	Key26    string  `json:"key26" gorm:"column:col_key26;comment:预留字段"`
	Key27    string  `json:"key27" gorm:"column:col_key27;comment:预留字段"`
	Key28    string  `json:"key28" gorm:"column:col_key28;comment:预留字段"`
	Key29    string  `json:"key29" gorm:"column:col_key29;comment:预留字段"`
	Key30    string  `json:"key30" gorm:"column:col_key30;comment:预留字段"`
}

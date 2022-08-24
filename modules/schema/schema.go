package schema

import (
	"altsub/base"
	"altsub/models"
	"encoding/json"
	"errors"
	"sort"
	"strings"

	"gorm.io/gorm"
)

type SchemaedEventItem struct {
	Description string      `json:"description"`
	Key         string      `json:"key"`
	CName       string      `json:"cname"`
	From        string      `json:"from"`
	FType       string      `json:"ftype"` // 字段类型（string、int、array、map等）
	SType       string      `json:"stype"` // 告警消息展示类型（pic、link、text等）
	Priority    int         `json:"priority"`
	Value       interface{} `json:"value"`
}

type SchemaedEvent []*SchemaedEventItem

type SchemaedEvents []SchemaedEvent

func Fetch(tx *gorm.DB, pq *models.PageQuery) (*models.MSchemas, error) {
	if tx == nil {
		err := errors.New("nil db object")
		base.NewLog("error", err, "获取schema失败", "schema:Fetch()")
		return nil, err
	}
	if pq == nil {
		pq = &models.PageQuery{}
		pq.Page = 1
		pq.Size = 10000
		pq.Order = "+id"
	}
	var ss = &models.MSchemas{
		TX:  tx,
		PQ:  *pq,
		All: []*models.MSchema{},
	}
	return ss, ss.Fetch()
}

func Add(schm *models.MSchema) (err error) {
	if schm == nil {
		err = errors.New("empty schema data")
		base.NewLog("error", err, "新增schema失败", "schema:Add()")
		return
	}
	if schm.TX == nil {
		err = errors.New("nil db object")
		base.NewLog("error", err, "新增schema失败", "schema:Add()")
		return
	}
	return schm.Add()
}

func Update(schm *models.MSchema) (err error) {
	if schm == nil || len(schm.Data) <= 0 {
		err = errors.New("empty schema")
		base.NewLog("error", err, "更新schema失败", "schema:Update()")
		return err
	}
	if err := schm.Get(); err != nil {
		base.NewLog("error", err, "更新schema失败", "schema:Update()")
		return err
	}
	return schm.Update()
}

func GetBySourceName(schm *models.MSchema, srcName string) (err error) {
	if schm.TX == nil {
		err = errors.New("nil db object")
		base.NewLog("error", err, "根据source名称获取schema失败", "schema:GetBySourceName()")
		return
	}
	if len(srcName) <= 0 {
		err = errors.New("empty source name")
		base.NewLog("error", err, "根据source名称获取schema失败", "schema:GetBySourceName()")
		return
	}
	var src = models.MSource{}
	src.TX = schm.TX
	src.Name = srcName
	if err = src.GetByName(); err != nil {
		base.NewLog("error", err, "根据source名称获取schema失败", "schema:GetBySourceName()")
		return
	}
	schm.SourceID = src.ID
	if err = schm.GetBySourceID(); err != nil {
		base.NewLog("error", err, "根据source名称获取schema", "schema:GetBySourceName()")
		return
	}
	return nil
}

func Get(schm *models.MSchema) error {
	if schm.TX == nil {
		err := errors.New("nil db object")
		base.NewLog("error", err, "获取schema失败", "schema:Get()")
		return err
	}
	return schm.Get()
}

func ParseEvent(schm *models.MSchema, ev *models.MEvent) (SchemaedEvents, error) {
	// schema数据自身解析
	var schmItems = &SchemaedEvent{}
	if err := json.Unmarshal(schm.Data, &schmItems); err != nil {
		base.NewLog("error", err, "事件解析失败", "schema:ParseEvent()")
		return nil, err
	}
	// 解析事件
	var rawev interface{}
	var rawevm map[string]interface{}
	var raweva []interface{}
	if err := json.Unmarshal(ev.Data, &rawev); err != nil {
		base.NewLog("error", err, "事件解析失败", "schema:ParseEvent()")
		return nil, err
	}
	switch schm.EvField {
	case ".":
	default: //
		evfs := strings.Split(schm.EvField, ".")
		for _, evf := range evfs {
			rawev = rawev.(map[string]interface{})[evf]
		}
	}
	var parsedEvma = SchemaedEvents{}
	switch schm.EvType {
	case "array":
		raweva, _ = rawev.([]interface{})
		for _, rawev := range raweva {
			ev, _ := rawev.(map[string]interface{})
			schmItems.parseEvent(ev)
			parsedEvma = append(parsedEvma, *schmItems)
		}
	default: // case "map":
		rawevm, _ = rawev.(map[string]interface{})
		ev := rawevm
		schmItems.parseEvent(ev)
		parsedEvma = []SchemaedEvent{*schmItems}
	}
	return parsedEvma, nil
}

func (SchemaedEvent *SchemaedEvent) parseEvent(ev map[string]interface{}) {
	for i, SchemaedEventItem := range *SchemaedEvent {
		var v = func(e map[string]interface{}) map[string]interface{} { return e }(ev)
		fieldsPathSlice := strings.Split(SchemaedEventItem.From, ".")
		for n, field := range fieldsPathSlice {
			if n+1 < len(fieldsPathSlice) {
				v, _ = v[field].(map[string]interface{})
			} else {
				(*SchemaedEvent)[i].Value = v[field]
			}
		}
	}
	// 对处理完成的事件根据字段优先级排序
	sort.Slice(*SchemaedEvent, func(i, j int) bool {
		return (*SchemaedEvent)[i].Priority > (*SchemaedEvent)[j].Priority
	})
}

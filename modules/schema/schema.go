package schema

import (
	"altsub/base"
	"altsub/models"
	"encoding/json"
	"errors"
	"sort"
	"strings"
)

type SchemaItem struct {
	Description string      `json:"description"`
	Name        string      `json:"name"`
	CName       string      `json:"cname"`
	From        string      `json:"from"`
	FType       string      `json:"ftype"` // 字段类型（string、int、array、map等）
	SType       string      `json:"stype"` // 告警消息展示类型（pic、link、text等）
	Priority    int         `json:"priority"`
	Value       interface{} `json:"-"`
}

type SchemaItems []*SchemaItem

func Add(schm *models.MSchema) (err error) {
	if schm == nil {
		err = errors.New("empty schema data")
		base.NewLog("error", err, "新增schema失败", "schema:Add()")
		return
	}
	if schm.BaseModel.TX == nil {
		err = errors.New("nil db object")
		base.NewLog("error", err, "新增schema失败", "schema:Add()")
		return
	}
	return schm.Add()
}

func GetBySourceName(schm *models.MSchema, srcName string) (err error) {
	if schm.BaseModel.TX == nil {
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
	src.BaseModel.TX = schm.BaseModel.TX
	src.Name = srcName
	if err = src.GetByName(); err != nil {
		base.NewLog("error", err, "根据source名称获取schema失败", "schema:GetBySourceName()")
		return
	}
	schm.SourceID = src.ID
	err = schm.GetBySourceID()
	base.NewLog("", err, "根据source名称获取schema", "schema:GetBySourceName()")
	return
}

func ParseEvent(schm *models.MSchema, ev *models.MEvent) error {
	// schema数据自身解析
	var schmItems = &SchemaItems{}
	if err := json.Unmarshal(schm.Data, &schmItems); err != nil {
		base.NewLog("error", err, "事件解析失败", "schema:ParseEvent()")
		return err
	}
	// 解析事件
	var rawev interface{}
	var rawevm map[string]interface{}
	var raweva []interface{}
	if err := json.Unmarshal(ev.Data, &rawev); err != nil {
		base.NewLog("error", err, "事件解析失败", "schema:ParseEvent()")
		return err
	}
	switch schm.EvField {
	case ".":
	default: //
		evfs := strings.Split(schm.EvField, ".")
		for _, evf := range evfs {
			rawev = rawev.(map[string]interface{})[evf]
		}
	}
	var parsedEvma = []SchemaItems{}
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
		parsedEvma = []SchemaItems{*schmItems}
	}
	// 事件处理-维护检测
	// 事件处理-抑制检测
	// 事件处理-认领检测
	// 事件处理-指派检测
	// 事件处理-告警发送

	return nil
}

func (schemaItems *SchemaItems) parseEvent(ev map[string]interface{}) {
	for i, schemaItem := range *schemaItems {
		var v = func(e map[string]interface{}) map[string]interface{} { return e }(ev)
		fieldsPathSlice := strings.Split(schemaItem.From, ".")
		for n, field := range fieldsPathSlice {
			if n+1 < len(fieldsPathSlice) {
				v, _ = v[field].(map[string]interface{})
			} else {
				(*schemaItems)[i].Value = v[field]
			}
		}
	}
	// 对处理完成的事件根据字段优先级排序
	sort.Slice(*schemaItems, func(i, j int) bool {
		return (*schemaItems)[i].Priority > (*schemaItems)[j].Priority
	})
}

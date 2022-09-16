package event

import (
	"altsub/base"
	"altsub/models"
	"altsub/modules/schema"
	"crypto/md5"
	"encoding/json"
	"fmt"
)

// 接收事件，写入kafka
func Receive(srcName string, ev models.MEvent) error {
	// 根据source名称生成topic，将事件写入指定topic
	return base.WriteToKafka(srcName, ev.Data)
}

// 从kafka消费
func ReadAndParseEventFromBufferForever(srcs ...string) {
	for _, src := range srcs {
		go func(src string) {
			// 从对应message buffer消费事件
			buf := base.ReadFromKafka(src)
			for {
				rawEv := <-buf
				base.NewLog("trace", nil, fmt.Sprintf("从kafka读取到数据 rawEv: %s", string(rawEv)), "event:ReadAndParseEventFromBufferForever()")
				// 根据source名称找对应schema解析事件
				var schm = models.MSchema{}
				schm.TX = base.DB()
				if err := schema.GetBySourceName(&schm, src); err != nil {
					base.NewLog("error", err, "根据 source name 获取 schema 失败", "ReadAndParseEventFromBufferForever()")
					continue
				}
				var ev = models.MEvent{Data: rawEv}
				if parsedEvs, err := schema.ParseEvent(&schm, &ev); err != nil {
					base.NewLog("error", err, "事件解析失败", "ReadAndParseEventFromBufferForever()")
					continue
				} else {
					// 原始事件入库
					if err := StoreRawToDb(&ev); err != nil {
						base.NewLog("warn", err, "事件入库失败", "ReadAndParseEventFromBufferForever()")
					}
					base.NewLog("trace", nil, fmt.Sprintf("事件解析成功：%#v", parsedEvs), "ReadAndParseEventFromBufferForever()")
					for _, parseEv := range parsedEvs {
						go func(parseEv schema.SchemaedEvent) {
							// 事件处理-匹配规则
							var rs []models.MRule = checkRules(src, parseEv)
							// 解析后数据入库
							go func(parseEv schema.SchemaedEvent, rs []models.MRule) {
								if err := StoreParsedEvToDb(parseEv, rs); err != nil {
									base.NewLog("warn", err, "解析后事件存储失败", "ReadAndParseEventFromBufferForever()")
									return
								}
							}(parseEv, rs)
							// 事件处理-发版检测
							var bigtype = ""
							var instance = ""
							var service = ""
							for _, item := range parseEv {
								if item.Key == "bigtype" {
									bigtype, _ = item.Value.(string)
								} else if item.Key == "instance" {
									instance, _ = item.Value.(string)
								} else if item.Key == "service" {
									service, _ = item.Value.(string)
								}
							}
							if checkPublish(bigtype, service, instance) {
								base.NewLog("info", nil, "事件检测 - 检测到服务发版", "ReadAndParseEventFromBufferForever()")
								return
							}
							// 事件处理-维护检测
							if checkMaintenance(rs) {
								base.NewLog("info", nil, "事件检测 - 检测到相关维护项", "ReadAndParseEventFromBufferForever()")
								//
								// TODO：检测到维护事项，是否需要记录入库？
								//
								return
							}
							// 事件处理-抑制检测

							// 事件处理-认领检测

							// 事件处理-指派(订阅)检测
							var rcvs []models.MReceiver
							rcvs = append(rcvs, checkSubscribe(rs)...)

							base.NewLog("debug", nil, "发送至全局接收者", "ReadAndParseEventFromBufferForever()")
							// send to global group
							auth := map[string]interface{}{
								"app_key":    base.Config.NotificationConf.DefaultDingTalkAppKey,
								"app_secret": base.Config.NotificationConf.DefaultDingTalkAppSecret,
								"chat_id":    base.Config.NotificationConf.DefaultDingtalkChatID,
							}
							ab, _ := json.Marshal(auth)
							hs := md5.New()
							hs.Write(ab)
							rcvs = append(rcvs, models.MReceiver{
								Type:     models.ReceiverTypeDingtalkApp,
								Auth:     ab,
								AuthHash: fmt.Sprintf("%x", hs.Sum(nil)),
							})

							base.NewLog("debug", nil, "获取服务默认接收群组", "ReadAndParseEventFromBufferForever()")
							// send to default group
							rcvs = append(rcvs, getDefaultServiceGroup(bigtype, service)...)

							if len(rcvs) > 0 {
								//接收者去重处理
								var rcvm = map[string]models.MReceiver{}
								for _, v := range rcvs {
									if _, ok := rcvm[v.AuthHash]; !ok {
										rcvm[v.AuthHash] = v
									}
								}
								rcvs = rcvs[:0]
								for _, rcv := range rcvm {
									rcvs = append(rcvs, rcv)
								}
								// 事件处理-告警发送
								Notice(parseEv, rcvs)
							}
							// 事件发送记录（事件未发送（维护、抑制）、事件已发送但失败、事件已发送且成功）
						}(parseEv)
					}
				}
			}
		}(src)
	}
}

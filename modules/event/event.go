package event

import (
	"altsub/base"
	"altsub/models"
	"altsub/modules/schema"
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
				// 根据source名称找对应schema解析事件
				var schm = models.MSchema{}
				schm.TX = base.DB()
				if err := schema.GetBySourceName(&schm, src); err != nil {
					base.NewLog("error", err, "从buffer读取并解析事件失败", "ReadAndParseEventFromBufferForever()")
					continue
				}
				var ev = models.MEvent{}
				ev.TX = schm.TX
				ev.Data = rawEv
				if parsedEvs, err := schema.ParseEvent(&schm, &ev); err != nil {
					base.NewLog("error", err, "事件解析失败", "ReadAndParseEventFromBufferForever()")
					continue
				} else {
					base.NewLog("trace", nil, fmt.Sprintf("事件解析：%#v", parsedEvs), "ReadAndParseEventFromBufferForever()")
					// 事件入库
					if err := StoreToDb(&ev); err != nil {
						base.NewLog("warn", err, "事件入库失败", "ReadAndParseEventFromBufferForever()")
					}
					// 事件处理-维护检测
					// 事件处理-抑制检测
					// 事件处理-认领检测
					// 事件处理-指派(订阅)检测
					// 事件处理-告警发送
					// 事件发送记录（事件未发送（维护、抑制）、事件已发送但失败、事件已发送且成功）

				}
			}
		}(src)
	}
}

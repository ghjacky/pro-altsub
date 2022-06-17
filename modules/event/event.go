package event

import (
	"altsub/base"
	"altsub/models"
	"altsub/modules/schema"
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
				if err := schema.ParseEvent(&schm, &ev); err != nil {
					base.NewLog("error", err, "事件解析失败", "ReadAndParseEventFromBufferForever()")
					continue
				}
				// 事件入库

				// 维护检测

				// 告警抑制检测

				// 订阅（指派）监测

				// 发送告警事件

				// 事件发送记录（事件未发送（维护、抑制）、事件已发送但失败、事件已发送且成功）

			}
		}(src)
	}
}

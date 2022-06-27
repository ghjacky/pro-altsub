package notification

import (
	"altsub/models"
	"altsub/modules/schema"
	"crypto/tls"
	"net/http"
)

type INotification interface {
	Notice(string) error
	RenderMsg() string
	ParseAuth([]byte)
	SetEvent(schema.SchemaedEvent)
}

func NewNotification(tp int) INotification {
	hc := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true,
			},
		},
	}
	switch tp {
	case models.ReceiverTypeDingtalkApp:
		return &DingTalk{
			Client: hc,
		}
	case models.ReceiverTypeDingtalkPersonal:
		return &DingPersonal{
			Client: hc,
		}
	case models.ReceiverTypeSMS:
		return &SMS{
			Client: hc,
		}
	case models.ReceiverTypeVoice:
		return &Voice{
			Client: hc,
		}
	default:
		return nil
	}
}

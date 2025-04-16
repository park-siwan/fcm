package sendMessage

import (
	"fmt"

	"github.com/appleboy/go-fcm"
	"goTest/ihplogger"
)

// FCMMessage FCM 메시지 구조체
type FCMMessage struct {
	Token string
	Title string
	Body  string
	Date  string
	Link  string
}

// SendMessageWithFCM appleboy/go-fcm 패키지를 사용한 FCM 메시지 전송
func (info *FCMMessage) SendMessageWithFCM() (success int, err error) {
	// FCM 서버 키
	serverKey := "AAAAVgKeCHE:APA91bGdghN4yRB8oSRsQ_4ME15QvmQdx4L4DeODIKkTZwEwXLxyzUHcHkInwlGwfbqJz7Vr_r5fmBPzAaEAO2WQakGuv0B6hsZmC4GEBrmP3vU9zhDKq3P-cm3CULsxXo2Kb20RAslt"

	// FCM 클라이언트 생성
	client, err := fcm.NewClient(serverKey)
	if err != nil {
		ihplogger.Log.Errorf("FCM 클라이언트 생성 실패: %v", err)
		return 0, fmt.Errorf("FCM 클라이언트 생성 실패: %v", err)
	}

	// 재시도 로직
	successCount := 0
	for i := 0; i < 10; i++ {
		// 메시지 생성
		msg := &fcm.Message{
			To: info.Token,
			Data: map[string]interface{}{
				"title": info.Title,
				"body":  info.Body,
				"date":  info.Date,
				"link":  info.Link,
			},
			Notification: &fcm.Notification{
				Title: info.Title,
				Body:  info.Body,
				Sound: "default",
			},
		}

		// iOS 설정 추가
		msg.ApnsConfig = map[string]interface{}{
			"headers": map[string]string{
				"apns-priority":  "10",
				"apns-push-type": "alert",
			},
			"payload": map[string]interface{}{
				"aps": map[string]interface{}{
					"alert": map[string]interface{}{
						"title": info.Title,
						"body":  info.Body,
					},
					"sound":             "default",
					"content-available": 1,
					"mutable-content":   1,
					"badge":             1,
				},
			},
		}

		// 메시지 전송
		response, err := client.Send(msg)
		if err != nil {
			ihplogger.Log.Errorf("메시지 전송 오류 (시도 %d/10): %v", i+1, err)
			continue
		}

		// 성공 여부 확인
		if response.Success == 1 {
			ihplogger.Log.Infof("메시지 전송 성공: %v", response)
			successCount = 1
			return successCount, nil
		}

		// 오류 상세 정보 로깅
		ihplogger.Log.Warnf("FCM 응답 (시도 %d/10): 성공=%d, 실패=%d", 
			i+1, response.Success, response.Failure)
		
		if len(response.Results) > 0 {
			for idx, result := range response.Results {
				if result.Error != nil {
					ihplogger.Log.Warnf("결과 %d: 오류=%v", idx, result.Error)
				}
			}
		}

		// 마지막 시도에서 실패한 경우
		if i >= 9 {
			errMsg := "10회 시도 후 메시지 전송 실패"
			if response.Failure > 0 && len(response.Results) > 0 && response.Results[0].Error != nil {
				errMsg = fmt.Sprintf("메시지 전송 실패: %v", response.Results[0].Error)
			}
			ihplogger.Log.Error(errMsg)
			return 0, fmt.Errorf(errMsg)
		}
	}

	return successCount, nil
} 
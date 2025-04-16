package sendMessage

import (
	"context"
	"fmt"

	firebase "firebase.google.com/go/v4"
	"firebase.google.com/go/v4/messaging"
	"goTest/ihplogger"
	"google.golang.org/api/option"
)

type Message struct {
	Token string
	Title string
	Body  string
	Date  string
	Link  string
}

func (info *Message) SendMessage() (success int, err error) {
	// Firebase Admin SDK 초기화
	ctx := context.Background()
	
	// Firebase 프로젝트 설정
	config := &firebase.Config{
		ProjectID: "inhandplus-d8fb1",
	}
	
	opt := option.WithCredentialsFile("../certs/inhandplus-d8fb1-firebase-adminsdk-9db5s-2e80d39c16.json")
	app, err := firebase.NewApp(ctx, config, opt)
	if err != nil {
		ihplogger.Log.Errorf("Firebase 앱 초기화 실패: %v", err)
		return 0, fmt.Errorf("Firebase 앱 초기화 실패: %v", err)
	}

	// FCM 클라이언트 생성
	client, err := app.Messaging(ctx)
	if err != nil {
		ihplogger.Log.Errorf("FCM 클라이언트 생성 실패: %v", err)
		return 0, fmt.Errorf("FCM 클라이언트 생성 실패: %v", err)
	}

	// APNS 설정
	apnsConfig := messaging.APNSConfig{
		Payload: &messaging.APNSPayload{
			Aps: &messaging.Aps{
				Alert: &messaging.ApsAlert{
					Title: info.Title,
					Body:  info.Body,
				},
				ContentAvailable: true,
			},
		},
		Headers: map[string]string{
			"apns-priority":  "10",
			"apns-push-type": "alert",
		},
	}

	// 재시도 로직
	successCount := 0
	for i := 0; i < 10; i++ {
		// 메시지 구성
		message := &messaging.Message{
			Token: info.Token,
			Data: map[string]string{
				"title": info.Title,
				"body":  info.Body,
				"date":  info.Date,
				"link":  info.Link,
			},
			APNS: &apnsConfig,
			Notification: &messaging.Notification{
				Title: info.Title,
				Body:  info.Body,
			},
		}

		// 메시지 전송
		response, err := client.Send(ctx, message)
		
		if err == nil && response != "" {
			ihplogger.Log.Infof("메시지 전송 성공: %s", response)
			successCount = 1
			return successCount, nil
		}

		if i >= 9 {
			if err != nil {
				ihplogger.Log.Errorf("메시지 전송 실패: %v", err)
				return 0, fmt.Errorf("메시지 전송 실패: %v", err)
			}
			ihplogger.Log.Error("10회 시도했으나 메시지 전송에 실패했습니다")
			return 0, fmt.Errorf("10회 시도했으나 메시지 전송에 실패했습니다")
		}
		ihplogger.Log.Warnf("메시지 전송 시도 %d/10 실패, 재시도 중...", i+1)
	}

	return successCount, nil
} 
package main

import (
	"context"
	"fmt"
	"log"

	firebase "firebase.google.com/go/v4"
	"firebase.google.com/go/v4/messaging"
	"google.golang.org/api/option"
)

func main() {
	// 디바이스 토큰
	token := "eZB1KmXSrkr8nxcjKvA7BF:APA91bFxaNoHZuOj3yLxaN0wyuXmDQb_skZl8APa5SAC9KgWb_a_PiZCOMVsWrvBlc4zr0sMN8Go745UZp5YQIS-Jjf12Fr5CPXde_-s-ZZhCcZva8JKI0A"

	// 컨텍스트 생성
	ctx := context.Background()
	
	// Firebase 설정
	config := &firebase.Config{
		ProjectID: "inhandplus-d8fb1",
	}
	
	// Firebase 앱 생성
	opt := option.WithCredentialsFile("./certs/serviceAccountKey.json")
	app, err := firebase.NewApp(ctx, config, opt)
	if err != nil {
		log.Fatalf("Firebase 앱 초기화 실패: %v", err)
	}
	
	// FCM 클라이언트 생성
	client, err := app.Messaging(ctx)
	if err != nil {
		log.Fatalf("FCM 클라이언트 생성 실패: %v", err)
	}

	// Badge 값 설정
	badgeValue := 1

	// 메시지 생성
	msg := &messaging.Message{
		Token: token,
		Data: map[string]string{
			"title": "테스트 알림",
			"body":  "이것은 테스트 알림입니다.",
			"date":  "2023-04-01",
			"link":  "https://example.com",
		},
		Notification: &messaging.Notification{
			Title: "테스트 알림",
			Body:  "이것은 테스트 알림입니다.",
		},
		// APNS 설정 추가
		APNS: &messaging.APNSConfig{
			Headers: map[string]string{
				"apns-priority":  "10",
				"apns-push-type": "alert",
			},
			Payload: &messaging.APNSPayload{
				Aps: &messaging.Aps{
					Alert: &messaging.ApsAlert{
						Title: "테스트 알림",
						Body:  "이것은 테스트 알림입니다.",
					},
					Sound:            "default",
					ContentAvailable: true,
					MutableContent:   true,
					Badge:            &badgeValue,
				},
			},
		},
	}

	// 메시지 전송
	response, err := client.Send(ctx, msg)
	if err != nil {
		log.Fatalf("메시지 전송 실패: %v", err)
	}

	fmt.Printf("메시지 전송 성공: %s\n", response)
}
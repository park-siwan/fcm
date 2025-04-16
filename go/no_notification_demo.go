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

	// Data-only 메시지 생성 (iOS에서는 시각적으로 표시되지 않음)
	// 백그라운드에서 처리됨
	msg := &messaging.Message{
		Token: token,
		Data: map[string]string{
			"title": "iOS에 표시되지 않는 메시지",
			"body":  "이 메시지는 Data-only 메시지로, iOS에서 시각적으로 표시되지 않습니다",
			"date":  "2023-04-01",
			"link":  "https://example.com",
		},
		// Notification 필드가 없으면 iOS에서 푸시 알림이 표시되지 않음
		// APNS 설정만으로는 충분하지 않음
		APNS: &messaging.APNSConfig{
			Headers: map[string]string{
				"apns-priority":  "5", // 낮은 우선순위 (백그라운드 메시지)
				"apns-push-type": "background",
			},
			Payload: &messaging.APNSPayload{
				Aps: &messaging.Aps{
					ContentAvailable: true,
					// 알림 필드가 없음 (Alert 없음)
				},
			},
		},
	}

	// 메시지 전송
	log.Println("iOS에 표시되지 않는 Data-only 메시지 전송 시작...")
	response, err := client.Send(ctx, msg)
	if err != nil {
		log.Fatalf("메시지 전송 실패: %v", err)
	}

	fmt.Printf("메시지 전송 성공: %s\n", response)
	log.Println("중요: 메시지가 전송되었지만 iOS 디바이스에서는 시각적으로 표시되지 않을 것입니다 (Data-only 메시지)")
} 
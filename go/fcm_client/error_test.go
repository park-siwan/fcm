package main

import (
	"context"
	"fmt"
	"log"

	"firebase.google.com/go/v4/messaging"
	fcm "github.com/appleboy/go-fcm"
)

// Message 메시지 구조체
type Message struct {
	Token string
	Title string
	Body  string
	Date  string
	Link  string
}

// SendMessage appleboy/go-fcm과 Firebase Admin SDK의 혼합 사용으로 인한 오류 재현 코드
func (info *Message) SendMessage() (resp *messaging.BatchResponse, err error) {
	ctx := context.Background()
	client, err := fcm.NewClient(
		ctx,
		fcm.WithCredentialsFile("./certs/serviceAccountKey.json"),
	)

	if err != nil {
		log.Printf("오류: %v", err)
		return resp, err
	}

	APNS := messaging.APNSConfig{
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

	for i := 0; i < 10; i++ {
		resp, err = client.Send(
			ctx,
			&messaging.Message{
				Token: info.Token,
				Data: map[string]string{
					"title": info.Title,
					"body":  info.Body,
					"date":  info.Date,
					"link":  info.Link,
				},
				APNS: &APNS,
				/* Notification: &messaging.Notification{
					Title: info.Title,
					Body:  info.Body,
				}, */
			},
		)

		if err != nil {
			log.Printf("메시지 전송 실패 (시도 %d): %v", i, err)
		}

		if resp != nil && resp.SuccessCount > 0 {
			break
		}

		if i >= 9 {
			if resp != nil && len(resp.Responses) > 0 {
				log.Printf("10회 시도 후 메시지 전송 실패: %v", resp.Responses[0].Error)
				return resp, fmt.Errorf("메시지 전송 실패: %v", resp.Responses[0].Error)
			}
			return resp, fmt.Errorf("10회 시도 후 메시지 전송 실패")
		}
	}

	return resp, nil
}

func main() {
	// 메시지 구성
	msg := &Message{
		Token: "eZB1KmXSrkr8nxcjKvA7BF:APA91bFxaNoHZuOj3yLxaN0wyuXmDQb_skZl8APa5SAC9KgWb_a_PiZCOMVsWrvBlc4zr0sMN8Go745UZp5YQIS-Jjf12Fr5CPXde_-s-ZZhCcZva8JKI0A",
		Title: "오류 재현 테스트",
		Body:  "이 메시지는 오류 재현을 위한 테스트입니다",
		Date:  "2023-04-01",
		Link:  "https://example.com",
	}

	// 오류 재현 테스트
	log.Println("FCM 메시지 전송 오류 재현 시작...")
	resp, err := msg.SendMessage()
	if err != nil {
		log.Printf("예상된 오류 발생: %v", err)
		return
	}

	log.Printf("메시지 전송 결과: %+v", resp)
} 
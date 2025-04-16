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
	
	// appleboy/go-fcm 라이브러리에서는 WithCredentialsFile() 메서드가 다른 의미로 사용됨
	// Firebase Admin SDK와 혼용하면 타입 불일치 오류 발생
	client, err := fcm.NewClient(
		ctx,
		fcm.WithCredentialsFile("./certs/serviceAccountKey.json"),
	)

	if err != nil {
		log.Printf("오류: %v", err)
		return resp, err
	}

	// 불호환성 1: Firebase Admin SDK의 APNSConfig와 appleboy/go-fcm의 혼합 사용
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

	// 불호환성 2: client.Send() 메서드는 firebase.google.com/go/v4/messaging.Message를 
	// 기대하지 않고 appleboy/go-fcm의 Message 타입을 기대함
	for i := 0; i < 3; i++ {
		// 이 부분에서 타입 불일치 오류 발생
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
			},
		)

		if err != nil {
			log.Printf("메시지 전송 실패 (시도 %d): %v", i, err)
			// 오류 발생하여 더 이상 시도하지 않음
			break
		}
	}

	// 이 부분은 실행되지 않음 (에러가 발생하기 때문)
	return resp, fmt.Errorf("라이브러리 간 타입 불일치로 인한 기대된 오류")
}

func main() {
	// 메시지 구성
	msg := &Message{
		Token: "eZB1KmXSrkr8nxcjKvA7BF:APA91bFxaNoHZuOj3yLxaN0wyuXmDQb_skZl8APa5SAC9KgWb_a_PiZCOMVsWrvBlc4zr0sMN8Go745UZp5YQIS-Jjf12Fr5CPXde_-s-ZZhCcZva8JKI0A",
		Title: "호환성 오류 테스트",
		Body:  "이 메시지는 appleboy/go-fcm과 firebase.google.com/go/v4/messaging 간의 호환성 문제 테스트입니다",
		Date:  "2023-04-01",
		Link:  "https://example.com",
	}

	// 오류 재현 테스트
	log.Println("FCM 메시지 전송 라이브러리 호환성 오류 재현 시작...")
	resp, err := msg.SendMessage()
	if err != nil {
		log.Printf("예상된 오류 발생: %v", err)
	} else {
		log.Printf("예상과 다르게 메시지 전송 성공(?): %+v", resp)
	}
} 
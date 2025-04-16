package main

import (
	"fmt"
	"os"

	"goTest/ihplogger"
	"goTest/sendMessage"
)

func main() {
	// 로그 초기화
	ihplogger.LogInit()
	ihplogger.Log.Info("FCM 메시지 전송 테스트 시작 (appleboy/go-fcm)")

	// FCM 메시지 구성
	message := sendMessage.FCMMessage{
		Token: "eZB1KmXSrkr8nxcjKvA7BF:APA91bFxaNoHZuOj3yLxaN0wyuXmDQb_skZl8APa5SAC9KgWb_a_PiZCOMVsWrvBlc4zr0sMN8Go745UZp5YQIS-Jjf12Fr5CPXde_-s-ZZhCcZva8JKI0A", 
		Title: "FCM 테스트 알림",
		Body:  "appleboy/go-fcm 패키지를 이용한 테스트 알림입니다.",
		Date:  "2023-04-01",
		Link:  "https://example.com",
	}

	// FCM 메시지 전송
	ihplogger.Log.Infof("메시지 전송 시작: %s", message.Title)
	success, err := message.SendMessageWithFCM()
	if err != nil {
		ihplogger.Log.Errorf("메시지 전송 오류: %v", err)
		fmt.Fprintf(os.Stderr, "오류: %v\n", err)
		os.Exit(1)
	}

	// 결과 출력
	ihplogger.Log.Infof("성공적으로 %d개의 메시지를 전송했습니다", success)
	fmt.Printf("성공적으로 %d개의 메시지를 전송했습니다.\n", success)
} 
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
	ihplogger.Log.Info("FCM 메시지 전송 프로그램 시작")

	// Firebase Admin SDK를 사용하는 새 구현에서는 서버 키가 필요하지 않습니다.
	// 대신 JSON 인증 파일을 사용합니다.

	// 실제 사용 시에는 유효한 디바이스 토큰을 여기에 입력해야 합니다.
	message := sendMessage.Message{
		Token: "eZB1KmXSrkr8nxcjKvA7BF:APA91bFxaNoHZuOj3yLxaN0wyuXmDQb_skZl8APa5SAC9KgWb_a_PiZCOMVsWrvBlc4zr0sMN8Go745UZp5YQIS-Jjf12Fr5CPXde_-s-ZZhCcZva8JKI0A", // 실제 디바이스 토큰으로 교체하세요
		Title: "테스트 알림",
		Body:  "알림 내용입니다",
		Date:  "2023-04-01",
		Link:  "https://example.com",
	}

	ihplogger.Log.Infof("메시지 전송 시작: %s", message.Title)
	success, err := message.SendMessage()
	if err != nil {
		ihplogger.Log.Errorf("메시지 전송 오류: %v", err)
		fmt.Fprintf(os.Stderr, "오류: %v\n", err)
		os.Exit(1)
	}

	ihplogger.Log.Infof("성공적으로 %d개의 메시지를 전송했습니다", success)
	fmt.Printf("성공적으로 %d개의 메시지를 전송했습니다.\n", success)
}
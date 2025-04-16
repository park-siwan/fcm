# FCM 푸시 알림 전송 예제

이 프로젝트는 Firebase Cloud Messaging(FCM)을 사용하여 iOS 및 Android 디바이스로 푸시 알림을 전송하는 Go 언어 예제입니다.

## 필요 조건

- Go 1.19 이상
- Firebase 프로젝트
- FCM 서버 키
- 유효한 디바이스 토큰

## 설치 방법

```bash
# 프로젝트 복제
git clone https://github.com/your-username/fcm-example.git
cd fcm-example

# 의존성 설치
go mod tidy
```

## 사용 방법

1. `cmd/main.go` 파일의 `serverKey` 변수에 Firebase 콘솔에서 얻은 서버 키를 입력하세요.
2. `message.Token` 값을 실제 디바이스에서 생성된 FCM 토큰으로 변경하세요.
3. 필요에 따라 메시지 제목, 내용, 날짜, 링크를 수정하세요.
4. 다음 명령어로 예제를 실행하세요:

```bash
go run cmd/main.go
```

## 코드 구조

- `sendMessage/fcm.go`: FCM 메시지 전송 로직이 구현된 패키지
- `cmd/main.go`: 메인 애플리케이션 엔트리 포인트

## 주요 기능

- 단일 디바이스로 푸시 알림 전송
- 자동 재시도 로직 구현 (최대 10회)
- 제목, 내용, 날짜, 링크 등의 데이터 전송

## 사용 예제

```go
package main

import (
    "fmt"
    "os"
    "goTest/sendMessage"
)

func main() {
    // Firebase 콘솔에서 얻은 서버 키
    serverKey := "YOUR_FCM_SERVER_KEY"

    message := sendMessage.Message{
        Token: "DEVICE_TOKEN",
        Title: "테스트 알림",
        Body:  "알림 내용입니다",
        Date:  "2023-04-01",
        Link:  "https://example.com",
    }

    success, err := message.SendMessage(serverKey)
    if err != nil {
        fmt.Fprintf(os.Stderr, "오류: %v\n", err)
        os.Exit(1)
    }

    fmt.Printf("성공적으로 %d개의 메시지를 전송했습니다.\n", success)
}
```

## 라이센스

MIT

## 문제 해결

- "404 Not Found" 오류: 유효하지 않은 디바이스 토큰을 사용하고 있을 수 있습니다. 디바이스 토큰을 확인하세요.
- "401 Unauthorized" 오류: 서버 키가 올바르지 않을 수 있습니다. Firebase 콘솔에서 서버 키를 확인하세요.

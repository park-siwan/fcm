const axios = require('axios');
const getAccessToken = require('./getAccessToken');
require('dotenv').config();

// FCM HTTP v1 API URL
const FCM_URL =
  'https://fcm.googleapis.com/v1/projects/inhandplus-d8fb1/messages:send';

// 토큰 배열 (플랫폼별로 구분)
const tokens = {
  ios: [
    'cdq_4SieIUrXkfb_h3R1AY:APA91bF_296uESJdOOr1fytxt_kX-4VnRa8EBGcEoAyYagv8eTr2hnxrp0AW8khRlfhQG0ID9j9-VHPA50QiV6Vlp-_bawm2TCqSwu4t1t6QXMDDV1luVdA',
  ],
  android: [
    'cCXn-pRQRnafokLmnG9SXU:APA91bH7ve4zu4zjhgo9cbRDu3uxwDMqkEdU_dy5ey9iMn-hkLDP66ayuOq4K-UFMekgSw5W2E-LjLB4yvWNGm31yToZb_RFmPS55iKK0N-2mlGhEmcyda8',
  ],
};

// 메시지 데이터 (공통 부분)
const notificationTitle = '응급상황(낙상감지)';
const notificationBody = '참여자의 낙상이 감지되었습니다.';
const targetUrl =
  'http://localhost:3000/ko/help/9154d7dfe98b46af9a39bbf373039fec';

// 플랫폼별 메시지 생성 함수
function createMessage(token, platform) {
  console.log(`\n[메시지 생성] 플랫폼: ${platform}`);

  // 안드로이드 전용 메시지
  if (platform === 'android') {
    console.log('[Android] 메시지 생성 시작');
    const androidMessage = {
      message: {
        token: token,
        data: {
          title: notificationTitle,
          body: notificationBody,
          target_url: targetUrl,
        },
      },
    };
    console.log('[Android] 메시지 생성 완료');
    return androidMessage;
  }
  // iOS 전용 메시지
  else if (platform === 'ios') {
    console.log('[iOS] 메시지 생성 시작');
    const iosMessage = {
      message: {
        token: token,

        data: {
          title: notificationTitle,
          body: notificationBody,
          target_url: targetUrl,
        },
        apns: {
          headers: {
            'apns-priority': '10',
            'apns-push-type': 'alert',
          },
          payload: {
            aps: {
              alert: {
                title: notificationTitle,
                body: notificationBody,
              },
              'mutable-content': 1,
              sound: 'default',
            },
          },
        },
      },
    };
    console.log('[iOS] 메시지 생성 완료');
    return iosMessage;
  }

  console.log('[ERROR] 알 수 없는 플랫폼:', platform);
  return null;
}

async function sendMessage() {
  try {
    console.log('\n=== FCM 메시지 전송 시작 ===');

    // Access Token 가져오기
    const accessToken = await getAccessToken();
    console.log('Access Token 획득 완료');

    // 전송 결과 저장
    const results = [];
    const allTokens = [];

    // 토큰 정보 출력
    console.log('\n[토큰 현황]');
    console.log('iOS 토큰 수:', tokens.ios.length);
    console.log('Android 토큰 수:', tokens.android.length);

    // iOS 토큰 처리
    for (const token of tokens.ios) {
      console.log('\n[iOS] 토큰 추가:', token.substring(0, 10) + '...');
      allTokens.push({ token, platform: 'ios' });
    }

    // 안드로이드 토큰 처리
    for (const token of tokens.android) {
      console.log('\n[Android] 토큰 추가:', token.substring(0, 10) + '...');
      allTokens.push({ token, platform: 'android' });
    }

    console.log('\n[전송 준비]');
    console.log('총 처리할 토큰 수:', allTokens.length);
    console.log(
      '플랫폼 별 토큰:',
      allTokens.map((t) => t.platform),
    );

    // 모든 토큰에 대해 메시지 전송
    for (const { token, platform } of allTokens) {
      console.log(
        `\n[전송 시작] ${platform.toUpperCase()} 토큰: ${token.substring(
          0,
          10,
        )}...`,
      );

      const message = createMessage(token, platform);
      if (!message) {
        console.error(`[ERROR] ${platform} 메시지 생성 실패`);
        continue;
      }

      // FCM API 호출
      try {
        const response = await axios.post(FCM_URL, message, {
          headers: {
            Authorization: `Bearer ${accessToken}`,
            'Content-Type': 'application/json',
          },
        });

        console.log(`[성공] ${platform} 메시지 전송 완료:`, response.data);
        results.push({ token, platform, success: true, data: response.data });
      } catch (err) {
        console.error(
          `[실패] ${platform} 메시지 전송 오류:`,
          err.response ? err.response.data : err.message,
        );

        if (err.response?.data?.error?.details) {
          console.error('[상세 오류]:', err.response.data.error.details);
        }

        results.push({
          token,
          platform,
          success: false,
          error: err.response ? err.response.data : err.message,
        });
      }
    }

    // 결과 요약
    console.log('\n=== 전송 결과 요약 ===');
    const summary = {
      total: allTokens.length,
      success: results.filter((r) => r.success).length,
      failure: results.filter((r) => !r.success).length,
      byPlatform: {
        ios: {
          total: results.filter((r) => r.platform === 'ios').length,
          success: results.filter((r) => r.platform === 'ios' && r.success)
            .length,
        },
        android: {
          total: results.filter((r) => r.platform === 'android').length,
          success: results.filter((r) => r.platform === 'android' && r.success)
            .length,
        },
      },
    };
    console.log(summary);
  } catch (error) {
    console.error('\n[치명적 오류]', error.message);
  }
}

// 실행
console.log('\n=== FCM 메시지 전송 프로그램 시작 ===');
sendMessage();

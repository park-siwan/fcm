const axios = require('axios');
const getAccessToken = require('./getAccessToken');
require('dotenv').config();

// FCM HTTP v1 API URL
const FCM_URL =
  'https://fcm.googleapis.com/v1/projects/inhandplus-d8fb1/messages:send';

// 토큰 배열 (플랫폼별로 구분)
const tokens = {
  ios: [
    'eyIF99Bh1k5jlhNANJ3EB-:APA91bHxsGa5-GfXY7jFWwtKUNeL7LBdcejd8BDJxIJDm7RQaDyr_ydyO-48f7R0oRzd4APWpGf8S7tStVKHt52rVo9MGKFkxtUWdDbd3ReeptFWybPn7Io',
  ],
  android: [],
};

// 메시지 데이터 (공통 부분)
const notificationTitle =
  '박시완응급상황(낙상감지)응급상황(낙상감지)응급상황(낙상감지)응급상황(낙상감지)응급상황(낙상감지)응급상황(낙상감지)응급상황(낙상감지)응급상황(낙상감지)';
const notificationBody =
  '박시완참여자의 낙상이 감지되었습니다.참여자의 낙상이 감지되었습니다.참여자의 낙상이 감지되었습니다.참여자의 낙상이 감지되었습니다.참여자의 낙상이 감지되었습니다.참여자의 낙상이 감지되었습니다.참여자의 낙상이 감지되었습니다.참여자의 낙상이 감지되었습니다.';
const targetUrl =
  'http://localhost:3000/ko/help/9154d7dfe98b46af9a39bbf373039fec';

// 플랫폼별 메시지 생성 함수
function createMessage(token, platform) {
  // 기본 데이터 (모든 플랫폼 공통)
  const dataPayload = {
    title: notificationTitle,
    body: notificationBody,
    target_url: targetUrl,
    click_action: 'FLUTTER_NOTIFICATION_CLICK', // Flutter 앱을 위한 액션
    channel_id: 'high_importance_channel', // 안드로이드용 채널 ID
  };

  // 안드로이드 전용 메시지
  if (platform === 'android') {
    return {
      message: {
        token: token,
        data: dataPayload,
        android: {
          priority: 'high',
          notification: {
            sound: 'default',
            default_sound: true,
            default_vibrate_timings: true,
            default_light_settings: true,
            notification_priority: 'PRIORITY_HIGH',
            visibility: 'PUBLIC',
            icon: 'notification_icon',
            color: '#FF0000',
          },
        },
      },
    };
  }
  // iOS 전용 메시지
  else if (platform === 'ios') {
    return {
      message: {
        token: token,
        notification: {
          title: notificationTitle,
          body: notificationBody,
        },
        data: {
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
              'content-available': 1,
              'mutable-content': 1,
              sound: 'default',
              badge: 1,
            },
            fcm_options: {
              image: 'https://login.care.inhandplus.com/images/logo.png',
            },
          },
        },
      },
    };
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
      allTokens.push({ token, platform: 'android' });
    }

    // 모든 토큰에 대해 메시지 전송
    for (const { token, platform } of allTokens) {
      const message = createMessage(token, platform);
      console.log(
        `${platform} 기기로 전송할 메시지:`,
        JSON.stringify(message, null, 2),
      );

      // FCM API 호출
      try {
        const response = await axios.post(FCM_URL, message, {
          headers: {
            Authorization: `Bearer ${accessToken}`,
            'Content-Type': 'application/json',
          },
        });

        console.log(
          `메시지가 성공적으로 전송되었습니다 (${platform} 토큰: ${token.substring(
            0,
            10,
          )}...):`,
          response.data,
        );
        results.push({ token, platform, success: true, data: response.data });
      } catch (err) {
        console.error(
          `메시지 전송 오류 (${platform} 토큰: ${token.substring(0, 10)}...):`,
          err.response ? err.response.data : err.message,
        );

        if (
          err.response &&
          err.response.data.error &&
          err.response.data.error.details
        ) {
          console.error('상세 오류:', err.response.data.error.details);
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
    console.log('모든 메시지 전송 결과:', {
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
    });
  } catch (error) {
    console.error('메시지 전송 중 일반 오류:', error.message);
  }
}

// 실행
console.log('\n=== FCM 메시지 전송 프로그램 시작 ===');
sendMessage();

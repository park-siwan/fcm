const axios = require('axios');
const getAccessToken = require('./getAccessToken');
require('dotenv').config();

// FCM HTTP v1 API URL
const FCM_URL =
  'https://fcm.googleapis.com/v1/projects/inhandplus-d8fb1/messages:send';

// 메시지 데이터
const message = {
  message: {
    token:
      // 'f3IVPuhNQne8MYRBfMJse5:APA91bF1-h011coO-D3_Hqsm7ZMmo8wBWblc5XWIKubbB-ZR1rljP0g6n6fNS-UNnmCiNbMzkzURgNoLZXMYiuRwSC9ChNrTI-yQ64XwPagg-Lze7KUwWrc', // 기기의 FCM 토큰 갤럭시 s10
      // 'dVXiSNn3Q3i3aPJpm4_o4-:APA91bEAKyRLdejLvBM4IGNBjRPtp9RHm-bpoCZ2nCMqzEkwPhkO6zZw-Ko8IMOkGH4UWeg0KS4G1LerINvp3GBWQGTfk_SaMe4zNqbswOQ9E8sJph4MozI', //Medium Phone API 35
      // 'dgnhiUDzR-efNsn0sDmdal:APA91bHlY7tkQAfplq6ONBbUxZt9t6yxHrmMKXj8NJ0txe-vyjF-QQYjs__mqyclkZPRGqYik9mxCVQO51HBRMrvar5dZPGlxjfAg2YoOqGnZkFtK7lzapk', //현성
      'eaxSepE0SWuSP3y1bjdaEi:APA91bEfgQBUUSk8A39sVCaZ7ifft-6huZxnnKqErUISB5mb49eBBIVNzkw_wap0t234YOeQXAM_Yyc3IYnriDVy3_yypMxfzlaV5ZEUysggLRW6hDNdexw', //시완 웹뷰2
    // notification: {
    //   title: '응급상황(낙상감지)',
    //   body: '참여자의 낙상이 감지되었습니다.',
    // },
    data: {
      title: '응급상황(낙상감지)',
      body: '참여자의 낙상이 감지되었습니다.',
      // target_url: 'http://localhost:3000/ko/participant-info',
      target_url:
        'http://localhost:3000/ko/help/9154d7dfe98b46af9a39bbf373039fec',
      // target_url: '',
    },
  },
};

async function sendMessage() {
  try {
    // Access Token 가져오기
    const accessToken = await getAccessToken();
    // console.log(accessToken);
    // FCM API 호출
    const response = await axios.post(FCM_URL, message, {
      headers: {
        Authorization: `Bearer ${accessToken}`,
        'Content-Type': 'application/json',
      },
    });

    console.log('Message sent successfully:', response.data);
  } catch (error) {
    console.error(
      'Error sending message:',
      error.response ? error.response.data : error.message,
    );
    if (error.response) {
      console.error(error.response.data.error.details);
    }
  }
}

sendMessage();

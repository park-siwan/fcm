const admin = require('firebase-admin');

// 서비스 계정 키 JSON 파일 경로
const serviceAccount = require('./serviceAccountKey.json');

// Firebase Admin SDK 초기화
admin.initializeApp({
  credential: admin.credential.cert(serviceAccount),
});

async function getAccessToken() {
  const token = await admin.credential.applicationDefault().getAccessToken();
  return token.access_token;
}

module.exports = getAccessToken;

import http from 'k6/http';
import { check, sleep } from 'k6';
import { Trend, Rate } from 'k6/metrics';

const uploadDuration = new Trend('upload_duration', true);
const errorRate = new Rate('error_rate');  

export const options = {
  scenarios: {
    concurrent_upload_and_analyze: {
      executor: 'constant-vus',
      vus: 100,
      duration: '5m',           // можно изменить на '10m'
      gracefulStop: '30s',
    },
  },
  thresholds: {
    http_req_duration: ['p(95)<5000'],
    upload_duration: ['p(95)<3000'],
    error_rate: ['rate<0.05'],            // порог для кастомной метрики (не более 5% ошибок)
    http_req_failed: ['rate<0.05'],       // порог для HTTP-ошибок (не более 5%)
  },
};

function generateUniqueUser(vu, iter) {
  const ts = Date.now();
  const uniqueSuffix = `${ts}_${vu}_${iter}`;
  // Генерируем номер телефона из 11 цифр, начинающийся с 79
  let phoneNum = 79000000000 + (ts % 1000000000) + vu * 1000 + iter;
  let phone = phoneNum.toString().slice(0, 11);
  return {
    username: `loadtest_${uniqueSuffix}`,
    email: `loadtest_${uniqueSuffix}@test.com`,
    password: 'Test123456!',
    first_name: 'Load',
    last_name: `Tester_${vu}_${iter}`,
    phone: phone,
  };
}

export default function () {
  const vu = __VU;
  const iter = __ITER;
  const user = generateUniqueUser(vu, iter);

  // ---------- РЕГИСТРАЦИЯ ----------
  let registerRes = http.post(
    `${__ENV.API_URL}/auth/sign-up`,
    JSON.stringify({
      username: user.username,
      email: user.email,
      password: user.password,
      first_name: user.first_name,
      last_name: user.last_name,
      phone: user.phone,
    }),
    { headers: { 'Content-Type': 'application/json' } }
  );

  if (registerRes.status !== 200 && registerRes.status !== 201 && registerRes.status !== 409) {
    errorRate.add(1);   // регистрация провалилась
    return;
  }

  // ---------- ЛОГИН ----------
  let loginRes = http.post(
    `${__ENV.API_URL}/auth/sign-in`,
    JSON.stringify({
      username: user.username,
      password: user.password,
    }),
    { headers: { 'Content-Type': 'application/json' } }
  );

  check(loginRes, { 'login successful': (r) => r.status === 200 });
  if (loginRes.status !== 200) {
    errorRate.add(1);   // логин провалился
    return;
  }
  const token = loginRes.json('token');

  // ---------- ЗАГРУЗКА ДОКУМЕНТА ----------
  const pdfContent = buildTestPDF();
  const formData = {
    file: http.file(pdfContent, 'test-document.pdf', 'application/pdf'),
  };

  let uploadRes = http.post(
    `${__ENV.API_URL}/api/documents`,
    formData,
    {
      headers: {
        'Authorization': `Bearer ${token}`,
      },
    }
  );

  check(uploadRes, { 'upload successful': (r) => r.status === 202 });
  uploadDuration.add(uploadRes.timings.duration);
  if (uploadRes.status !== 202) {
    errorRate.add(1);   // загрузка документа провалилась
    return;
  }
  errorRate.add(0);   // отмечаем, что итерация прошла без ошибок
  sleep(Math.random() * 2 + 1);
}


function buildTestPDF() {
  return `%PDF-1.0
1 0 obj<</Type/Catalog/Pages 2 0 R>>endobj
2 0 obj<</Type/Pages/Kids[3 0 R]/Count 1>>endobj
3 0 obj<</Type/Page/MediaBox[0 0 612 792]/Parent 2 0 R/Resources<<>>>>endobj
xref
0 4
0000000000 65535 f 
0000000009 00000 n 
0000000058 00000 n 
0000000118 00000 n 
trailer<</Size 4/Root 1 0 R>>
startxref
179
%%EOF`;
}

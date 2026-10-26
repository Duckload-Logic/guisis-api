import http from 'k6/http';
import { check, sleep } from 'k6';

export const options = {
  stages: [
    { duration: '5s', target: 10 },   // Warm-up 10 VUs
    { duration: '15s', target: 20 },  // Baseline 20 VUs
    { duration: '10s', target: 50 },  // Stress spike 50 VUs
    { duration: '10s', target: 50 },  // Hold 50 VUs
    { duration: '5s', target: 0 },    // Ramp-down
  ],
  thresholds: {
    http_req_failed: ['rate<0.05'],    // Allow max 5% error during stress
    http_req_duration: ['p(95)<2000'], // 95% of requests under 2s
  },
};

const API_URL = __ENV.API_URL || 'http://localhost:8080';
const AI_URL = __ENV.AI_URL || 'http://localhost:8000';
const AI_API_KEY = __ENV.AI_API_KEY || '1yOojwxW0POZMQxQ';
const API_TOKEN = __ENV.API_TOKEN || '';
const IIR_ID = __ENV.IIR_ID || '1';

const CONCERNS = [
  'tulong naman po, emergency suicidal thoughts',
  'how to get student id and campus map saan po',
  'delayed po ang graduation ko may grades issue',
  'ordinary routine check for library hours please',
  'sobrang stress at depress na ako sa grades ko par',
];

export default function () {
  const rand = Math.random();

  const authHeaders = {
    headers: {
      'Content-Type': 'application/json',
      'Authorization': `Bearer ${API_TOKEN}`,
      'X-VU-ID': String(__VU),
    },
  };

  const aiHeaders = {
    headers: {
      'Content-Type': 'application/json',
      'X-API-Key': AI_API_KEY,
      'X-VU-ID': String(__VU),
    },
  };

  if (rand < 0.40) {
    // 1. AI Classification (40% weight) - CPU intensive
    const randomText = CONCERNS[Math.floor(Math.random() * CONCERNS.length)];
    const payload = JSON.stringify({ text: randomText });

    const res = http.post(
      `${AI_URL}/api/v1/classify`,
      payload,
      Object.assign({}, aiHeaders, { tags: { name: 'AI_Classify' } })
    );

    check(res, {
      'AI: status is 200': (r) => r.status === 200,
      'AI: has level': (r) => {
        try {
          const body = r.json();
          return body && body.level !== undefined;
        } catch (_) {
          return false;
        }
      },
    });

  } else if (rand < 0.70) {
    // 2. Database Lookup (30% weight) - MySQL read connection
    const res = http.get(
      `${API_URL}/api/v1/students/lookups/genders`,
      Object.assign({}, authHeaders, { tags: { name: 'DB_Lookup' } })
    );

    check(res, {
      'DB Lookup: status is 200': (r) => r.status === 200,
    });

  } else if (rand < 0.95) {
    // 3. IIR Draft Process (25% weight) - DB & Session Check
    const res = http.get(
      `${API_URL}/api/v1/students/inventory/records/iir/draft`,
      Object.assign({}, authHeaders, { tags: { name: 'IIR_Draft' } })
    );

    check(res, {
      'IIR Draft: status is 200/404': (r) => 
        r.status === 200 || r.status === 404,
    });

  } else {
    // 4. PDF Generation (5% weight) - Gotenberg / Chromium
    const res = http.get(
      `${API_URL}/api/v1/students/inventory/records/iir/${IIR_ID}/download`,
      Object.assign({}, authHeaders, { tags: { name: 'PDF_Download' } })
    );

    check(res, {
      'Gotenberg: status is 200': (r) => r.status === 200,
    });
  }

  sleep(1);
}

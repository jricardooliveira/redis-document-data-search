import http from 'k6/http';
import { check, sleep } from 'k6';

export let options = {
  vus: 50, // virtual users
  duration: '30s', // test duration
};

// Example indexed fields for customers and events:
const customerQueries = [
  'email=foo1@example.com',
  'phone=1234567890',
  'visitor_id=abc123',
];

const eventQueries = [
  'visitor_id=abc123',
  'call_id=call_001',
  'chat_id=chat_001',
  'external_id=ext_001',
  'form2lead_id=f2l_001',
  'tickets_id=ticket_001',
];

const BASE_URL = 'http://localhost:8080';

export default function () {
  // Randomly pick customer or event
  if (Math.random() < 0.5) {
    // Customer search
    const q = customerQueries[Math.floor(Math.random() * customerQueries.length)];
    let res = http.get(`${BASE_URL}/search_customers?${q}`);
    check(res, {
      'customer search status is 200': (r) => r.status === 200,
      'customer search has results': (r) => r.body && r.body.includes('results'),
    });
  } else {
    // Event search
    const q = eventQueries[Math.floor(Math.random() * eventQueries.length)];
    let res = http.get(`${BASE_URL}/search_events?${q}`);
    check(res, {
      'event search status is 200': (r) => r.status === 200,
      'event search has results': (r) => r.body && r.body.includes('results'),
    });
  }
  sleep(1);
}

import http from 'k6/http';
import { check, sleep } from 'k6';
import { SharedArray } from 'k6/data';

export let options = {
  vus: 50,
  duration: '30s',
};

// Load CSVs into memory (relative to k6 run directory)
const customerSamples = new SharedArray('customers', function () {
  return open('./customer_sample.csv')
    .split('\n')
    .slice(1) // skip header
    .filter(Boolean)
    .map(line => {
      const [key, email, phone, visitor_id] = line.split(',').map(s => s.replace(/^"|"$/g, ''));
      return { key, email, phone, visitor_id };
    });
});

const eventSamples = new SharedArray('events', function () {
  return open('./event_sample.csv')
    .split('\n')
    .slice(1)
    .filter(Boolean)
    .map(line => {
      const [key, visitor_id, call_id, chat_id, external_id, form2lead_id, tickets_id] = line.split(',').map(s => s.replace(/^"|"$/g, ''));
      return { key, visitor_id, call_id, chat_id, external_id, form2lead_id, tickets_id };
    });
});

const BASE_URL = 'http://localhost:8080';

function pickRandom(arr) {
  return arr[Math.floor(Math.random() * arr.length)];
}

export default function () {
  if (Math.random() < 0.5 && customerSamples.length > 0) {
    // Random customer sample and indexed field
    const rec = pickRandom(customerSamples);
    const fields = ['email', 'phone', 'visitor_id'];
    const field = pickRandom(fields);
    const value = rec[field];
    if (value) {
      let res = http.get(`${BASE_URL}/search_customers?${field}=${encodeURIComponent(value)}`);
      check(res, {
        'customer search status is 200': (r) => r.status === 200,
        'customer search has results': (r) => r.body && r.body.includes('results'),
      });
    }
  } else if (eventSamples.length > 0) {
    // Random event sample and indexed field
    const rec = pickRandom(eventSamples);
    const fields = ['visitor_id', 'call_id', 'chat_id', 'external_id', 'form2lead_id', 'tickets_id'];
    const field = pickRandom(fields);
    const value = rec[field];
    if (value) {
      let res = http.get(`${BASE_URL}/search_events?${field}=${encodeURIComponent(value)}`);
      check(res, {
        'event search status is 200': (r) => r.status === 200,
        'event search has results': (r) => r.body && r.body.includes('results'),
      });
    }
  }
  sleep(1);
}

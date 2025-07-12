import http from 'k6/http';
import { check, sleep } from 'k6';
import { SharedArray } from 'k6/data';

export let options = {
  vus: 50,
  duration: '30s',
};

const BASE_URL = 'http://localhost:8080';

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

export function setup() {
    // Query /healthz before test starts
    const res = http.get(`${BASE_URL}/healthz`);
    if (res.status === 200) {
      try {
        const health = JSON.parse(res.body);
        console.log('----- DATABASE DEBUG INFO -----');
        console.log(`Customer count: ${health.customer_count}`);
        console.log(`Event count:    ${health.event_count}`);
        if (health.redis_memory) {
          console.log(`Redis/Valkey used memory: ${health.redis_memory.used_memory_bytes} bytes (${health.redis_memory.used_memory_human})`);
        }
        console.log('--------------------------------');
      } catch (e) {
        console.log('Error parsing /healthz response:', e);
        console.log('Raw /healthz response body:', res.body);
      }
    } else {
      console.log('/healthz endpoint failed with status', res.status);
      console.log('Raw /healthz response body:', res.body);
    }
  }

function pickRandom(arr) {
  return arr[Math.floor(Math.random() * arr.length)];
}

export default function () {
  if (Math.random() < 0.5 && customerSamples.length > 0) {
    // Random customer sample and indexed field (only non-empty fields)
    const rec = pickRandom(customerSamples);
    const fields = ['email', 'phone', 'visitor_id'];
    const nonEmptyFields = fields.filter(f => rec[f] && rec[f].trim() !== '' && rec[f].toLowerCase() !== 'null');
    if (nonEmptyFields.length > 0) {
      const field = pickRandom(nonEmptyFields);
      const value = rec[field];
      const url = `${BASE_URL}/search_customers?${field}=${encodeURIComponent(value)}`;
      let res = http.get(url);
      if (res.status !== 200) {
        console.log(`[CUSTOMER FAIL] Status: ${res.status} URL: ${url} Body: ${res.body && res.body.slice(0, 200)}`);
      }
      check(res, {
        'customer search status is 200': (r) => r.status === 200,
        'customer search has results': (r) => r.body && r.body.includes('results'),
      });
    }
  } else if (eventSamples.length > 0) {
    // Random event sample and indexed field (only non-empty fields)
    const rec = pickRandom(eventSamples);
    const fields = ['visitor_id', 'call_id', 'chat_id', 'external_id', 'form2lead_id', 'tickets_id'];
    const nonEmptyFields = fields.filter(f => rec[f] && rec[f].trim() !== '' && rec[f].toLowerCase() !== 'null');
    if (nonEmptyFields.length > 0) {
      const field = pickRandom(nonEmptyFields);
      const value = rec[field];
      const url = `${BASE_URL}/search_events?${field}=${encodeURIComponent(value)}`;
      let res = http.get(url);
      if (res.status !== 200) {
        console.log(`[EVENT FAIL] Status: ${res.status} URL: ${url} Body: ${res.body && res.body.slice(0, 200)}`);
      }
      check(res, {
        'event search status is 200': (r) => r.status === 200,
        'event search has results': (r) => r.body && r.body.includes('results'),
      });
    }
  }
  sleep(1);
}

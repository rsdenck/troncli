import http from 'k6/http';
import { check, sleep } from 'k6';

const OLLAMA_HOST = 'http://192.168.130.25:11434';
const NVIDIA_API = 'https://integrate.api.nvidia.com/v1/chat/completions';
const NVIDIA_KEY = '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';

export default function() {
  // Test Ollama Generate
  const ollamaPayload = JSON.stringify({
    model: 'gemma:2b',
    prompt: 'test load',
    stream: false
  });
  
  const ollamaRes = http.post(`${OLLAMA_HOST}/api/generate`, ollamaPayload, {
    headers: { 'Content-Type': 'application/json' },
    timeout: '360s'
  });
  
  check(ollamaRes, {
    'Ollama status is 200': (r) => r.status === 200,
    'Ollama has response': (r) => r.body.includes('response'),
  });
  
  sleep(1);
  
  // Test Ollama Tags
  const tagsRes = http.get(`${OLLAMA_HOST}/api/tags`, { timeout: '30s' });
  check(tagsRes, {
    'Ollama Tags status is 200': (r) => r.status === 200,
  });
  
  sleep(1);
  
  // Test NVIDIA Build
  const nvidiaPayload = JSON.stringify({
    model: 'minimaxai/minimax-m2.7',
    messages: [{ role: 'user', content: 'test load' }],
    max_tokens: 50
  });
  
  const nvidiaRes = http.post(NVIDIA_API, nvidiaPayload, {
    headers: {
      'Content-Type': 'application/json',
      'Authorization': `Bearer ${NVIDIA_KEY}`
    },
    timeout: '60s'
  });
  
  check(nvidiaRes, {
    'NVIDIA status is 200': (r) => r.status === 200,
    'NVIDIA has choices': (r) => r.body.includes('choices'),
  });
}

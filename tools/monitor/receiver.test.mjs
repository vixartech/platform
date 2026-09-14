import assert from 'node:assert/strict';
import { once } from 'node:events';
import { test } from 'node:test';
import { newReceiver } from './receiver.mjs';

async function fixture(t) {
  const server = newReceiver();
  server.listen(0, '127.0.0.1');
  await once(server, 'listening');
  t.after(
    () =>
      new Promise((resolve) => {
        server.close(resolve);
        server.closeAllConnections();
      }),
  );
  const url = `http://127.0.0.1:${server.address().port}`;
  const post = (body) =>
    fetch(`${url}/alerts`, {
      method: 'POST',
      body: JSON.stringify(body),
      signal: AbortSignal.timeout(2000),
    });
  const events = async () => {
    const response = await fetch(`${url}/events`);
    assert.equal(response.headers.get('cache-control'), 'no-store');
    return response.json();
  };
  return { server, post, events, url };
}

test('receiver retains at most 100 metadata events, not payload/labels/annotations', async (t) => {
  const { server, post, events } = await fixture(t);
  assert.equal(server.maxConnections, 8);
  assert.equal(server.requestTimeout, 5000);
  const alert = (name) => ({
    status: 'firing',
    labels: { alertname: name, secret: 'PRIVATE' },
    annotations: { summary: 'PRIVATE' },
  });
  assert.equal(
    (await post({ alerts: Array.from({ length: 100 }, () => alert('First')) })).status,
    204,
  );
  assert.equal(
    (await post({ alerts: [alert('Last'), { ...alert('Last'), status: 'resolved' }] })).status,
    204,
  );
  const result = await events();
  assert.equal(result.length, 100);
  assert.equal(result.at(-1).status, 'resolved');
  assert.equal(result.at(-1).name, 'Last');
  assert.ok(!JSON.stringify(result).includes('PRIVATE'));
  assert.deepEqual(Object.keys(result[0]).sort(), ['at', 'name', 'status']);
});

test('invalid batches are rejected atomically and routes are narrow', async (t) => {
  const { post, events, url } = await fixture(t);
  const valid = { status: 'firing', labels: { alertname: 'Valid' } };
  for (const invalid of [
    null,
    {},
    { alerts: null },
    { alerts: Array(101).fill(valid) },
    { alerts: [valid, null] },
    { alerts: [{ ...valid, status: 'unknown' }] },
    { alerts: [{ status: 'firing', labels: { alertname: '<script>' } }] },
  ]) {
    assert.equal((await post(invalid)).status, 400);
    assert.deepEqual(await events(), []);
  }
  assert.equal((await fetch(`${url}/alerts`)).status, 404);
  assert.equal((await fetch(`${url}/events?debug=1`)).status, 404);
  assert.equal((await fetch(`${url}/alerts`, { method: 'POST', body: '{broken' })).status, 400);
  assert.deepEqual(await events(), []);
});

test('oversize bodies cannot enter the retained event buffer', async (t) => {
  const { post, events } = await fixture(t);
  try {
    const response = await post({ alerts: [], padding: 'x'.repeat(65536) });
    assert.equal(response.status, 413);
  } catch (error) {
    // The receiver may close an over-limit connection before the response is read.
    assert.equal(error.cause?.code, 'UND_ERR_SOCKET');
  }
  assert.deepEqual(await events(), []);
});

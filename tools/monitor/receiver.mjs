// Local acceptance sink only: bounded in-memory metadata, never raw webhook bodies.
import { createServer } from 'node:http';
import { resolve } from 'node:path';
import { pathToFileURL } from 'node:url';

export function newReceiver() {
  const events = [];
  const server = createServer(
    { requestTimeout: 5000, headersTimeout: 3000, maxHeaderSize: 8192 },
    async (req, res) => {
      res.setHeader('Cache-Control', 'no-store');
      if (req.method === 'GET' && req.url === '/events') {
        res.setHeader('Content-Type', 'application/json');
        res.end(JSON.stringify(events));
        return;
      }
      if (req.method !== 'POST' || req.url !== '/alerts') {
        res.writeHead(404).end();
        return;
      }
      try {
        let bytes = 0;
        const chunks = [];
        for await (const chunk of req) {
          bytes += chunk.length;
          if (bytes > 65536) {
            res.writeHead(413).end();
            req.destroy();
            return;
          }
          chunks.push(chunk);
        }
        const body = JSON.parse(Buffer.concat(chunks).toString('utf8'));
        if (!Array.isArray(body.alerts) || body.alerts.length > 100)
          throw new Error('invalid alerts');
        const next = body.alerts.map((alert) => {
          if (
            !['firing', 'resolved'].includes(alert.status) ||
            typeof alert.labels?.alertname !== 'string' ||
            !/^[a-zA-Z][a-zA-Z0-9_]{0,79}$/.test(alert.labels.alertname)
          )
            throw new Error('invalid alert');
          return {
            name: alert.labels.alertname,
            status: alert.status,
            at: new Date().toISOString(),
          };
        });
        events.push(...next);
        events.splice(0, Math.max(0, events.length - 100));
        res.writeHead(204).end();
      } catch {
        res.writeHead(400).end();
      }
    },
  );
  server.maxConnections = 8;
  return server;
}

if (process.argv[1] && import.meta.url === pathToFileURL(resolve(process.argv[1])).href) {
  const server = newReceiver();
  server.listen(8080, '0.0.0.0');
  process.on('SIGTERM', () => {
    server.close();
    server.closeAllConnections();
  });
}

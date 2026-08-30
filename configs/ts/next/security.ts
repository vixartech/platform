export type SecurityHeader = {
  key: string;
  value: string;
};

export type SecurityHeadersOptions = {
  development: boolean;
  cacheControl?: string;
  childSources?: readonly string[];
  connectSources?: readonly string[];
  frameAncestors?: readonly string[];
  frameSources?: readonly string[];
  imageSources?: readonly string[];
  referrerPolicy?: 'no-referrer' | 'strict-origin-when-cross-origin';
  scriptSources?: readonly string[];
  strictTransportSecurity?: boolean;
};

export function createSecurityHeaders(options: SecurityHeadersOptions): SecurityHeader[] {
  const headers: SecurityHeader[] = [];

  if (options.cacheControl) {
    headers.push({ key: 'Cache-Control', value: options.cacheControl });
  }

  headers.push(
    { key: 'Content-Security-Policy', value: contentSecurityPolicy(options) },
    { key: 'Cross-Origin-Opener-Policy', value: 'same-origin' },
    { key: 'Cross-Origin-Resource-Policy', value: 'same-origin' },
    { key: 'Permissions-Policy', value: 'camera=(), geolocation=(), microphone=()' },
    {
      key: 'Referrer-Policy',
      value: options.referrerPolicy ?? 'strict-origin-when-cross-origin',
    },
    { key: 'X-Content-Type-Options', value: 'nosniff' },
  );

  if (!options.frameAncestors?.length) {
    headers.push({ key: 'X-Frame-Options', value: 'DENY' });
  }

  if (options.strictTransportSecurity ?? !options.development) {
    headers.push({
      key: 'Strict-Transport-Security',
      value: 'max-age=63072000; includeSubDomains',
    });
  }

  return headers;
}

function contentSecurityPolicy(options: SecurityHeadersOptions): string {
  const childSources = ["'self'", ...(options.childSources ?? [])];
  const connectSources = [
    "'self'",
    ...(options.connectSources ?? []),
    ...(options.development ? ['http:', 'ws:'] : []),
  ];
  const frameAncestors = options.frameAncestors?.length
    ? options.frameAncestors
    : ["'none'"];
  const frameSources = ["'self'", ...(options.frameSources ?? [])];
  const imageSources = ["'self'", 'data:', ...(options.imageSources ?? [])];
  const scriptSources = [
    "'self'",
    "'unsafe-inline'",
    ...(options.scriptSources ?? []),
    ...(options.development ? ["'unsafe-eval'"] : []),
  ];

  return [
    "base-uri 'self'",
    "default-src 'self'",
    "form-action 'self'",
    `frame-ancestors ${frameAncestors.join(' ')}`,
    "object-src 'none'",
    `child-src ${childSources.join(' ')}`,
    `connect-src ${connectSources.join(' ')}`,
    "font-src 'self' data:",
    `frame-src ${frameSources.join(' ')}`,
    `img-src ${imageSources.join(' ')}`,
    `script-src ${scriptSources.join(' ')}`,
    "style-src 'self' 'unsafe-inline'",
  ].join('; ');
}

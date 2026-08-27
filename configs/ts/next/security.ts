export type SecurityHeader = {
  key: string;
  value: string;
};

export type SecurityHeadersOptions = {
  development: boolean;
  cacheControl?: string;
  imageSources?: readonly string[];
  referrerPolicy?: 'no-referrer' | 'strict-origin-when-cross-origin';
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
    { key: 'X-Frame-Options', value: 'DENY' },
  );

  if (options.strictTransportSecurity ?? !options.development) {
    headers.push({
      key: 'Strict-Transport-Security',
      value: 'max-age=63072000; includeSubDomains',
    });
  }

  return headers;
}

function contentSecurityPolicy(options: SecurityHeadersOptions): string {
  const connectSources = [
    "'self'",
    ...(options.development ? ['http:', 'ws:'] : []),
  ];
  const imageSources = ["'self'", 'data:', ...(options.imageSources ?? [])];

  return [
    "base-uri 'self'",
    "default-src 'self'",
    "form-action 'self'",
    "frame-ancestors 'none'",
    "object-src 'none'",
    `connect-src ${connectSources.join(' ')}`,
    "font-src 'self' data:",
    `img-src ${imageSources.join(' ')}`,
    `script-src 'self' 'unsafe-inline'${options.development ? " 'unsafe-eval'" : ''}`,
    "style-src 'self' 'unsafe-inline'",
  ].join('; ');
}

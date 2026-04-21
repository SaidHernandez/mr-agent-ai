// Standard error response shape — use this in every controller error path.
// See: skills/controller/SKILL.md

export interface ErrorResponse {
  success: false;
  code: string;      // SCREAMING_SNAKE_CASE — machine-readable
  message: string;   // Human-readable, safe to display in UI
  detail: string;    // Technical context for debugging — never PII
}

export interface SuccessResponse<T = unknown> {
  success: true;
  data: T;
}

export type ApiResponse<T = unknown> = SuccessResponse<T> | ErrorResponse;

// ── Error code registry ───────────────────────────────────────────────────────
export const ERROR_CODES = {
  VALIDATION_ERROR:      'VALIDATION_ERROR',
  UNAUTHORIZED:          'UNAUTHORIZED',
  FORBIDDEN:             'FORBIDDEN',
  NOT_FOUND:             'NOT_FOUND',
  CONFLICT:              'CONFLICT',
  UNPROCESSABLE:         'UNPROCESSABLE',
  RATE_LIMIT_EXCEEDED:   'RATE_LIMIT_EXCEEDED',
  INTERNAL_ERROR:        'INTERNAL_ERROR',
} as const;

export type ErrorCode = (typeof ERROR_CODES)[keyof typeof ERROR_CODES];

// ── HTTP status map ───────────────────────────────────────────────────────────
export const HTTP_STATUS: Record<ErrorCode, number> = {
  VALIDATION_ERROR:    400,
  UNAUTHORIZED:        401,
  FORBIDDEN:           403,
  NOT_FOUND:           404,
  CONFLICT:            409,
  UNPROCESSABLE:       422,
  RATE_LIMIT_EXCEEDED: 429,
  INTERNAL_ERROR:      500,
};

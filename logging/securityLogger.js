
const winston = require('winston');

const securityLogger = winston.createLogger({
  level: 'info',
  format: winston.format.combine(
    winston.format.timestamp(),
    winston.format.json()
  ),
  transports: [
    new winston.transports.File({ filename: 'logs/security.log' }),
    new winston.transports.Console(),
  ],
});


function logFailedLogin(email, ip) {
  securityLogger.warn('FAILED_LOGIN_ATTEMPT', {
    email,
    ip,
    time: new Date().toISOString(),
  });
}

function logUnauthorizedAccess(userId, role, route, ip) {
  securityLogger.warn('UNAUTHORIZED_ACCESS_ATTEMPT', {
    userId,
    role,
    route,
    ip,
    time: new Date().toISOString(),
  });
}

function logSuspiciousInput(field, ip, route) {
  securityLogger.warn('SUSPICIOUS_INPUT_BLOCKED', {
    field, 
    route,
    ip,
    time: new Date().toISOString(),
  });
}

function logRateLimitHit(ip, route) {
  securityLogger.warn('RATE_LIMIT_EXCEEDED', {
    ip,
    route,
    time: new Date().toISOString(),
  });
}

module.exports = {
  securityLogger,
  logFailedLogin,
  logUnauthorizedAccess,
  logSuspiciousInput,
  logRateLimitHit,
};

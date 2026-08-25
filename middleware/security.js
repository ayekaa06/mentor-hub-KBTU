
const helmet = require('helmet');
const rateLimit = require('express-rate-limit');
const cors = require('cors');
const xss = require('xss-clean');
const mongoSanitize = require('express-mongo-sanitize'); 

function applySecurity(app) {

  app.use(helmet());

  
  app.use(
    cors({
      origin: process.env.FRONTEND_URL, 
      credentials: true,
    })
  );

  
  const limiter = rateLimit({
    windowMs: 15 * 60 * 1000, // 15 минут
    max: 100, 
    message: { error: 'Слишком много запросов, попробуйте позже' },
  });
  app.use('/api/', limiter);

  
  const loginLimiter = rateLimit({
    windowMs: 15 * 60 * 1000,
    max: 5,
    message: { error: 'Слишком много попыток входа, попробуйте через 15 минут' },
  });
  app.use('/api/auth/login', loginLimiter);

  
  app.use(xss());

  
  app.use(mongoSanitize());

  
}


module.exports = { applySecurity };


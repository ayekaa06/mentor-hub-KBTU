
const speakeasy = require('speakeasy');
const QRCode = require('qrcode');


async function generate2FASecret(adminEmail) {
  const secret = speakeasy.generateSecret({
    name: `MentorHub KBTU (${adminEmail})`,
  });

  const qrCodeImage = await QRCode.toDataURL(secret.otpauth_url);

  return {
    secret: secret.base32,
    qrCodeImage,
  };
}


function verify2FACode(userSecret, codeFromUser) {
  return speakeasy.totp.verify({
    secret: userSecret, 
    encoding: 'base32',
    token: codeFromUser, 
    window: 1, 
  });
}


function require2FA(req, res, next) {
  const code = req.headers['x-2fa-code'];

  if (!code) {
    return res.status(401).json({ error: 'Требуется код двухфакторной аутентификации' });
  }

  
  const adminSecret = req.user.twoFactorSecret;

  const isValid = verify2FACode(adminSecret, code);

  if (!isValid) {
    return res.status(403).json({ error: 'Неверный код 2FA' });
  }

  next();
}

module.exports = { generate2FASecret, verify2FACode, require2FA };

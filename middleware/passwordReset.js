
const crypto = require('crypto');
const { hashPassword } = require('../utils/password');

// В реальном проекте это должно храниться в БД (таблица password_resets),
// здесь — упрощённая in-memory версия для демонстрации логики.
const resetTokens = new Map(); // token -> { userId, expiresAt }

const TOKEN_TTL_MS = 15 * 60 * 1000; // 15 минут


function requestPasswordReset(userId, sendEmailFn, userEmail) {
  const token = crypto.randomBytes(32).toString('hex'); // случайный токен, не угадать
  const expiresAt = Date.now() + TOKEN_TTL_MS;

  resetTokens.set(token, { userId, expiresAt });

  const resetLink = `${process.env.FRONTEND_URL}/reset-password?token=${token}`;

  if (sendEmailFn) {
    sendEmailFn(userEmail, resetLink);
  }

  return { message: 'Если email существует, ссылка для сброса пароля отправлена' };
}


async function confirmPasswordReset(token, newPassword) {
  const record = resetTokens.get(token);

  if (!record) {
    throw new Error('Токен недействителен или уже использован');
  }

  if (Date.now() > record.expiresAt) {
    resetTokens.delete(token);
    throw new Error('Срок действия токена истёк, запросите сброс пароля заново');
  }

  const newHash = await hashPassword(newPassword);

  resetTokens.delete(token); 

  return { userId: record.userId, newPasswordHash: newHash };
}

module.exports = { requestPasswordReset, confirmPasswordReset };

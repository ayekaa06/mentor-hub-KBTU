process.env.JWT_ACCESS_SECRET = 'test_secret_1';
process.env.JWT_REFRESH_SECRET = 'test_secret_2';

const { hashPassword, comparePassword, isPasswordStrong } = require('./utils/password');
const { generateTokens, verifyToken } = require('./middleware/auth');
const { requireRole } = require('./middleware/rbac');

async function runTests() {
  console.log('--- Проверка пароля ---');
  const hashed = await hashPassword('MyPassword123');
  console.log('Хеш:', hashed);
  console.log('Пароль совпадает:', await comparePassword('MyPassword123', hashed));
  console.log('Неверный пароль:', await comparePassword('wrong', hashed));
  console.log('Пароль сложный:', isPasswordStrong('MyPassword123'));

  console.log('\n--- Проверка JWT ---');
  const user = { id: '123', role: 'student' };
  const { accessToken } = generateTokens(user);
  console.log('Токен создан:', accessToken.slice(0, 30) + '...');

  console.log('\n--- Проверка RBAC ---');
  const fakeReq = { user: { id: '123', role: 'student' } };
  const fakeRes = {
    status: (code) => ({ json: (msg) => console.log('Ответ:', code, msg) }),
  };
  const middleware = requireRole('admin');
  middleware(fakeReq, fakeRes, () => console.log('Доступ разрешён (не должно быть)'));
}

runTests();

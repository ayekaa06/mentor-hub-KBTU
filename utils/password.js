
const bcrypt = require('bcrypt');


const SALT_ROUNDS = 12;


async function hashPassword(plainPassword) {
  return await bcrypt.hash(plainPassword, SALT_ROUNDS);
}


async function comparePassword(plainPassword, hashedPassword) {
  return await bcrypt.compare(plainPassword, hashedPassword);
}


function isPasswordStrong(password) {
  const regex = /^(?=.*[a-z])(?=.*[A-Z])(?=.*\d).{8,}$/;
  return regex.test(password);
}

module.exports = { hashPassword, comparePassword, isPasswordStrong };

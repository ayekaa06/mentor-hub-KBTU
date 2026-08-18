
const jwt = require('jsonwebtoken');

const ACCESS_SECRET = process.env.JWT_ACCESS_SECRET;
const REFRESH_SECRET = process.env.JWT_REFRESH_SECRET;


function generateTokens(user) {
  const payload = {
    id: user.id,
    role: user.role, // 'student' | 'mentor' | 'admin'
  };

  const accessToken = jwt.sign(payload, ACCESS_SECRET, { expiresIn: '15m' });
  const refreshToken = jwt.sign(payload, REFRESH_SECRET, { expiresIn: '7d' });

  return { accessToken, refreshToken };
}


function verifyToken(req, res, next) {
  const authHeader = req.headers.authorization; // "Bearer <token>"

  if (!authHeader || !authHeader.startsWith('Bearer ')) {
    return res.status(401).json({ error: 'Токен не предоставлен' });
  }

  const token = authHeader.split(' ')[1];

  try {
    const decoded = jwt.verify(token, ACCESS_SECRET);
    req.user = decoded; // { id, role }
    next();
  } catch (err) {
    if (err.name === 'TokenExpiredError') {
      return res.status(401).json({ error: 'Токен истёк, обновите его' });
    }
    return res.status(403).json({ error: 'Недействительный токен' });
  }
}


function refreshAccessToken(refreshToken) {
  const decoded = jwt.verify(refreshToken, REFRESH_SECRET);
  const newAccessToken = jwt.sign(
    { id: decoded.id, role: decoded.role },
    ACCESS_SECRET,
    { expiresIn: '15m' }
  );
  return newAccessToken;
}

module.exports = { generateTokens, verifyToken, refreshAccessToken };

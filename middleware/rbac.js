
function requireRole(...allowedRoles) {
  return (req, res, next) => {
    if (!req.user) {
      return res.status(401).json({ error: 'Не авторизован' });
    }

    if (!allowedRoles.includes(req.user.role)) {
      return res.status(403).json({
        error: 'Недостаточно прав для выполнения этого действия',
      });
    }

    next();
  };
}


function requireOwnershipOrAdmin(req, res, next) {
  const targetUserId = req.params.userId || req.params.id;

  if (req.user.role === 'admin' || req.user.id === targetUserId) {
    return next();
  }

  return res.status(403).json({ error: 'Доступ запрещён' });
}

module.exports = { requireRole, requireOwnershipOrAdmin };

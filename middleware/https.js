
function enforceHttps(req, res, next) {
  const isSecure =
    req.secure || req.headers['x-forwarded-proto'] === 'https';

  if (!isSecure && process.env.NODE_ENV === 'production') {
    return res.redirect(301, `https://${req.headers.host}${req.url}`);
  }
  next();
}

module.exports = { enforceHttps };


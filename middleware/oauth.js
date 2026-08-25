
const passport = require('passport');
const GoogleStrategy = require('passport-google-oauth20').Strategy;
const { generateTokens } = require('./auth');

passport.use(
  new GoogleStrategy(
    {
      clientID: process.env.GOOGLE_CLIENT_ID,
      clientSecret: process.env.GOOGLE_CLIENT_SECRET,
      callbackURL: '/api/auth/google/callback',
    },
    async (accessToken, refreshToken, profile, done) => {
      try {
        const email = profile.emails[0].value;

        // Разрешаем вход только с университетской почты
        if (!email.endsWith('@kbtu.kz')) {
          return done(null, false, { message: 'Разрешён вход только с почты @kbtu.kz' });
        }

        
        return done(null, { email, name: profile.displayName });
      } catch (err) {
        return done(err, null);
      }
    }
  )
);


module.exports = passport;

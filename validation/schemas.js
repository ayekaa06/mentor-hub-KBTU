
const Joi = require('joi');


const registerSchema = Joi.object({
  fullName: Joi.string().trim().min(2).max(100).required().messages({
    'string.min': 'Имя слишком короткое',
    'any.required': 'Имя обязательно',
  }),

  email: Joi.string()
    .trim()
    .lowercase()
    .email()
    .pattern(/@kbtu\.kz$/)
    .required()
    .messages({
      'string.pattern.base': 'Регистрация доступна только с почты @kbtu.kz',
      'string.email': 'Некорректный email',
    }),

  password: Joi.string()
    .min(8)
    .pattern(/^(?=.*[a-z])(?=.*[A-Z])(?=.*\d)/)
    .required()
    .messages({
      'string.pattern.base':
        'Пароль должен содержать минимум 8 символов, заглавную, строчную буквы и цифру',
    }),

  role: Joi.string().valid('student', 'mentor').required(),
  course: Joi.number().integer().min(1).max(6).when('role', {
    is: 'student',
    then: Joi.required(),
  }),
});


const loginSchema = Joi.object({
  email: Joi.string().trim().lowercase().email().required(),
  password: Joi.string().required(),
});


const mentorRequestSchema = Joi.object({
  mentorId: Joi.string().uuid().required(),
  message: Joi.string().trim().max(500).allow('').required(),
  goal: Joi.string()
    .valid('adaptation', 'career', 'academic', 'other')
    .required(),
});


const reviewSchema = Joi.object({
  mentorId: Joi.string().uuid().required(),
  rating: Joi.number().integer().min(1).max(5).required(),
  comment: Joi.string().trim().max(1000).allow(''),
});


function validate(schema) {
  return (req, res, next) => {
    const { error, value } = schema.validate(req.body, {
      abortEarly: false, 
      stripUnknown: true, 
    });

    if (error) {
      return res.status(400).json({
        error: 'Ошибка валидации',
        details: error.details.map((d) => d.message),
      });
    }

    req.body = value; 
    next();
  };
}

module.exports = {
  registerSchema,
  loginSchema,
  mentorRequestSchema,
  reviewSchema,
  validate,
};

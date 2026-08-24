//
//  LocalizationManager.swift
//  MentorsHub
//
//  Глобальный менеджер языка интерфейса. Юзер переключает язык на
//  экране профиля — currentLanguage сохраняется в UserDefaults и
//  применяется во всём приложении сразу (не только для хэндбука).
//
//  Как пользоваться в View:
//      private let localization = LocalizationManager.shared
//      Text(localization.text("mentee_calendar_section"))
//

import Foundation

@Observable
class LocalizationManager {
    static let shared = LocalizationManager()

    var currentLanguage: Language {
        didSet { UserDefaults.standard.set(currentLanguage.rawValue, forKey: "appLanguage") }
    }

    private init() {
        if let saved = UserDefaults.standard.string(forKey: "appLanguage"),
           let lang = Language(rawValue: saved) {
            currentLanguage = lang
        } else {
            currentLanguage = .ru
        }
    }

    func text(_ key: String) -> String {
        strings[key]?[currentLanguage] ?? key
    }

    private let strings: [String: [Language: String]] = [
        // Auth
        "auth_title": [.ru: "Mentors Hub", .kz: "Mentors Hub", .en: "Mentors Hub"],
        "auth_username_placeholder": [.ru: "Имя пользователя", .kz: "Пайдаланушы аты", .en: "Username"],
        "auth_password_placeholder": [.ru: "Пароль", .kz: "Құпия сөз", .en: "Password"],
        "auth_login_button": [.ru: "Войти", .kz: "Кіру", .en: "Log in"],
        "auth_login_error": [.ru: "Неверный email или пароль", .kz: "Email немесе құпия сөз қате", .en: "Invalid email or password"],
        "auth_register_button": [.ru: "Регистрация", .kz: "Тіркелу", .en: "Sign up"],
        "auth_forgot_button": [.ru: "Забыли пароль?", .kz: "Құпия сөзді ұмыттыңыз ба?", .en: "Forgot password?"],

        // Registration
        "register_title": [.ru: "Регистрация", .kz: "Тіркелу", .en: "Sign up"],
        "register_name_placeholder": [.ru: "Имя", .kz: "Аты", .en: "Name"],
        "register_email_placeholder": [.ru: "Электронная почта", .kz: "Электрондық пошта", .en: "Email"],
        "register_password_placeholder": [.ru: "Пароль", .kz: "Құпия сөз", .en: "Password"],
        "register_passcheck_placeholder": [.ru: "Повторите пароль", .kz: "Құпия сөзді қайталаңыз", .en: "Repeat password"],
        "register_role_label": [.ru: "Роль", .kz: "Рөл", .en: "Role"],
        "register_specialty_placeholder": [.ru: "Специальность (например 26BDIS)", .kz: "Мамандық (мысалы 26BDIS)", .en: "Specialty (e.g. 26BDIS)"],
        "register_submit_button": [.ru: "Зарегистрироваться", .kz: "Тіркелу", .en: "Create account"],
        "register_error_exists": [.ru: "Пользователь с таким email уже существует", .kz: "Бұл email тіркелген", .en: "This email is already registered"],
        "register_have_account": [.ru: "Уже есть аккаунт? Войти", .kz: "Аккаунт бар ма? Кіру", .en: "Already have an account? Log in"],

        // Forgot password
        "forgot_title": [.ru: "Восстановите пароль", .kz: "Құпия сөзді қалпына келтіріңіз", .en: "Reset password"],
        "forgot_email_placeholder": [.ru: "Электронная почта", .kz: "Электрондық пошта", .en: "Email"],
        "forgot_newpass_placeholder": [.ru: "Придумайте новый пароль", .kz: "Жаңа құпия сөз ойлап табыңыз", .en: "Create a new password"],
        "forgot_confirm_placeholder": [.ru: "Подтвердите новый пароль", .kz: "Жаңа құпия сөзді растаңыз", .en: "Confirm new password"],
        "forgot_submit_button": [.ru: "Подтвердить", .kz: "Растау", .en: "Confirm"],

        // Profile / settings
        "profile_title": [.ru: "Профиль", .kz: "Профиль", .en: "Profile"],
        "profile_language_section": [.ru: "Язык приложения", .kz: "Қосымша тілі", .en: "App language"],
        "profile_logout_button": [.ru: "Выйти", .kz: "Шығу", .en: "Log out"],

        // Roles
        "role_mentee": [.ru: "Менти", .kz: "Менти", .en: "Mentee"],
        "role_mentor": [.ru: "Ментор", .kz: "Ментор", .en: "Mentor"],
        "role_advisor": [.ru: "Эдвайзер", .kz: "Эдвайзер", .en: "Advisor"],
        "role_vice": [.ru: "Вице-хэд", .kz: "Вице-хэд", .en: "Vice head"],
        "role_head": [.ru: "Хэд", .kz: "Хэд", .en: "Head"],
        "role_admin": [.ru: "Админ", .kz: "Админ", .en: "Admin"],

        // Mentee screen
        "mentee_main_title": [.ru: "Главная", .kz: "Басты бет", .en: "Home"],
        "mentee_rup_section": [.ru: "Учебный план", .kz: "Оқу жоспары", .en: "Curriculum"],
        "mentee_rup_empty": [.ru: "РУП пока не загружен", .kz: "ОЖ әлі жүктелмеген", .en: "Curriculum not uploaded yet"],
        "mentee_specialty_label": [.ru: "Специальность", .kz: "Мамандық", .en: "Specialty"],
        "mentee_handbook_section": [.ru: "Хэндбук", .kz: "Анықтамалық", .en: "Handbook"],
        "mentee_handbook_language": [.ru: "Язык хэндбука", .kz: "Анықтамалық тілі", .en: "Handbook language"],
        "mentee_open_handbook": [.ru: "Открыть хэндбук", .kz: "Анықтамалықты ашу", .en: "Open handbook"],
        "mentee_handbook_unavailable": [.ru: "Хэндбук на этом языке недоступен", .kz: "Бұл тілде анықтамалық жоқ", .en: "Handbook not available in this language"],
        "mentee_calendar_section": [.ru: "Академический календарь", .kz: "Академиялық күнтізбе", .en: "Academic calendar"],

        // Mentor screen
        "mentor_main_title": [.ru: "Главная", .kz: "Басты бет", .en: "Home"],
        "mentor_mentees_section": [.ru: "Мои менти", .kz: "Менің ментилерім", .en: "My mentees"],
        "mentor_randomizer_button": [.ru: "Получить менти (рандомайзер)", .kz: "Менти алу", .en: "Get a mentee (randomizer)"],
        "mentor_meetings_section": [.ru: "Встречи", .kz: "Кездесулер", .en: "Meetings"],
        "mentor_checkin_link": [.ru: "Отметить встречу за семестр", .kz: "Семестрлік кездесуді белгілеу", .en: "Check in this semester's meeting"],
        "mentor_rules_section": [.ru: "Правила менторства", .kz: "Менторлық ережелері", .en: "Mentorship rules"],
        "mentor_specialty_label": [.ru: "Моя специальность", .kz: "Менің мамандығым", .en: "My specialty"],
        "mentor_mentees_empty": [.ru: "Менти пока не закреплены", .kz: "Әзірге менти бекітілмеген", .en: "No mentees assigned yet"],
        "mentor_randomizer_footer": [.ru: "Менти достаются случайно, только по вашей специальности", .kz: "Менти кездейсоқ түрде, тек сіздің мамандығыңыз бойынша беріледі", .en: "Mentees are assigned randomly, only from your specialty"],

        // Advisor screen
        "advisor_main_title": [.ru: "Главная", .kz: "Басты бет", .en: "Home"],
        "advisor_mentors_section": [.ru: "Мои менторы", .kz: "Менің менторларым", .en: "My mentors"],
        "advisor_complaint_action": [.ru: "Жалоба", .kz: "Шағым", .en: "Complaint"],

        // Head/Vice screen
        "head_main_title": [.ru: "Главная", .kz: "Басты бет", .en: "Home"],
        "head_stats_section": [.ru: "Статистика", .kz: "Статистика", .en: "Statistics"],
        "head_stats_mentors": [.ru: "Менторов", .kz: "Менторлар", .en: "Mentors"],
        "head_stats_mentees": [.ru: "Менти", .kz: "Ментилер", .en: "Mentees"],
        "head_stats_complaints": [.ru: "Жалоб на рассмотрении", .kz: "Қаралып жатқан шағымдар", .en: "Pending complaints"],
        "head_management_section": [.ru: "Управление", .kz: "Басқару", .en: "Management"],
        "head_complaints_link": [.ru: "Жалобы", .kz: "Шағымдар", .en: "Complaints"],
        "head_surveys_link": [.ru: "Опросы менти о менторах", .kz: "Ментилердің менторлар туралы сауалнамасы", .en: "Mentee surveys about mentors"],
        "head_documents_link": [.ru: "Документы и хэндбук", .kz: "Құжаттар мен анықтамалық", .en: "Documents and handbook"],

        // Admin screen
        "admin_main_title": [.ru: "Главная", .kz: "Басты бет", .en: "Home"],
        "admin_management_section": [.ru: "Управление", .kz: "Басқару", .en: "Management"],

        // Evaluation form
        "evaluation_title": [.ru: "Оценка ментора", .kz: "Менторды бағалау", .en: "Evaluate mentor"],
        "evaluation_target_section": [.ru: "Кого оцениваем", .kz: "Кімді бағалаймыз", .en: "Evaluating"],
        "evaluation_score_section": [.ru: "Активность (1 — плохо, 5 — отлично)", .kz: "Белсенділік (1 — нашар, 5 — өте жақсы)", .en: "Activity (1 — poor, 5 — excellent)"],
        "evaluation_comment_section": [.ru: "Комментарий", .kz: "Пікір", .en: "Comment"],
        "evaluation_submit_button": [.ru: "Сохранить оценку", .kz: "Бағаны сақтау", .en: "Save evaluation"],

        // Complaint form/list
        "complaint_form_title": [.ru: "Новая жалоба", .kz: "Жаңа шағым", .en: "New complaint"],
        "complaint_about_section": [.ru: "Жалоба на", .kz: "Шағым", .en: "Complaint about"],
        "complaint_description_section": [.ru: "Описание проблемы", .kz: "Мәселенің сипаттамасы", .en: "Problem description"],
        "complaint_submit_button": [.ru: "Отправить жалобу", .kz: "Шағымды жіберу", .en: "Submit complaint"],
        "complaint_list_title": [.ru: "Жалобы", .kz: "Шағымдар", .en: "Complaints"],
        "complaint_status_pending": [.ru: "На рассмотрении", .kz: "Қаралуда", .en: "Pending"],
        "complaint_status_resolved": [.ru: "Решено", .kz: "Шешілді", .en: "Resolved"],
        "complaint_status_dismissed": [.ru: "Отклонено", .kz: "Қабылданбады", .en: "Dismissed"],
        "complaint_resolve_button": [.ru: "Решено", .kz: "Шешілді", .en: "Resolve"],
        "complaint_dismiss_button": [.ru: "Отклонить", .kz: "Қабылданбады", .en: "Dismiss"],

        // Survey
        "survey_title": [.ru: "Опросы", .kz: "Сауалнамалар", .en: "Surveys"],
        "survey_new_section": [.ru: "Новый опрос", .kz: "Жаңа сауалнама", .en: "New survey"],
        "survey_mentor_picker": [.ru: "Ментор", .kz: "Ментор", .en: "Mentor"],
        "survey_mentor_pick_placeholder": [.ru: "Выберите", .kz: "Таңдаңыз", .en: "Choose"],
        "survey_question_placeholder": [.ru: "Вопрос менти", .kz: "Ментиге сұрақ", .en: "Question for mentees"],
        "survey_create_button": [.ru: "Создать опрос", .kz: "Сауалнама құру", .en: "Create survey"],
        "survey_results_section": [.ru: "Результаты", .kz: "Нәтижелер", .en: "Results"],
        "survey_no_responses": [.ru: "Ответов пока нет", .kz: "Жауаптар әлі жоқ", .en: "No responses yet"],

        // Document management
        "document_title": [.ru: "Документы", .kz: "Құжаттар", .en: "Documents"],
        "document_new_section": [.ru: "Новый документ", .kz: "Жаңа құжат", .en: "New document"],
        "document_title_placeholder": [.ru: "Название", .kz: "Атауы", .en: "Title"],
        "document_type_label": [.ru: "Тип", .kz: "Түрі", .en: "Type"],
        "document_type_handbook": [.ru: "Хэндбук", .kz: "Анықтамалық", .en: "Handbook"],
        "document_type_codex": [.ru: "Кодекс", .kz: "Кодекс", .en: "Codex"],
        "document_type_process": [.ru: "Процесс запуска", .kz: "Іске қосу процесі", .en: "Launch process"],
        "document_type_rup": [.ru: "РУП", .kz: "ОЖ", .en: "Curriculum"],
        "document_language_label": [.ru: "Язык", .kz: "Тіл", .en: "Language"],
        "document_specialty_placeholder": [.ru: "Специальность (например 26BDIS)", .kz: "Мамандық (мысалы 26BDIS)", .en: "Specialty (e.g. 26BDIS)"],
        "document_upload_button": [.ru: "Выбрать файл и загрузить", .kz: "Файлды таңдап жүктеу", .en: "Choose file and upload"],
        "document_list_section": [.ru: "Документы", .kz: "Құжаттар", .en: "Documents"],

        // Meeting check-in
        "meeting_title": [.ru: "Отметка встречи", .kz: "Кездесуді белгілеу", .en: "Meeting check-in"],
        "meeting_semester_section": [.ru: "Семестр", .kz: "Семестр", .en: "Semester"],
        "meeting_semester_placeholder": [.ru: "Например: Осень 2026", .kz: "Мысалы: Күз 2026", .en: "e.g. Fall 2026"],
        "meeting_photo_section": [.ru: "Фото со встречи", .kz: "Кездесу суреті", .en: "Meeting photo"],
        "meeting_pick_photo": [.ru: "Выбрать фото", .kz: "Суретті таңдау", .en: "Choose photo"],
        "meeting_submit_button": [.ru: "Отметить встречу", .kz: "Кездесуді белгілеу", .en: "Check in"],
        "meeting_history_section": [.ru: "История", .kz: "Тарих", .en: "History"],
        "meeting_approved": [.ru: "Подтверждено", .kz: "Расталды", .en: "Approved"],
        "meeting_pending": [.ru: "На проверке", .kz: "Тексерілуде", .en: "Pending review"],

        // Check-in review (advisor/head/vice)
        "checkin_review_title": [.ru: "Проверка встреч", .kz: "Кездесулерді тексеру", .en: "Meeting check-ins"],
        "checkin_review_approve_button": [.ru: "Подтвердить", .kz: "Растау", .en: "Approve"],
        "checkin_review_empty": [.ru: "Пока нет отметок о встречах", .kz: "Әзірге кездесу белгілері жоқ", .en: "No check-ins yet"],

        // Mentee detail
        "mentee_detail_title": [.ru: "Менти", .kz: "Менти", .en: "Mentee"],
        "mentee_detail_info_section": [.ru: "Информация", .kz: "Ақпарат", .en: "Info"],
        "mentee_detail_email": [.ru: "Email", .kz: "Email", .en: "Email"],
        "mentee_detail_group": [.ru: "Группа", .kz: "Топ", .en: "Group"],
        "mentee_detail_specialty": [.ru: "Специальность", .kz: "Мамандық", .en: "Specialty"],
        "mentee_detail_assigned_date": [.ru: "Дата закрепления", .kz: "Бекітілген күні", .en: "Assigned on"],
        "mentee_detail_status_section": [.ru: "Статус", .kz: "Мәртебе", .en: "Status"],

        // User management (admin)
        "usermgmt_title": [.ru: "Пользователи", .kz: "Пайдаланушылар", .en: "Users"],
        "usermgmt_link": [.ru: "Пользователи и роли", .kz: "Пайдаланушылар мен рөлдер", .en: "Users and roles"],
        "usermgmt_role_label": [.ru: "Роль", .kz: "Рөл", .en: "Role"],
    ]
}

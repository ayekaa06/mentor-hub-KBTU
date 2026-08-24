//
//  SharedRepositories.swift
//  MentorsHub
//
//  Пока нет бэкенда, все экраны в рамках одного запуска приложения
//  должны работать с одними и теми же mock-репозиториями — иначе
//  жалоба, поданная эдвайзером, не будет видна хеду, и т.д.
//
//  Когда появится реальный API — просто заменяешь правую часть
//  каждой константы на API-реализацию, остальной код (ViewModel'и)
//  не трогаешь вообще, они зависят только от протокола.
//

import Foundation

enum SharedRepositories {
    static let auth: AuthRepository = MockAuthRepository()
    static let mentor: MentorRepository = MockMentorRepository()
    static let complaint: ComplaintRepository = MockComplaintRepository()
    static let survey: SurveyRepository = MockSurveyRepository()
    static let evaluation: EvaluationRepository = MockEvaluationRepository()
    static let meeting: MeetingRepository = MockMeetingRepository()
    static let documents: DocumentManagementRepository = MockDocumentManagementRepository()
}

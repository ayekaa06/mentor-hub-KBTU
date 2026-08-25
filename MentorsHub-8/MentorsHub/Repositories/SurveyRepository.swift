//
//  SurveyRepository.swift
//  MentorsHub
//

import Foundation

protocol SurveyRepository {
    func create(_ survey: Survey)
    func getAll() -> [Survey]
}

final class MockSurveyRepository: SurveyRepository {
    private var surveys: [Survey] = [
        Survey(id: UUID(), aboutMentorId: UUID(), aboutMentorName: "Ержан С.",
               question: "Насколько вы довольны работой ментора в этом семестре?",
               responses: [
                SurveyResponse(id: UUID(), answer: "Очень доволен, всегда на связи"),
                SurveyResponse(id: UUID(), answer: "Редко отвечает в чате")
               ],
               createdDate: Date())
    ]

    func create(_ survey: Survey) {
        surveys.append(survey)
    }

    func getAll() -> [Survey] {
        surveys.sorted { $0.createdDate > $1.createdDate }
    }
}

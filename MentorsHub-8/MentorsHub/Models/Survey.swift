//
//  Survey.swift
//  MentorsHub
//

import Foundation

struct SurveyResponse: Identifiable, Codable, Hashable {
    let id: UUID
    var answer: String
}

struct Survey: Identifiable, Codable, Hashable {
    let id: UUID
    var aboutMentorId: UUID
    var aboutMentorName: String
    var question: String
    var responses: [SurveyResponse]
    var createdDate: Date
}

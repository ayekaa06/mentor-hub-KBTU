//
//  AdvisorRepository.swift
//  MentorsHub
//

import Foundation

protocol AdvisorRepository {
    func getMyMentors() -> [User]
}

struct MockAdvisorRepository: AdvisorRepository {
    func getMyMentors() -> [User] {
        [
            User(id: UUID(), name: "Ержан С.", email: "e@s.kz", role: .mentor, specialty: nil),
            User(id: UUID(), name: "Мадина Р.", email: "m@r.kz", role: .mentor, specialty: nil)
        ]
    }
}

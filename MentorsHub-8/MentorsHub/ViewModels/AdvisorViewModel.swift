//
//  AdvisorViewModel.swift
//  MentorsHub
//

import Foundation

@Observable
class AdvisorViewModel {
    private let repository: AdvisorRepository
    var myMentors: [User] = []

    init(repository: AdvisorRepository = MockAdvisorRepository()) {
        self.repository = repository
        myMentors = repository.getMyMentors()
    }
}

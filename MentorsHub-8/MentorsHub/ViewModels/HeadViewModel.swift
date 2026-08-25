//
//  HeadViewModel.swift
//  MentorsHub
//

import Foundation

@Observable
class HeadViewModel {
    private let repository: HeadRepository
    var totalMentors: Int = 0
    var totalMentees: Int = 0
    var pendingComplaints: Int = 0

    init(repository: HeadRepository = MockHeadRepository()) {
        self.repository = repository
        load()
    }

    func load() {
        let stats = repository.getStats()
        totalMentors = stats.mentors
        totalMentees = stats.mentees
        pendingComplaints = stats.pendingComplaints
    }
}

//
//  MentorViewModel.swift
//  MentorsHub
//

import Foundation

@Observable
class MentorViewModel {
    private let repository: MentorRepository
    private let mentorSpecialty: String

    var myMentees: [Mentee] = []
    var hasUnassignedMentees: Bool = false

    init(mentorSpecialty: String, repository: MentorRepository = SharedRepositories.mentor) {
        self.mentorSpecialty = mentorSpecialty
        self.repository = repository
        loadContent()
    }

    func loadContent() {
        myMentees = repository.getMyMentees()
        hasUnassignedMentees = repository.hasUnassignedMentees(forSpecialty: mentorSpecialty)
    }

    func assignRandomMentee() {
        guard let new = repository.assignRandomMentee(forSpecialty: mentorSpecialty) else { return }
        myMentees.append(new)
        hasUnassignedMentees = repository.hasUnassignedMentees(forSpecialty: mentorSpecialty)
    }

    func updateStatus(for mentee: Mentee, to status: MenteeStatus) {
        repository.updateStatus(menteeId: mentee.id, status: status)
        loadContent()
    }
}

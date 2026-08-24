//
//  ComplaintViewModel.swift
//  MentorsHub
//

import Foundation

@Observable
class ComplaintViewModel {
    private let repository: ComplaintRepository
    var complaints: [Complaint] = []
    var description: String = ""

    init(repository: ComplaintRepository = SharedRepositories.complaint) {
        self.repository = repository
        load()
    }

    func load() {
        complaints = repository.getAll()
    }

    func submit(from: User, about: User) {
        guard !description.isEmpty else { return }
        let complaint = Complaint(id: UUID(), fromUserId: from.id, fromUserName: from.name,
                                   aboutUserId: about.id, aboutUserName: about.name,
                                   description: description, status: .pending, date: Date())
        repository.submit(complaint)
        description = ""
        load()
    }

    func resolve(_ complaint: Complaint, status: ComplaintStatus) {
        repository.resolve(id: complaint.id, status: status)
        load()
    }
}

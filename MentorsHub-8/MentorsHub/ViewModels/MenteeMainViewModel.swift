//
//  MenteeMainViewModel.swift
//  MentorsHub
//

import Foundation

@Observable
class MenteeMainViewModel {
    private let repository: MenteeContentRepository
    private let specialty: String

    var documents: [Document] = []
    var events: [AcademicEvent] = []
    var selectedHandbookLanguage: Language = .ru

    init(specialty: String, repository: MenteeContentRepository = MockMenteeContentRepository()) {
        self.specialty = specialty
        self.repository = repository
        loadContent()
    }

    func loadContent() {
        documents = repository.getDocuments(for: specialty)
        events = repository.getCalendarEvents()
    }

    var otherDocuments: [Document] {
        documents.filter { $0.type != .handbook }
    }

    var currentHandbook: Document? {
        documents.first { $0.type == .handbook && $0.language == selectedHandbookLanguage }
    }
}

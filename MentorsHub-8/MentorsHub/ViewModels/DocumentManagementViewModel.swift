//
//  DocumentManagementViewModel.swift
//  MentorsHub
//

import Foundation

@Observable
class DocumentManagementViewModel {
    private let repository: DocumentManagementRepository
    var documents: [Document] = []

    var newTitle: String = ""
    var newType: DocumentType = .handbook
    var newLanguage: Language = .ru
    var newSpecialty: String = ""

    init(repository: DocumentManagementRepository = SharedRepositories.documents) {
        self.repository = repository
        load()
    }

    func load() {
        documents = repository.getAllDocuments()
    }

    /// Вызывается после того как пользователь выбрал файл через .fileImporter
    func upload(pickedURL: URL) {
        guard !newTitle.isEmpty else { return }
        let doc = Document(
            id: UUID(),
            title: newTitle,
            type: newType,
            fileURL: pickedURL,
            specialty: newType == .rup ? newSpecialty : nil,
            language: newType == .handbook ? newLanguage : nil
        )
        repository.upload(doc)
        newTitle = ""
        newSpecialty = ""
        load()
    }
}

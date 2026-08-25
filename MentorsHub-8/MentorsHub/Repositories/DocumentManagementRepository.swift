//
//  DocumentManagementRepository.swift
//  MentorsHub
//
//  Отдельно от MenteeContentRepository: тот отдаёт документы read-only
//  (для менти/ментора), этот умеет ещё и добавлять новые (для
//  админа/хеда/вице, у кого canManageDocuments == true).
//

import Foundation

protocol DocumentManagementRepository {
    func getAllDocuments() -> [Document]
    func upload(_ document: Document)
}

final class MockDocumentManagementRepository: DocumentManagementRepository {
    private var documents: [Document] = [
        Document(id: UUID(), title: "Кодекс поведения ментора", type: .codex,
                 fileURL: Bundle.main.url(forResource: "Kodeks_Mentora", withExtension: "docx")
                    ?? URL(string: "https://example.com/codex.pdf")!,
                 specialty: nil, language: nil),
        Document(id: UUID(), title: "Процесс запуска менторской группы", type: .process,
                 fileURL: Bundle.main.url(forResource: "Process_Zapuska", withExtension: "docx")
                    ?? URL(string: "https://example.com/process.pdf")!,
                 specialty: nil, language: nil),
        Document(id: UUID(), title: "Хэндбук", type: .handbook,
                 fileURL: URL(string: "https://example.com/handbook_ru.pdf")!, specialty: nil, language: .ru),
        Document(id: UUID(), title: "Хэндбук", type: .handbook,
                 fileURL: URL(string: "https://example.com/handbook_kz.pdf")!, specialty: nil, language: .kz),
        Document(id: UUID(), title: "Хэндбук", type: .handbook,
                 fileURL: URL(string: "https://example.com/handbook_en.pdf")!, specialty: nil, language: .en)
    ]

    func getAllDocuments() -> [Document] { documents }

    func upload(_ document: Document) {
        documents.append(document)
    }
}

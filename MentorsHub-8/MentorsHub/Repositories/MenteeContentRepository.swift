//
//  MenteeContentRepository.swift
//  MentorsHub
//

import Foundation

protocol MenteeContentRepository {
    func getDocuments(for specialty: String) -> [Document]
    func getCalendarEvents() -> [AcademicEvent]
}

struct MockMenteeContentRepository: MenteeContentRepository {
    func getDocuments(for specialty: String) -> [Document] {
        [
            Document(id: UUID(), title: "РУП \(specialty)", type: .rup,
                     fileURL: URL(string: "https://example.com/rup.pdf")!,
                     specialty: specialty, language: nil),
            Document(id: UUID(), title: "Хэндбук", type: .handbook,
                     fileURL: URL(string: "https://example.com/handbook_ru.pdf")!,
                     specialty: nil, language: .ru),
            Document(id: UUID(), title: "Хэндбук", type: .handbook,
                     fileURL: URL(string: "https://example.com/handbook_kz.pdf")!,
                     specialty: nil, language: .kz),
            Document(id: UUID(), title: "Хэндбук", type: .handbook,
                     fileURL: URL(string: "https://example.com/handbook_en.pdf")!,
                     specialty: nil, language: .en),
            Document(id: UUID(), title: "Кодекс поведения ментора", type: .codex,
                     fileURL: Bundle.main.url(forResource: "Kodeks_Mentora", withExtension: "docx")
                        ?? URL(string: "https://example.com/codex.pdf")!,
                     specialty: nil, language: nil),
            Document(id: UUID(), title: "Процесс запуска менторской группы", type: .process,
                     fileURL: Bundle.main.url(forResource: "Process_Zapuska", withExtension: "docx")
                        ?? URL(string: "https://example.com/process.pdf")!,
                     specialty: nil, language: nil)
        ]
    }

    func getCalendarEvents() -> [AcademicEvent] {
        [
            AcademicEvent(id: UUID(), title: "Начало сессии", date: Date().addingTimeInterval(86400 * 30)),
            AcademicEvent(id: UUID(), title: "Каникулы", date: Date().addingTimeInterval(86400 * 60))
        ]
    }
}

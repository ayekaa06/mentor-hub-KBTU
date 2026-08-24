//
//  Document.swift
//  MentorsHub
//
//  Created by Abylai  on 29.07.2026.
//

import Foundation

enum DocumentType: String, Codable, CaseIterable {
    case handbook
    case codex
    case process
    case rup
}

struct Document: Identifiable, Codable, Hashable {

    let id: UUID
    var title: String
    var type: DocumentType
    var fileURL: URL
    var specialty: String?
    var language: Language?
}

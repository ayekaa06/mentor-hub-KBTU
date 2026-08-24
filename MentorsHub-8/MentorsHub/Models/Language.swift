//
//  Language.swift
//  MentorsHub
//
//  Created by Abylai  on 29.07.2026.
//

enum Language: String, Codable, CaseIterable, Identifiable {
    case ru, kz, en
    
    var id: String { rawValue }
    
    var displayName: String {
        switch self {
        case .ru: "Русский"
        case .kz: "Қазақша"
        case .en: "English"
        }
    }
}

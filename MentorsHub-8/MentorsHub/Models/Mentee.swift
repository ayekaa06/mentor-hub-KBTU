//
//  Mentee.swift
//  MentorsHub
//
//  Created by Abylai  on 05.08.2026.
//

import SwiftUI

enum MenteeStatus: String, Codable, CaseIterable {
    case active, inactive
    
    var displayName: String {
        switch self {
        case .active: "Активный"
        case .inactive: "Неактивный"
        }
    }
}

struct Mentee: Identifiable, Codable, Hashable {
    let id: UUID
    var user: User
    var specialty: String
    var group: String
    var assignedTo: UUID?
    var status: MenteeStatus
    var dateAssigned: Date?
}

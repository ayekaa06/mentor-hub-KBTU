//
//  User.swift
//  MentorsHub
//
//  Created by Abylai  on 05.08.2026.
//
import SwiftUI

enum Role: String, Codable, CaseIterable, Identifiable {
    case mentee, mentor, advisor, vice, head, admin

    var id: String { rawValue }

    var localizationKey: String {
        switch self {
        case .mentee: "role_mentee"
        case .mentor: "role_mentor"
        case .advisor: "role_advisor"
        case .vice: "role_vice"
        case .head: "role_head"
        case .admin: "role_admin"
        }
    }
}

struct User: Identifiable, Codable, Hashable {
    let id: UUID
    var name: String
    var email: String
    var role: Role
    var specialty: String?
}

extension User {
    var canHaveMentees: Bool { [.mentor, .advisor, .vice, .head].contains(role) }
    var canEvaluateMentors: Bool { [.advisor, .vice, .head].contains(role) }
    var canManageDocuments: Bool { [.admin, .vice, .head].contains(role) }
    var canManageFAQ: Bool { [.admin, .vice, .head].contains(role) }
    var canResolveComplaints: Bool { [.vice, .head].contains(role) }
    var canRunSurveys: Bool { [.vice, .head].contains(role) }
}

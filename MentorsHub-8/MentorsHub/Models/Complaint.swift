//
//  Complaint.swift
//  MentorsHub
//

import Foundation

enum ComplaintStatus: String, Codable, CaseIterable {
    case pending, resolved, dismissed

    var displayName: String {
        switch self {
        case .pending: "На рассмотрении"
        case .resolved: "Решено"
        case .dismissed: "Отклонено"
        }
    }
}

struct Complaint: Identifiable, Codable, Hashable {
    let id: UUID
    var fromUserId: UUID
    var fromUserName: String
    var aboutUserId: UUID
    var aboutUserName: String
    var description: String
    var status: ComplaintStatus
    var date: Date
}

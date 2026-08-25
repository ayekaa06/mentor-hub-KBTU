//
//  MeetingCheckIn.swift
//  MentorsHub
//

import Foundation

struct MeetingCheckIn: Identifiable, Codable, Hashable {
    let id: UUID
    var mentorId: UUID
    var mentorName: String
    var semester: String
    var date: Date
    var approved: Bool
    var photoData: Data?
}

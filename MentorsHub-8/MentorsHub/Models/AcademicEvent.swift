//
//  AcademicEvent.swift
//  MentorsHub
//
//  Renamed from Calendar.swift — "Calendar" clashed with Foundation.Calendar
//

import Foundation

struct AcademicEvent: Identifiable, Codable, Hashable {
    let id: UUID
    var title: String
    var date: Date
}

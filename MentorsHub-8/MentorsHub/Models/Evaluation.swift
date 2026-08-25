//
//  Evaluation.swift
//  MentorsHub
//
//  Created by Abylai  on 05.08.2026.
//

import SwiftUI

struct Evaluation: Identifiable, Codable, Hashable {
    let id: UUID
    var evaluatorId: UUID
    var mentorId: UUID
    var month: Date
    var activityScore: Int
    var comment: String
}

//
//  MeetingRepository.swift
//  MentorsHub
//

import Foundation

protocol MeetingRepository {
    func checkIn(_ checkIn: MeetingCheckIn)
    func getCheckIns(forMentor mentorId: UUID) -> [MeetingCheckIn]
    func getAllCheckIns() -> [MeetingCheckIn]
    func approve(id: UUID)
}

final class MockMeetingRepository: MeetingRepository {
    private var checkIns: [MeetingCheckIn] = []

    func checkIn(_ checkIn: MeetingCheckIn) {
        checkIns.append(checkIn)
    }

    func getCheckIns(forMentor mentorId: UUID) -> [MeetingCheckIn] {
        checkIns.filter { $0.mentorId == mentorId }
    }

    func getAllCheckIns() -> [MeetingCheckIn] {
        checkIns.sorted { $0.date > $1.date }
    }

    func approve(id: UUID) {
        if let idx = checkIns.firstIndex(where: { $0.id == id }) {
            checkIns[idx].approved = true
        }
    }
}

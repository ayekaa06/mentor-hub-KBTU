//
//  CheckInReviewViewModel.swift
//  MentorsHub
//

import Foundation

@Observable
class CheckInReviewViewModel {
    private let repository: MeetingRepository
    var checkIns: [MeetingCheckIn] = []

    init(repository: MeetingRepository = SharedRepositories.meeting) {
        self.repository = repository
        load()
    }

    func load() {
        checkIns = repository.getAllCheckIns()
    }

    func approve(_ checkIn: MeetingCheckIn) {
        repository.approve(id: checkIn.id)
        load()
    }
}

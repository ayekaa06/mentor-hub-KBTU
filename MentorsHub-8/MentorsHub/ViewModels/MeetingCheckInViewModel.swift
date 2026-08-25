//
//  MeetingCheckInViewModel.swift
//  MentorsHub
//

import Foundation
import SwiftUI
import PhotosUI

@Observable
class MeetingCheckInViewModel {
    private let repository: MeetingRepository
    let mentor: User

    var semester: String = "Осень 2026"
    var selectedPhotoItem: PhotosPickerItem? {
        didSet { loadPhoto() }
    }
    var photoData: Data?
    var checkIns: [MeetingCheckIn] = []
    var didSubmit = false

    init(mentor: User, repository: MeetingRepository = SharedRepositories.meeting) {
        self.mentor = mentor
        self.repository = repository
        load()
    }

    func load() {
        checkIns = repository.getCheckIns(forMentor: mentor.id)
    }

    private func loadPhoto() {
        guard let item = selectedPhotoItem else { return }
        Task {
            if let data = try? await item.loadTransferable(type: Data.self) {
                await MainActor.run { self.photoData = data }
            }
        }
    }

    func submit() {
        guard photoData != nil else { return }
        let checkIn = MeetingCheckIn(id: UUID(), mentorId: mentor.id, mentorName: mentor.name,
                                      semester: semester, date: Date(), approved: false, photoData: photoData)
        repository.checkIn(checkIn)
        didSubmit = true
        load()
    }
}

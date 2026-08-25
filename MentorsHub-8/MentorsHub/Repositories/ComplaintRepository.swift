//
//  ComplaintRepository.swift
//  MentorsHub
//

import Foundation

protocol ComplaintRepository {
    func submit(_ complaint: Complaint)
    func getAll() -> [Complaint]
    func resolve(id: UUID, status: ComplaintStatus)
}

final class MockComplaintRepository: ComplaintRepository {
    private var complaints: [Complaint] = [
        Complaint(id: UUID(), fromUserId: UUID(), fromUserName: "Нурлан С. (эдвайзер)",
                  aboutUserId: UUID(), aboutUserName: "Ержан С. (ментор)",
                  description: "Не выходит на связь с менти уже 3 недели.",
                  status: .pending, date: Date())
    ]

    func submit(_ complaint: Complaint) {
        complaints.append(complaint)
    }

    func getAll() -> [Complaint] {
        complaints.sorted { $0.date > $1.date }
    }

    func resolve(id: UUID, status: ComplaintStatus) {
        if let idx = complaints.firstIndex(where: { $0.id == id }) {
            complaints[idx].status = status
        }
    }
}

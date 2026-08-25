//
//  HeadRepository.swift
//  MentorsHub
//
//  Раньше цифры были захардкожены. Теперь считаются из тех же
//  shared-репозиториев, которыми пользуются остальные роли — так что
//  цифры реально меняются, когда ментор берёт менти или кто-то подаёт жалобу.
//

import Foundation

protocol HeadRepository {
    func getStats() -> (mentors: Int, mentees: Int, pendingComplaints: Int)
}

struct MockHeadRepository: HeadRepository {
    func getStats() -> (mentors: Int, mentees: Int, pendingComplaints: Int) {
        let mentorsCount = SharedRepositories.auth.getAllUsers().filter { $0.role == .mentor }.count
        let menteesCount = SharedRepositories.mentor.getAllMenteesCount()
        let pendingCount = SharedRepositories.complaint.getAll().filter { $0.status == .pending }.count
        return (mentors: mentorsCount, mentees: menteesCount, pendingComplaints: pendingCount)
    }
}

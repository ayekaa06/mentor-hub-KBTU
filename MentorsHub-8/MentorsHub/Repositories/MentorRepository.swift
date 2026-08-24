//
//  MentorRepository.swift
//  MentorsHub
//

import Foundation

protocol MentorRepository {
    func getMyMentees() -> [Mentee]
    /// Есть ли свободные менти именно по специальности этого ментора
    func hasUnassignedMentees(forSpecialty specialty: String) -> Bool
    /// Рандомайзер отдаёт менти только из той же специальности, что у ментора
    func assignRandomMentee(forSpecialty specialty: String) -> Mentee?
    func updateStatus(menteeId: UUID, status: MenteeStatus)
    /// Все менти в системе (закреплённые + свободные) — для агрегированной статистики хеда/вице
    func getAllMenteesCount() -> Int
}

final class MockMentorRepository: MentorRepository {
    private var assigned: [Mentee] = [
        Mentee(id: UUID(),
               user: User(id: UUID(), name: "Алина К.", email: "a@k.kz", role: .mentee, specialty: "26BDIS"),
               specialty: "26BDIS", group: "IS-2601", assignedTo: nil, status: .active, dateAssigned: Date())
    ]

    private var unassigned: [Mentee] = [
        Mentee(id: UUID(),
               user: User(id: UUID(), name: "Данияр Т.", email: "d@t.kz", role: .mentee, specialty: "26BDCS"),
               specialty: "26BDCS", group: "CS-2602", assignedTo: nil, status: .active, dateAssigned: nil),
        Mentee(id: UUID(),
               user: User(id: UUID(), name: "Сая М.", email: "s@m.kz", role: .mentee, specialty: "26BDIS"),
               specialty: "26BDIS", group: "IS-2601", assignedTo: nil, status: .active, dateAssigned: nil),
        Mentee(id: UUID(),
               user: User(id: UUID(), name: "Тимур Б.", email: "t@b.kz", role: .mentee, specialty: "26BDCS"),
               specialty: "26BDCS", group: "CS-2602", assignedTo: nil, status: .active, dateAssigned: nil)
    ]

    func getMyMentees() -> [Mentee] { assigned }

    func hasUnassignedMentees(forSpecialty specialty: String) -> Bool {
        unassigned.contains { $0.specialty == specialty }
    }

    func assignRandomMentee(forSpecialty specialty: String) -> Mentee? {
        let candidateIndices = unassigned.indices.filter { unassigned[$0].specialty == specialty }
        guard let idx = candidateIndices.randomElement() else { return nil }
        var mentee = unassigned.remove(at: idx)
        mentee.dateAssigned = Date()
        assigned.append(mentee)
        return mentee
    }

    func updateStatus(menteeId: UUID, status: MenteeStatus) {
        if let idx = assigned.firstIndex(where: { $0.id == menteeId }) {
            assigned[idx].status = status
        }
    }

    func getAllMenteesCount() -> Int {
        assigned.count + unassigned.count
    }
}

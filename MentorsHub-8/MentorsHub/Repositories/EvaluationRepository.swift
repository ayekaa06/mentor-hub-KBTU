//
//  EvaluationRepository.swift
//  MentorsHub
//

import Foundation

protocol EvaluationRepository {
    func submit(_ evaluation: Evaluation)
    func getEvaluations(forMentor mentorId: UUID) -> [Evaluation]
}

final class MockEvaluationRepository: EvaluationRepository {
    private var evaluations: [Evaluation] = []

    func submit(_ evaluation: Evaluation) {
        evaluations.append(evaluation)
    }

    func getEvaluations(forMentor mentorId: UUID) -> [Evaluation] {
        evaluations.filter { $0.mentorId == mentorId }
    }
}

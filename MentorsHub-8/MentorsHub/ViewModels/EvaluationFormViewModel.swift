//
//  EvaluationFormViewModel.swift
//  MentorsHub
//

import Foundation

@Observable
class EvaluationFormViewModel {
    private let repository: EvaluationRepository
    let mentor: User
    let evaluatorId: UUID

    var activityScore: Int = 3
    var comment: String = ""
    var didSubmit = false

    init(mentor: User, evaluatorId: UUID, repository: EvaluationRepository = SharedRepositories.evaluation) {
        self.mentor = mentor
        self.evaluatorId = evaluatorId
        self.repository = repository
    }

    func submit() {
        let evaluation = Evaluation(id: UUID(), evaluatorId: evaluatorId, mentorId: mentor.id,
                                     month: Date(), activityScore: activityScore, comment: comment)
        repository.submit(evaluation)
        didSubmit = true
    }
}

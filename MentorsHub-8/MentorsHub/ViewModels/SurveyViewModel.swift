//
//  SurveyViewModel.swift
//  MentorsHub
//

import Foundation

@Observable
class SurveyViewModel {
    private let repository: SurveyRepository
    var surveys: [Survey] = []

    var newQuestion: String = ""
    var newAboutMentor: User?

    init(repository: SurveyRepository = SharedRepositories.survey) {
        self.repository = repository
        load()
    }

    func load() {
        surveys = repository.getAll()
    }

    func createSurvey(mentors: [User]) {
        guard let mentor = newAboutMentor, !newQuestion.isEmpty else { return }
        let survey = Survey(id: UUID(), aboutMentorId: mentor.id, aboutMentorName: mentor.name,
                             question: newQuestion, responses: [], createdDate: Date())
        repository.create(survey)
        newQuestion = ""
        newAboutMentor = nil
        load()
    }
}

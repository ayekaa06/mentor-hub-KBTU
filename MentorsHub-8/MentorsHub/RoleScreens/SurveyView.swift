//
//  SurveyView.swift
//  MentorsHub
//

import SwiftUI

struct SurveyView: View {
    @State private var viewModel = SurveyViewModel()
    let mentors: [User]
    private let localization = LocalizationManager.shared
    private let accentColor = Color(red: 0.0, green: 0.20, blue: 0.44)

    var body: some View {
        List {
            Section(localization.text("survey_new_section")) {
                Picker(localization.text("survey_mentor_picker"), selection: $viewModel.newAboutMentor) {
                    Text(localization.text("survey_mentor_pick_placeholder")).tag(User?.none)
                    ForEach(mentors) { mentor in
                        Text(mentor.name).tag(User?.some(mentor))
                    }
                }
                TextField(localization.text("survey_question_placeholder"), text: $viewModel.newQuestion)
                Button(localization.text("survey_create_button")) {
                    viewModel.createSurvey(mentors: mentors)
                }
                .disabled(viewModel.newQuestion.isEmpty || viewModel.newAboutMentor == nil)
            }

            Section(localization.text("survey_results_section")) {
                ForEach(viewModel.surveys) { survey in
                    VStack(alignment: .leading, spacing: 4) {
                        Text(survey.aboutMentorName).font(.headline)
                        Text(survey.question).font(.subheadline).foregroundStyle(.secondary)
                        if survey.responses.isEmpty {
                            Text(localization.text("survey_no_responses")).font(.caption).foregroundStyle(.secondary)
                        } else {
                            ForEach(survey.responses) { response in
                                Text("• \(response.answer)").font(.caption)
                            }
                        }
                    }
                    .padding(.vertical, 4)
                }
            }
        }
        .navigationTitle(localization.text("survey_title"))
        .tint(accentColor)
        .onAppear { viewModel.load() }
    }
}

#Preview {
    NavigationStack {
        SurveyView(mentors: [User(id: UUID(), name: "Ержан С.", email: "e@s.kz", role: .mentor, specialty: nil)])
    }
}

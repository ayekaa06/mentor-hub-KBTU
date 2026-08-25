//
//  EvaluationFormView.swift
//  MentorsHub
//

import SwiftUI

struct EvaluationFormView: View {
    @Environment(\.dismiss) private var dismiss
    @State private var viewModel: EvaluationFormViewModel
    private let localization = LocalizationManager.shared
    private let accentColor = Color(red: 0.0, green: 0.20, blue: 0.44)

    init(mentor: User, evaluatorId: UUID) {
        _viewModel = State(initialValue: EvaluationFormViewModel(mentor: mentor, evaluatorId: evaluatorId))
    }

    var body: some View {
        Form {
            Section(localization.text("evaluation_target_section")) {
                Text(viewModel.mentor.name)
                    .foregroundStyle(.secondary)
            }

            Section(localization.text("evaluation_score_section")) {
                Picker("", selection: $viewModel.activityScore) {
                    ForEach(1...5, id: \.self) { score in
                        Text("\(score)").tag(score)
                    }
                }
                .pickerStyle(.segmented)
            }

            Section(localization.text("evaluation_comment_section")) {
                TextEditor(text: $viewModel.comment)
                    .frame(minHeight: 100)
            }

            Section {
                Button(localization.text("evaluation_submit_button")) {
                    viewModel.submit()
                    dismiss()
                }
                .disabled(viewModel.comment.isEmpty)
            }
        }
        .navigationTitle(localization.text("evaluation_title"))
        .tint(accentColor)
    }
}

#Preview {
    NavigationStack {
        EvaluationFormView(
            mentor: User(id: UUID(), name: "Ержан С.", email: "e@s.kz", role: .mentor, specialty: nil),
            evaluatorId: UUID()
        )
    }
}

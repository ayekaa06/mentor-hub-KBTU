//
//  AdvisorScreen.swift
//  MentorsHub
//

import SwiftUI

struct AdvisorScreen: View {
    let user: User
    let onLogout: () -> Void
    @State private var viewModel = AdvisorViewModel()
    private let localization = LocalizationManager.shared
    private let accentColor = Color(red: 0.0, green: 0.20, blue: 0.44)

    var body: some View {
        List {
            Section(localization.text("advisor_mentors_section")) {
                ForEach(viewModel.myMentors) { mentor in
                    NavigationLink(mentor.name) {
                        EvaluationFormView(mentor: mentor, evaluatorId: user.id)
                    }
                    .swipeActions {
                        NavigationLink {
                            ComplaintFormView(fromUser: user, aboutUser: mentor)
                        } label: {
                            Label(localization.text("advisor_complaint_action"), systemImage: "exclamationmark.bubble")
                        }
                        .tint(.red)
                    }
                }
            }

            Section {
                NavigationLink(localization.text("checkin_review_title")) {
                    CheckInReviewView()
                }
            }
        }
        .navigationTitle(localization.text("advisor_main_title"))
        .toolbar {
            ToolbarItem(placement: .navigationBarTrailing) {
                NavigationLink {
                    ProfileScreen(user: user, onLogout: onLogout)
                } label: {
                    Image(systemName: "person.circle")
                }
            }
        }
        .tint(accentColor)
    }
}

#Preview {
    NavigationStack {
        AdvisorScreen(user: User(id: UUID(), name: "Тест", email: "test@test.kz", role: .advisor, specialty: nil), onLogout: {})
    }
}

//
//  HeadScreen.swift
//  MentorsHub
//
//  Общий экран для ролей .vice и .head — права у них почти одинаковые.
//

import SwiftUI

struct HeadScreen: View {
    let user: User
    let onLogout: () -> Void
    @State private var viewModel = HeadViewModel()
    @State private var advisorViewModel = AdvisorViewModel() // ради списка менторов для опросов
    private let localization = LocalizationManager.shared
    private let accentColor = Color(red: 0.0, green: 0.20, blue: 0.44)

    var body: some View {
        List {
            Section(localization.text("head_stats_section")) {
                HStack { Text(localization.text("head_stats_mentors")); Spacer(); Text("\(viewModel.totalMentors)") }
                HStack { Text(localization.text("head_stats_mentees")); Spacer(); Text("\(viewModel.totalMentees)") }
                HStack { Text(localization.text("head_stats_complaints")); Spacer(); Text("\(viewModel.pendingComplaints)") }
            }

            Section(localization.text("head_management_section")) {
                NavigationLink(localization.text("head_complaints_link")) {
                    ComplaintListView()
                }
                NavigationLink(localization.text("head_surveys_link")) {
                    SurveyView(mentors: advisorViewModel.myMentors)
                }
                NavigationLink(localization.text("head_documents_link")) {
                    DocumentManagementView()
                }
                NavigationLink(localization.text("checkin_review_title")) {
                    CheckInReviewView()
                }
                NavigationLink(localization.text("usermgmt_link")) {
                    UserManagementView()
                }
            }
        }
        .navigationTitle(localization.text("head_main_title"))
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
        .onAppear { viewModel.load() }
    }
}

#Preview {
    NavigationStack {
        HeadScreen(user: User(id: UUID(), name: "Тест", email: "test@test.kz", role: .head, specialty: nil), onLogout: {})
    }
}

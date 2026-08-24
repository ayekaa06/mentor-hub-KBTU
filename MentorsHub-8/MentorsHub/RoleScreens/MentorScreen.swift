//
//  MentorScreen.swift
//  MentorsHub
//

import SwiftUI
import QuickLook

struct MentorScreen: View {
    let user: User
    let onLogout: () -> Void
    @State private var viewModel: MentorViewModel
    @State private var contentViewModel = MenteeMainViewModel(specialty: "") // ради Document-репозитория (кодекс/процесс)
    @State private var previewURL: URL?
    private let localization = LocalizationManager.shared
    private let accentColor = Color(red: 0.0, green: 0.20, blue: 0.44)

    init(user: User, onLogout: @escaping () -> Void) {
        self.user = user
        self.onLogout = onLogout
        _viewModel = State(initialValue: MentorViewModel(mentorSpecialty: user.specialty ?? ""))
    }

    private var mySpecialtyName: String {
        Specialty.all.first(where: { $0.code == user.specialty })?.displayName ?? (user.specialty ?? "—")
    }

    var body: some View {
        List {
            Section {
                HStack {
                    Text(localization.text("mentor_specialty_label"))
                    Spacer()
                    Text(mySpecialtyName)
                        .foregroundStyle(.secondary)
                }
            }

            Section(localization.text("mentor_mentees_section")) {
                if viewModel.myMentees.isEmpty {
                    Text(localization.text("mentor_mentees_empty"))
                        .foregroundStyle(.secondary)
                }
                ForEach(viewModel.myMentees) { mentee in
                    NavigationLink {
                        MenteeDetailView(mentee: mentee, viewModel: viewModel)
                    } label: {
                        HStack {
                            VStack(alignment: .leading) {
                                Text(mentee.user.name)
                                Text(mentee.group)
                                    .font(.caption)
                                    .foregroundStyle(.secondary)
                            }
                            Spacer()
                            Circle()
                                .fill(mentee.status == .active ? .green : .gray)
                                .frame(width: 8, height: 8)
                        }
                    }
                }
            }

            if viewModel.hasUnassignedMentees {
                Section {
                    Button(localization.text("mentor_randomizer_button")) {
                        viewModel.assignRandomMentee()
                    }
                } footer: {
                    Text(localization.text("mentor_randomizer_footer"))
                }
            }

            Section(localization.text("mentor_meetings_section")) {
                NavigationLink(localization.text("mentor_checkin_link")) {
                    MeetingCheckInView(mentor: user)
                }
            }

            Section(localization.text("mentor_rules_section")) {
                ForEach(contentViewModel.documents.filter { $0.type == .codex || $0.type == .process }) { doc in
                    Button {
                        previewURL = doc.fileURL
                    } label: {
                        Label(doc.title, systemImage: "doc.text")
                    }
                }
            }
        }
        .navigationTitle(localization.text("mentor_main_title"))
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
        .quickLookPreview($previewURL)
    }
}

#Preview {
    NavigationStack {
        MentorScreen(user: User(id: UUID(), name: "Тест", email: "test@test.kz", role: .mentor, specialty: "26BDIS"), onLogout: {})
    }
}

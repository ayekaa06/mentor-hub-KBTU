//
//  MenteeScreen.swift
//  MentorsHub
//
//  Created by Abylai  on 27.07.2026.
//

import SwiftUI
import QuickLook

struct MenteeScreen: View {
    let user: User
    let onLogout: () -> Void
    @State private var viewModel: MenteeMainViewModel
    @State private var previewURL: URL?
    private let localization = LocalizationManager.shared
    private let accentColor = Color(red: 0.0, green: 0.20, blue: 0.44)

    init(user: User, onLogout: @escaping () -> Void) {
        self.user = user
        self.onLogout = onLogout
        _viewModel = State(initialValue: MenteeMainViewModel(specialty: user.specialty ?? ""))
    }

    private var mySpecialtyName: String {
        Specialty.all.first(where: { $0.code == user.specialty })?.displayName ?? (user.specialty ?? "—")
    }

    private var rupDocuments: [Document] {
        viewModel.otherDocuments.filter { $0.type == .rup }
    }

    var body: some View {
        List {
            Section {
                HStack {
                    Text(localization.text("mentee_specialty_label"))
                    Spacer()
                    Text(mySpecialtyName)
                        .foregroundStyle(.secondary)
                }
            }

            Section(localization.text("mentee_rup_section")) {
                if rupDocuments.isEmpty {
                    Text(localization.text("mentee_rup_empty"))
                        .foregroundStyle(.secondary)
                }
                ForEach(rupDocuments) { doc in
                    Button {
                        previewURL = doc.fileURL
                    } label: {
                        Label(doc.title, systemImage: "doc.text")
                    }
                }
            }

            Section(localization.text("mentee_handbook_section")) {
                Picker(localization.text("mentee_handbook_language"), selection: $viewModel.selectedHandbookLanguage) {
                    ForEach(Language.allCases) { lang in
                        Text(lang.displayName).tag(lang)
                    }
                }
                .pickerStyle(.segmented)

                if let handbook = viewModel.currentHandbook {
                    Button {
                        previewURL = handbook.fileURL
                    } label: {
                        Label(localization.text("mentee_open_handbook"), systemImage: "questionmark.circle")
                    }
                } else {
                    Text(localization.text("mentee_handbook_unavailable"))
                        .foregroundStyle(.secondary)
                }
            }

            Section(localization.text("mentee_calendar_section")) {
                ForEach(viewModel.events) { event in
                    HStack {
                        Text(event.title)
                        Spacer()
                        Text(event.date.formatted(date: .abbreviated, time: .omitted))
                            .foregroundStyle(.secondary)
                    }
                }
            }
        }
        .navigationTitle(localization.text("mentee_main_title"))
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
        MenteeScreen(user: User(id: UUID(), name: "Тест", email: "test@test.kz", role: .mentee, specialty: "26BDIS"), onLogout: {})
    }
}

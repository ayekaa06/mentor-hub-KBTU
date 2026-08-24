//
//  AdminScreen.swift
//  MentorsHub
//

import SwiftUI

struct AdminScreen: View {
    let user: User
    let onLogout: () -> Void
    private let localization = LocalizationManager.shared
    private let accentColor = Color(red: 0.0, green: 0.20, blue: 0.44)

    var body: some View {
        List {
            Section(localization.text("admin_management_section")) {
                NavigationLink(localization.text("head_documents_link")) {
                    DocumentManagementView()
                }
                NavigationLink(localization.text("usermgmt_link")) {
                    UserManagementView()
                }
            }
        }
        .navigationTitle(localization.text("admin_main_title"))
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
        AdminScreen(user: User(id: UUID(), name: "Тест", email: "test@test.kz", role: .admin, specialty: nil), onLogout: {})
    }
}

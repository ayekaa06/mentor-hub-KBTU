//
//  UserManagementView.swift
//  MentorsHub
//

import SwiftUI

struct UserManagementView: View {
    @State private var viewModel = UserManagementViewModel()
    private let localization = LocalizationManager.shared
    private let accentColor = Color(red: 0.0, green: 0.20, blue: 0.44)

    var body: some View {
        List {
            ForEach(viewModel.users) { user in
                VStack(alignment: .leading, spacing: 6) {
                    Text(user.name).font(.headline)
                    Text(user.email).font(.caption).foregroundStyle(.secondary)

                    Picker(localization.text("usermgmt_role_label"), selection: Binding(
                        get: { user.role },
                        set: { viewModel.updateRole(user, to: $0) }
                    )) {
                        ForEach(Role.allCases) { role in
                            Text(localization.text(role.localizationKey)).tag(role)
                        }
                    }
                    .pickerStyle(.menu)
                }
                .padding(.vertical, 4)
            }
        }
        .navigationTitle(localization.text("usermgmt_title"))
        .tint(accentColor)
        .onAppear { viewModel.load() }
    }
}

#Preview {
    NavigationStack {
        UserManagementView()
    }
}

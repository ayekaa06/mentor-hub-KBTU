//
//  ProfileScreen.swift
//  MentorsHub
//
//  Единый экран для всех ролей — язык интерфейса + выход из аккаунта.
//  Открывается через иконку в toolbar на каждом главном экране.
//

import SwiftUI

struct ProfileScreen: View {
    let user: User
    let onLogout: () -> Void
    private let localization = LocalizationManager.shared
    private let accentColor = Color(red: 0.0, green: 0.20, blue: 0.44)

    var body: some View {
        List {
            Section {
                VStack(alignment: .leading, spacing: 4) {
                    Text(user.name)
                        .font(.headline)
                    Text(localization.text(user.role.localizationKey))
                        .font(.subheadline)
                        .foregroundStyle(.secondary)
                }
                .padding(.vertical, 4)
            }

            Section(localization.text("profile_language_section")) {
                Picker(localization.text("profile_language_section"), selection: Binding(
                    get: { localization.currentLanguage },
                    set: { localization.currentLanguage = $0 }
                )) {
                    ForEach(Language.allCases) { lang in
                        Text(lang.displayName).tag(lang)
                    }
                }
                .pickerStyle(.segmented)
            }

            Section {
                Button(role: .destructive) {
                    onLogout()
                } label: {
                    Text(localization.text("profile_logout_button"))
                        .frame(maxWidth: .infinity)
                }
            }
        }
        .navigationTitle(localization.text("profile_title"))
        .tint(accentColor)
    }
}

#Preview {
    NavigationStack {
        ProfileScreen(
            user: User(id: UUID(), name: "Тест", email: "test@test.kz", role: .mentor, specialty: nil),
            onLogout: {}
        )
    }
}

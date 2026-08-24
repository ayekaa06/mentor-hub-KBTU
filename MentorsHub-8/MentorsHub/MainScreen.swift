//
//  MainScreen.swift
//  MentorsHub
//
//  Created by Abylai  on 24.07.2026.
//
//  Роутер по роли — единственная задача: получить User и показать
//  правильный экран, прокинув onLogout дальше.
//

import SwiftUI

struct MainScreen: View {
    let user: User
    let onLogout: () -> Void

    var body: some View {
        Group {
            switch user.role {
            case .mentee:
                MenteeScreen(user: user, onLogout: onLogout)
            case .mentor:
                MentorScreen(user: user, onLogout: onLogout)
            case .advisor:
                AdvisorScreen(user: user, onLogout: onLogout)
            case .vice, .head:
                HeadScreen(user: user, onLogout: onLogout)
            case .admin:
                AdminScreen(user: user, onLogout: onLogout)
            }
        }
        .navigationBarBackButtonHidden(true)
    }
}

#Preview {
    NavigationStack {
        MainScreen(user: User(id: UUID(), name: "Тест", email: "test@test.kz", role: .mentee, specialty: "26BDIS"), onLogout: {})
    }
}

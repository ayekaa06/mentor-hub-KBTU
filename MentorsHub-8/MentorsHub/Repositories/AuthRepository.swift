//
//  AuthRepository.swift
//  MentorsHub
//
//  Пока без бэкенда — MockAuthRepository. Когда API будет готов,
//  создаёшь APIAuthRepository: AuthRepository с реальным сетевым запросом
//  и подменяешь SharedRepositories.auth — экраны трогать не придётся.
//

import Foundation

protocol AuthRepository {
    func login(email: String, password: String) -> User?
    /// Возвращает nil если email уже занят
    func register(name: String, email: String, password: String, role: Role, specialty: String?) -> User?
    func getAllUsers() -> [User]
    func updateRole(userId: UUID, newRole: Role)
}

final class MockAuthRepository: AuthRepository {
    // class, а не struct — чтобы регистрация реально сохраняла нового
    // юзера в общей "базе" на время сессии приложения (см. SharedRepositories.auth)
    private var accounts: [String: (password: String, user: User)] = [
        "mentee@test.kz": ("password123", User(id: UUID(), name: "Айгерим (менти)", email: "mentee@test.kz", role: .mentee, specialty: "26BDIS")),
        "mentor@test.kz": ("password123", User(id: UUID(), name: "Абылай (ментор)", email: "mentor@test.kz", role: .mentor, specialty: "26BDIS")),
        "advisor@test.kz": ("password123", User(id: UUID(), name: "Нурлан (эдвайзер)", email: "advisor@test.kz", role: .advisor, specialty: nil)),
        "vice@test.kz": ("password123", User(id: UUID(), name: "Николь (вице-хэд)", email: "vice@test.kz", role: .vice, specialty: nil)),
        "head@test.kz": ("password123", User(id: UUID(), name: "Аружан (хэд)", email: "head@test.kz", role: .head, specialty: nil)),
        "admin@test.kz": ("password123", User(id: UUID(), name: "Админ", email: "admin@test.kz", role: .admin, specialty: nil))
    ]

    func login(email: String, password: String) -> User? {
        guard let entry = accounts[email.lowercased()], entry.password == password else {
            return nil
        }
        return entry.user
    }

    func register(name: String, email: String, password: String, role: Role, specialty: String?) -> User? {
        let key = email.lowercased()
        guard accounts[key] == nil else { return nil } // уже занят
        let user = User(id: UUID(), name: name, email: key, role: role, specialty: specialty)
        accounts[key] = (password, user)
        return user
    }

    func getAllUsers() -> [User] {
        accounts.values.map { $0.user }.sorted { $0.name < $1.name }
    }

    func updateRole(userId: UUID, newRole: Role) {
        guard let key = accounts.first(where: { $0.value.user.id == userId })?.key else { return }
        accounts[key]?.user.role = newRole
    }
}

//
//  UserManagementViewModel.swift
//  MentorsHub
//

import Foundation

@Observable
class UserManagementViewModel {
    private let repository: AuthRepository
    var users: [User] = []

    init(repository: AuthRepository = SharedRepositories.auth) {
        self.repository = repository
        load()
    }

    func load() {
        users = repository.getAllUsers()
    }

    func updateRole(_ user: User, to role: Role) {
        repository.updateRole(userId: user.id, newRole: role)
        load()
    }
}

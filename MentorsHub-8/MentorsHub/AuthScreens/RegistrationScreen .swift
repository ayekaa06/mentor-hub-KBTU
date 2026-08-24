//
//  RegistrationScreen .swift
//  MentorsHub
//
//  Created by Abylai  on 18.07.2026.
//

import SwiftUI

struct RegistrationScreen_: View {
    @Binding var path: NavigationPath

    enum Field: Hashable {
        case name
        case email
        case password
        case passcheck
        case specialty
    }

    @State private var name: String = ""
    @State private var email: String = ""
    @State private var password: String = ""
    @State private var passcheck: String = ""
    @State private var role: Role = .mentee
    @State private var specialty: String = ""
    @State private var errorMessage: String?

    @FocusState private var focusedField: Field?
    private let authRepository: AuthRepository = SharedRepositories.auth
    private let localization = LocalizationManager.shared

    private var isFormValid: Bool {
        !name.isEmpty && !email.isEmpty && !password.isEmpty && password == passcheck &&
        password.count >= 8 && passcheck.count >= 8 &&
        (![Role.mentee, .mentor].contains(role) || !specialty.isEmpty)
    }

    var body: some View {
        VStack {
            Spacer()

            Text(localization.text("register_title"))
                .font(.system(size: 34, design: .serif))
                .foregroundColor(Color(red: 0.0, green: 0.20, blue: 0.44))

            VStack(spacing: 15) {
                TextField(localization.text("register_name_placeholder"), text: $name)
                    .focused($focusedField, equals: .name)
                    .padding(.horizontal, 5)
                    .font(.system(size: 20))
                    .frame(width: 300, height: 40)
                    .background(Color.gray.opacity(0.2))
                    .cornerRadius(10)
                    .textContentType(.name)

                TextField(localization.text("register_email_placeholder"), text: $email)
                    .focused($focusedField, equals: .email)
                    .padding(.horizontal, 5)
                    .font(.system(size: 20))
                    .frame(width: 300, height: 40)
                    .background(Color.gray.opacity(0.2))
                    .cornerRadius(10)
                    .textContentType(.emailAddress)
                    .textInputAutocapitalization(.never)
                    .autocorrectionDisabled()

                SecureField(localization.text("register_password_placeholder"), text: $password)
                    .focused($focusedField, equals: .password)
                    .padding(.horizontal, 5)
                    .font(.system(size: 20))
                    .frame(width: 300, height: 40)
                    .background(Color.gray.opacity(0.2))
                    .cornerRadius(10)
                    .textContentType(.newPassword)

                SecureField(localization.text("register_passcheck_placeholder"), text: $passcheck)
                    .focused($focusedField, equals: .passcheck)
                    .padding(.horizontal, 5)
                    .font(.system(size: 20))
                    .frame(width: 300, height: 40)
                    .background(Color.gray.opacity(0.2))
                    .cornerRadius(10)
                    .textContentType(.newPassword)

                Picker(localization.text("register_role_label"), selection: $role) {
                    ForEach(Role.allCases) { r in
                        Text(localization.text(r.localizationKey)).tag(r)
                    }
                }
                .pickerStyle(.menu)
                .frame(width: 300)

                if role == .mentee || role == .mentor {
                    Picker(localization.text("register_specialty_placeholder"), selection: $specialty) {
                        Text(localization.text("register_specialty_placeholder")).tag("")
                        ForEach(Specialty.all) { spec in
                            Text(spec.displayName).tag(spec.code)
                        }
                    }
                    .pickerStyle(.menu)
                    .frame(width: 300)
                }
            }
            .padding(.vertical, 20)

            Button {
                guard let user = authRepository.register(
                    name: name, email: email, password: password,
                    role: role, specialty: (role == .mentee || role == .mentor) ? specialty : nil
                ) else {
                    errorMessage = localization.text("register_error_exists")
                    return
                }
                errorMessage = nil
                path.append(user)
            } label: {
                Text(localization.text("register_submit_button"))
            }
            .frame(width: 300, height: 50)
            .foregroundStyle(.white)
            .background(isFormValid ? Color(red: 0.0, green: 0.20, blue: 0.44) : Color.gray)
            .cornerRadius(10)
            .disabled(!isFormValid)

            if let errorMessage {
                Text(errorMessage)
                    .font(.footnote)
                    .foregroundStyle(.red)
            }

            Spacer()

            Button(localization.text("register_have_account")) {
                path.removeLast()
            }
            .padding(.bottom, 10)
        }
        .ignoresSafeArea(.keyboard, edges: .bottom)
        .foregroundStyle(Color(red: 0.0, green: 0.20, blue: 0.44))
        .frame(maxWidth: .infinity, maxHeight: .infinity)
        .contentShape(Rectangle())
        .onTapGesture {
            focusedField = nil
        }
    }
}

#Preview {
    NavigationStack {
        RegistrationScreen_(path: .constant(NavigationPath()))
    }
}

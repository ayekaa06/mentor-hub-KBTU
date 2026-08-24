//
//  AuthScreen.swift
//  MentorsHub
//
//  Created by Abylai  on 17.07.2026.
//

import SwiftUI
import Combine

struct AuthScreen: View {
    @State private var path = NavigationPath()
    @State private var email: String = ""
    @State private var password: String = ""
    @State private var passVisible: Bool = false
    @State private var loginError: String?
    @FocusState private var emailFocused: Bool
    @FocusState private var passFocused: Bool
    private let authRepository: AuthRepository = SharedRepositories.auth
    private let localization = LocalizationManager.shared

    var body: some View {
        NavigationStack(path: $path){
            VStack{
                Spacer()
                Text(localization.text("auth_title"))
                    .font(.system(size: 34, design: .serif))
                    .foregroundColor(Color(red: 0.0, green: 0.20, blue: 0.44))
                VStack{
                    TextField(localization.text("register_email_placeholder"), text: $email)
                        .padding(.horizontal, 5)
                        .font(.system(size: 20))
                        .frame(width: 300, height: 40)
                        .background(Color.gray.opacity(0.2))
                        .cornerRadius(10)
                        .focused($emailFocused)
                        .padding(5)
                        .textContentType(.username)
                        .textInputAutocapitalization(.never)
                        .autocorrectionDisabled()
                    ZStack(alignment: .trailing){
                        Group{
                            if passVisible{
                                TextField(localization.text("auth_password_placeholder"), text: $password)
                            }
                            else{
                                SecureField(localization.text("auth_password_placeholder"), text: $password)
                            }
                        }
                        .padding(.horizontal, 5)
                        .padding(.trailing, 35)
                        .font(.system(size: 20))
                        .frame(width: 300,height: 40)
                        .background(Color.gray.opacity(0.2))
                        .cornerRadius(10)
                        .focused($passFocused)
                        .textContentType(.password)

                        Button{
                            passVisible.toggle()
                        }label: {
                            Image(systemName: passVisible ? "eye.slash" : "eye" )
                                .foregroundStyle(.gray)

                        }.padding(.trailing, 8)
                    }
                }
                .padding(10)
                Button{
                    if let user = authRepository.login(email: email, password: password) {
                        loginError = nil
                        path.append(user)
                    } else {
                        loginError = localization.text("auth_login_error")
                    }
                }label:{
                    Text(localization.text("auth_login_button"))
                    .frame(width: 300, height: 50)
                    .foregroundStyle(Color(.white))
                    .background(email.isEmpty || password.isEmpty ? Color.gray :  Color(red: 0.0, green: 0.20, blue: 0.44))
                    .cornerRadius(10)
                    .contentShape(Rectangle())
                }
                .disabled(email.isEmpty || password.isEmpty)

                if let loginError {
                    Text(loginError)
                        .font(.footnote)
                        .foregroundStyle(.red)
                }

                Spacer()
                Button{
                    path.append("регистрация")
                }label: {
                    Text(localization.text("auth_register_button"))
                    .frame(width: 300, height: 50)
                    .background(Color(red: 0.0, green: 0.20, blue: 0.44))
                    .foregroundStyle(Color(.white))
                    .cornerRadius(10)
                    .padding(10)
                    .contentShape(Rectangle())
                }
                Button(localization.text("auth_forgot_button")){
                    path.append("забыли пароль")
                }
            }
            .foregroundStyle(Color(red: 0.0, green: 0.20, blue: 0.44))
            .ignoresSafeArea(.keyboard, edges: .bottom)
            .contentShape(Rectangle())
            .onTapGesture {
                emailFocused = false
                passFocused = false
            }
            .navigationDestination(for: String.self){ value in
                if value == "регистрация"{
                    RegistrationScreen_(path: $path)
                }else if value == "забыли пароль"{
                    ForgotPassScreen()
                }
            }
            .navigationDestination(for: User.self) { user in
                MainScreen(user: user, onLogout: {
                    email = ""
                    password = ""
                    path = NavigationPath()
                })
            }
        }
    }
}

#Preview {
    AuthScreen()
}

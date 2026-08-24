//
//  ForgotPassScreen.swift
//  MentorsHub
//
//  Created by Abylai  on 22.07.2026.
//

import SwiftUI

enum ForgotPassField: Hashable {
    case email
    case newpass
    case passconfirm
}

struct ForgotPassScreen: View {
    @Environment(\.dismiss) private var dismiss
    @State private var newpass = ""
    @State private var passconfirm = ""
    @State private var email: String = ""
    @State private var passVisible = false
    @State private var confirmVisible = false
    @FocusState private var isFocused: ForgotPassField?
    private let localization = LocalizationManager.shared

    private var isFormValid: Bool {
        !email.isEmpty && !newpass.isEmpty && !passconfirm.isEmpty
    }

    var body: some View {
        VStack {
            Spacer()
            Text(localization.text("forgot_title"))
                .font(.system(size: 30, design: .serif))
                .foregroundColor(Color(red: 0.0, green: 0.20, blue: 0.44))
                .padding(.bottom, 20)
            TextField(localization.text("forgot_email_placeholder"), text: $email)
                .padding(.horizontal, 5)
                .font(.system(size: 20))
                .frame(width: 310, height: 40)
                .background(Color.gray.opacity(0.2))
                .cornerRadius(10)
                .focused($isFocused, equals: .email)
                .textContentType(.username)
                .textInputAutocapitalization(.never)
                .autocorrectionDisabled()

            ZStack(alignment: .trailing) {
                Group {
                    if passVisible {
                        TextField(localization.text("forgot_newpass_placeholder"), text: $newpass)
                    } else {
                        SecureField(localization.text("forgot_newpass_placeholder"), text: $newpass)
                    }
                }
                .focused($isFocused, equals: .newpass)
                .padding(.horizontal, 5)
                .font(.system(size: 20))
                .frame(width: 310, height: 40)
                .background(Color.gray.opacity(0.2))
                .cornerRadius(10)
                .textContentType(.newPassword)

                Button {
                    passVisible.toggle()
                } label: {
                    Image(systemName: passVisible ? "eye.slash" : "eye")
                        .foregroundStyle(.gray)
                        .padding(7)
                }
            }
            .padding(.top, 5)

            ZStack(alignment: .trailing) {
                Group {
                    if confirmVisible {
                        TextField(localization.text("forgot_confirm_placeholder"), text: $passconfirm)
                    } else {
                        SecureField(localization.text("forgot_confirm_placeholder"), text: $passconfirm)
                    }
                }
                .focused($isFocused, equals: .passconfirm)
                .padding(.horizontal, 5)
                .font(.system(size: 20))
                .frame(width: 310, height: 40)
                .background(Color.gray.opacity(0.2))
                .cornerRadius(10)
                .textContentType(.newPassword)

                Button {
                    confirmVisible.toggle()
                } label: {
                    Image(systemName: confirmVisible ? "eye.slash" : "eye")
                        .foregroundStyle(.gray)
                        .padding(7)
                }
            }
            .padding(.top, 5)

            Spacer()

            Button {
                dismiss()
            } label: {
                Text(localization.text("forgot_submit_button"))
                    .frame(width: 300, height: 50)
                    .background(isFormValid ? Color(red: 0.0, green: 0.20, blue: 0.44) : Color(.gray))
                    .foregroundStyle(Color(.white))
                    .cornerRadius(10)
                    .padding(10)
                    .contentShape(Rectangle())
            }
            .disabled(!isFormValid)
            .contentShape(Rectangle())
        }
        .ignoresSafeArea(.keyboard, edges: .bottom)
        .contentShape(Rectangle())
        .onTapGesture {
            isFocused = nil
        }
    }
}

#Preview {
    ForgotPassScreen()
}

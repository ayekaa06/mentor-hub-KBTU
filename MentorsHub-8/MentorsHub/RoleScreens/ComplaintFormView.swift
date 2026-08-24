//
//  ComplaintFormView.swift
//  MentorsHub
//

import SwiftUI

struct ComplaintFormView: View {
    @Environment(\.dismiss) private var dismiss
    @State private var viewModel = ComplaintViewModel()
    let fromUser: User
    let aboutUser: User
    private let localization = LocalizationManager.shared
    private let accentColor = Color(red: 0.0, green: 0.20, blue: 0.44)

    var body: some View {
        Form {
            Section(localization.text("complaint_about_section")) {
                Text(aboutUser.name)
                    .foregroundStyle(.secondary)
            }

            Section(localization.text("complaint_description_section")) {
                TextEditor(text: $viewModel.description)
                    .frame(minHeight: 120)
            }

            Section {
                Button(localization.text("complaint_submit_button")) {
                    viewModel.submit(from: fromUser, about: aboutUser)
                    dismiss()
                }
                .disabled(viewModel.description.isEmpty)
            }
        }
        .navigationTitle(localization.text("complaint_form_title"))
        .tint(accentColor)
    }
}

#Preview {
    NavigationStack {
        ComplaintFormView(
            fromUser: User(id: UUID(), name: "Нурлан С.", email: "n@s.kz", role: .advisor, specialty: nil),
            aboutUser: User(id: UUID(), name: "Ержан С.", email: "e@s.kz", role: .mentor, specialty: nil)
        )
    }
}

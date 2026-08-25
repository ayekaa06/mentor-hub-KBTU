//
//  MeetingCheckInView.swift
//  MentorsHub
//

import SwiftUI
import PhotosUI

struct MeetingCheckInView: View {
    @Environment(\.dismiss) private var dismiss
    @State private var viewModel: MeetingCheckInViewModel
    private let localization = LocalizationManager.shared
    private let accentColor = Color(red: 0.0, green: 0.20, blue: 0.44)

    init(mentor: User) {
        _viewModel = State(initialValue: MeetingCheckInViewModel(mentor: mentor))
    }

    var body: some View {
        Form {
            Section(localization.text("meeting_semester_section")) {
                TextField(localization.text("meeting_semester_placeholder"), text: $viewModel.semester)
            }

            Section(localization.text("meeting_photo_section")) {
                PhotosPicker(selection: $viewModel.selectedPhotoItem, matching: .images) {
                    Label(localization.text("meeting_pick_photo"), systemImage: "photo")
                }

                if let data = viewModel.photoData, let uiImage = UIImage(data: data) {
                    Image(uiImage: uiImage)
                        .resizable()
                        .scaledToFit()
                        .frame(maxHeight: 200)
                        .cornerRadius(10)
                }
            }

            Section {
                Button(localization.text("meeting_submit_button")) {
                    viewModel.submit()
                    dismiss()
                }
                .disabled(viewModel.photoData == nil)
            }

            if !viewModel.checkIns.isEmpty {
                Section(localization.text("meeting_history_section")) {
                    ForEach(viewModel.checkIns) { checkIn in
                        HStack {
                            Text(checkIn.semester)
                            Spacer()
                            Text(localization.text(checkIn.approved ? "meeting_approved" : "meeting_pending"))
                                .font(.caption)
                                .foregroundStyle(checkIn.approved ? .green : .orange)
                        }
                    }
                }
            }
        }
        .navigationTitle(localization.text("meeting_title"))
        .tint(accentColor)
    }
}

#Preview {
    NavigationStack {
        MeetingCheckInView(mentor: User(id: UUID(), name: "Тест", email: "test@test.kz", role: .mentor, specialty: nil))
    }
}

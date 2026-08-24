//
//  CheckInReviewView.swift
//  MentorsHub
//

import SwiftUI

struct CheckInReviewView: View {
    @State private var viewModel = CheckInReviewViewModel()
    private let localization = LocalizationManager.shared
    private let accentColor = Color(red: 0.0, green: 0.20, blue: 0.44)

    var body: some View {
        List {
            if viewModel.checkIns.isEmpty {
                Text(localization.text("checkin_review_empty"))
                    .foregroundStyle(.secondary)
            }
            ForEach(viewModel.checkIns) { checkIn in
                VStack(alignment: .leading, spacing: 6) {
                    HStack {
                        Text(checkIn.mentorName).font(.headline)
                        Spacer()
                        Text(localization.text(checkIn.approved ? "meeting_approved" : "meeting_pending"))
                            .font(.caption)
                            .foregroundStyle(checkIn.approved ? .green : .orange)
                    }
                    Text(checkIn.semester)
                        .font(.subheadline)
                        .foregroundStyle(.secondary)
                    Text(checkIn.date.formatted(date: .abbreviated, time: .omitted))
                        .font(.caption)
                        .foregroundStyle(.secondary)

                    if let data = checkIn.photoData, let uiImage = UIImage(data: data) {
                        Image(uiImage: uiImage)
                            .resizable()
                            .scaledToFit()
                            .frame(maxHeight: 160)
                            .cornerRadius(10)
                    }

                    if !checkIn.approved {
                        Button(localization.text("checkin_review_approve_button")) {
                            viewModel.approve(checkIn)
                        }
                        .buttonStyle(.bordered)
                        .padding(.top, 4)
                    }
                }
                .padding(.vertical, 4)
            }
        }
        .navigationTitle(localization.text("checkin_review_title"))
        .tint(accentColor)
        .onAppear { viewModel.load() }
    }
}

#Preview {
    NavigationStack {
        CheckInReviewView()
    }
}

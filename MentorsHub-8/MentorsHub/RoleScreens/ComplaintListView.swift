//
//  ComplaintListView.swift
//  MentorsHub
//

import SwiftUI

struct ComplaintListView: View {
    @State private var viewModel = ComplaintViewModel()
    private let localization = LocalizationManager.shared
    private let accentColor = Color(red: 0.0, green: 0.20, blue: 0.44)

    var body: some View {
        List {
            ForEach(viewModel.complaints) { complaint in
                VStack(alignment: .leading, spacing: 6) {
                    HStack {
                        Text(complaint.aboutUserName)
                            .font(.headline)
                        Spacer()
                        Text(localization.text(statusKey(for: complaint.status)))
                            .font(.caption)
                            .foregroundStyle(complaint.status == .pending ? .orange : .secondary)
                    }
                    Text(complaint.fromUserName)
                        .font(.caption)
                        .foregroundStyle(.secondary)
                    Text(complaint.description)
                        .font(.subheadline)

                    if complaint.status == .pending {
                        HStack {
                            Button(localization.text("complaint_resolve_button")) { viewModel.resolve(complaint, status: .resolved) }
                                .buttonStyle(.bordered)
                            Button(localization.text("complaint_dismiss_button")) { viewModel.resolve(complaint, status: .dismissed) }
                                .buttonStyle(.bordered)
                        }
                        .padding(.top, 4)
                    }
                }
                .padding(.vertical, 4)
            }
        }
        .navigationTitle(localization.text("complaint_list_title"))
        .tint(accentColor)
        .onAppear { viewModel.load() }
    }

    private func statusKey(for status: ComplaintStatus) -> String {
        switch status {
        case .pending: "complaint_status_pending"
        case .resolved: "complaint_status_resolved"
        case .dismissed: "complaint_status_dismissed"
        }
    }
}

#Preview {
    NavigationStack {
        ComplaintListView()
    }
}

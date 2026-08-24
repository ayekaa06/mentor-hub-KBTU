//
//  MenteeDetailView.swift
//  MentorsHub
//

import SwiftUI

struct MenteeDetailView: View {
    let mentee: Mentee
    let viewModel: MentorViewModel
    private let localization = LocalizationManager.shared
    private let accentColor = Color(red: 0.0, green: 0.20, blue: 0.44)

    private var currentStatus: MenteeStatus {
        viewModel.myMentees.first(where: { $0.id == mentee.id })?.status ?? mentee.status
    }

    var body: some View {
        List {
            Section(localization.text("mentee_detail_info_section")) {
                HStack {
                    Text(localization.text("mentee_detail_email"))
                    Spacer()
                    Text(mentee.user.email).foregroundStyle(.secondary)
                }
                HStack {
                    Text(localization.text("mentee_detail_group"))
                    Spacer()
                    Text(mentee.group).foregroundStyle(.secondary)
                }
                HStack {
                    Text(localization.text("mentee_detail_specialty"))
                    Spacer()
                    Text(mentee.specialty).foregroundStyle(.secondary)
                }
                if let date = mentee.dateAssigned {
                    HStack {
                        Text(localization.text("mentee_detail_assigned_date"))
                        Spacer()
                        Text(date.formatted(date: .abbreviated, time: .omitted)).foregroundStyle(.secondary)
                    }
                }
            }

            Section(localization.text("mentee_detail_status_section")) {
                Picker(localization.text("mentee_detail_status_section"), selection: Binding(
                    get: { currentStatus },
                    set: { viewModel.updateStatus(for: mentee, to: $0) }
                )) {
                    ForEach(MenteeStatus.allCases, id: \.self) { status in
                        Text(status.displayName).tag(status)
                    }
                }
                .pickerStyle(.segmented)
            }
        }
        .navigationTitle(mentee.user.name)
        .tint(accentColor)
    }
}

#Preview {
    NavigationStack {
        MenteeDetailView(
            mentee: Mentee(id: UUID(), user: User(id: UUID(), name: "Алина К.", email: "a@k.kz", role: .mentee, specialty: "26BDIS"), specialty: "26BDIS", group: "IS-2601", assignedTo: nil, status: .active, dateAssigned: Date()),
            viewModel: MentorViewModel(mentorSpecialty: "26BDIS")
        )
    }
}

//
//  DocumentManagementView.swift
//  MentorsHub
//

import SwiftUI
import UniformTypeIdentifiers

struct DocumentManagementView: View {
    @State private var viewModel = DocumentManagementViewModel()
    @State private var showFileImporter = false
    private let localization = LocalizationManager.shared
    private let accentColor = Color(red: 0.0, green: 0.20, blue: 0.44)

    var body: some View {
        List {
            Section(localization.text("document_new_section")) {
                TextField(localization.text("document_title_placeholder"), text: $viewModel.newTitle)
                Picker(localization.text("document_type_label"), selection: $viewModel.newType) {
                    Text(localization.text("document_type_handbook")).tag(DocumentType.handbook)
                    Text(localization.text("document_type_codex")).tag(DocumentType.codex)
                    Text(localization.text("document_type_process")).tag(DocumentType.process)
                    Text(localization.text("document_type_rup")).tag(DocumentType.rup)
                }
                if viewModel.newType == .handbook {
                    Picker(localization.text("document_language_label"), selection: $viewModel.newLanguage) {
                        ForEach(Language.allCases) { lang in
                            Text(lang.displayName).tag(lang)
                        }
                    }
                }
                if viewModel.newType == .rup {
                    TextField(localization.text("document_specialty_placeholder"), text: $viewModel.newSpecialty)
                }
                Button(localization.text("document_upload_button")) {
                    showFileImporter = true
                }
                .disabled(viewModel.newTitle.isEmpty)
            }

            Section(localization.text("document_list_section")) {
                ForEach(viewModel.documents) { doc in
                    HStack {
                        Label(doc.title, systemImage: "doc.text")
                        Spacer()
                        if let lang = doc.language {
                            Text(lang.displayName).font(.caption).foregroundStyle(.secondary)
                        }
                        if let specialty = doc.specialty {
                            Text(specialty).font(.caption).foregroundStyle(.secondary)
                        }
                    }
                }
            }
        }
        .navigationTitle(localization.text("document_title"))
        .tint(accentColor)
        .onAppear { viewModel.load() }
        .fileImporter(isPresented: $showFileImporter, allowedContentTypes: [.pdf]) { result in
            if case .success(let url) = result {
                viewModel.upload(pickedURL: url)
            }
        }
    }
}

#Preview {
    NavigationStack {
        DocumentManagementView()
    }
}

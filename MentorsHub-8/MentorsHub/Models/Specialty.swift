//
//  Specialty.swift
//  MentorsHub
//
//  Общий справочник специальностей — используется и для менти
//  (какая у него специальность), и для ментора (по какой специальности
//  он менторит — от этого зависит, каких менти отдаёт рандомайзер).
//

import Foundation

struct Specialty: Identifiable, Hashable {
    let code: String
    let displayName: String
    var id: String { code }

    static let all: [Specialty] = [
        Specialty(code: "26BDIS", displayName: "ИС — Информационные системы"),
        Specialty(code: "26BDCS", displayName: "ВТиПО — Выч. техника и ПО")
    ]
}

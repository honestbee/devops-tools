#!/bin/bash
terraform import module.ghac.github_team.{{ .Team.Slug }} {{ $.Team.ID }}

{{- range $index, $user := .SortedUsers}} 
terraform import module.ghac.github_team_membership.{{ $.Team.Slug }}-{{ $user }} {{ $.Team.ID }}:{{ $user }}
{{- end -}}

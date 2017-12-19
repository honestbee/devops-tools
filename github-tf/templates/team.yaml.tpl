teams:
{{- range $tr := .TeamRoles}}
- name:         "{{.Team.Name}}"
  id:           {{.Team.ID}}
  description:  "{{.Team.Description}}" 
  slug:         {{.Team.Slug}}
  privacy:      {{.Team.Privacy}}
  user_roles:
{{- range $user := .SortedUsers}}
    {{ $user }}: {{ index $tr.UserRoles $user -}}
{{end -}}
{{end -}}

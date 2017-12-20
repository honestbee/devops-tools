resource "github_team" "{{.Slug}}" {
  name        = "{{.Name}}"
  description = "{{.Description}}"
  privacy     = "{{.Privacy}}"
  # {{ .SlugPrefix }}
  # {{ .SlugSuffix }}
}
{{ range $user, $role := .UserRoles }}
resource "github_team_membership" "{{ $.Slug }}-{{ $user }}" {
  team_id  = "{{ $.ID }}"
  username = "{{ $user }}"
  role = "{{ $role }}"
}
{{ end }}
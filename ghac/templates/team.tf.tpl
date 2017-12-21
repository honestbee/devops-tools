resource "github_team" "{{.Slug}}" {
  name        = "{{.Name}}"
  description = "{{.Description}}"
  privacy     = "{{ default "closed" .Privacy }}"
  # Test Prefix: {{ .SlugPrefix }}
  # Test Suffix: {{ .SlugSuffix }}
  # Test Sprig.default: {{ default "foo" .Privacy }} 
}
{{ range $user, $role := .UserRoles }}
resource "github_team_membership" "{{ $.Slug }}-{{ $user }}" {
  team_id  = "{{ $.ID }}"
  username = "{{ $user }}"
  role = "{{ $role }}"
}
{{ end }}
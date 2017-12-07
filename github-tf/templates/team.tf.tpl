resource "github_team" "{{.Team.Slug}}" {
  name        = "{{.Team.Slug}}"
  description = "{{.Team.Description}}"
  privacy     = "{{.Team.Privacy}}"
}

module "{{.Team.Slug}}" {
  source     = "git@github.com:honestbee/tf-modules.git?ref=master//github/team"
  user_roles = {
{{range $user := .SortedUsers}}    {{ $user }} = "{{ index $.UserRoles $user }}"
{{end}}  }

  team_id = "{{.Team.ID}}"
}

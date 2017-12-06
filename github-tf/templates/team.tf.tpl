module "{{.Team.Slug}}" {
  source     = "git@github.com:honestbee/tf-modules.git?ref=master//github/team"
  user_roles = {
{{range $user, $role := .UserRoles}}		{{$user }} = "{{$role}}"
{{end}}  }

  team_id = "{{.Team.ID}}"
}

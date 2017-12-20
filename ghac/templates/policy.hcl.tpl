path "secret/*" {
    capabilities = ["list"]
}

path "secret/{{.Env}}/{{.ShortName}}*" {
    capabilities = ["create", "read", "update", "list"]
}

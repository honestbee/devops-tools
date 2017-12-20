path "secret/*" {
    capabilities = ["list"]
}

path "secret/{{ env }}/{{ $team }}*" {
    capabilities = ["create", "read", "update", "list"]
}

# go-project

Squelette d'une application e-commerce en Go, organisé selon les couches
domain, application, infrastructure et delivery. Les fichiers de structure
ne contiennent pas encore la logique métier.

## Prérequis

- Go 1.27 ou supérieur

## Lancer le point d'entrée historique

```sh
go run .
```

## Vérifier et compiler le projet

```sh
go test ./...
go build ./...
```

Le serveur s'exécute depuis `cmd/server` ; il s'agit pour le moment d'un
point d'entrée vide :

```sh
go run ./cmd/server
```

Les autres exécutables sont disponibles dans `cmd/client`, `cmd/admin`,
`cmd/ssh-client` et `cmd/ssh-admin`.

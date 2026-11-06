# go-project

Squelette d'une application e-commerce en Go, organisé selon les couches
domain, application, infrastructure et delivery. Les fichiers de structure
ne contiennent pas encore la logique métier.

## Base de données

La base PostgreSQL et son schéma sont initialisés avec Docker Compose. Copiez
`.env.example` vers `.env`, puis démarrez le service :

```sh
docker compose up -d db
```

Les scripts `migrations/001` à `migrations/005` sont exécutés dans l'ordre
lors de la première création du volume `postgres_data`. Si vous modifiez les
migrations après cette initialisation, elles ne sont pas rejouées
automatiquement : utilisez un outil de migration ou une base/volume neuf.
Les identifiants de développement sont configurables dans `.env`.

Arrêter le service sans supprimer les données :

```sh
docker compose down
```

## Prérequis

- Go 1.27 ou supérieur
- Docker Compose

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

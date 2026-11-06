# go-project

Squelette d'une application e-commerce en Go, organisé selon les couches
domain, application, infrastructure et delivery. Les fichiers de structure
ne contiennent pas encore la logique métier.

Les conventions de code et d'architecture sont décrites dans
[CONVENTIONS.md](CONVENTIONS.md).

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

Le serveur lit également les variables `POSTGRES_HOST`, `POSTGRES_PORT`,
`POSTGRES_DB`, `POSTGRES_USER`, `POSTGRES_PASSWORD` et `POSTGRES_SSLMODE`.
Il démarre le serveur HTTP sur `HTTP_ADDR` (par défaut `:8080`) et expose
`GET /health` ainsi que les routes d'authentification ci-dessous. Après avoir
démarré la base, lancez `go run ./cmd/server` puis ouvrez
`http://localhost:8080/health`. Les valeurs par défaut de PostgreSQL
correspondent à la configuration locale Docker Compose ; remplacez le mot de
passe de développement hors environnement local.

### Authentification et utilisateurs

- `POST /api/auth/register` — crée un compte client (`email`, `password`).
- `POST /api/auth/login` — vérifie les identifiants et renvoie un jeton Bearer
  valable 24 heures.
- `GET /api/users/me` — renvoie le profil authentifié.
- `POST /api/auth/logout` — révoque le jeton Bearer utilisé.

Exemple :

```sh
curl -X POST http://localhost:8080/api/auth/register \
  -H "Content-Type: application/json" \
  -d '{"email":"ada@example.com","password":"correct-horse-battery"}'
```

La connexion s'effectue avec les mêmes champs JSON :

```sh
curl -X POST http://localhost:8080/api/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email":"ada@example.com","password":"correct-horse-battery"}'
```

Le jeton `access_token` renvoyé par la connexion s'envoie dans l'en-tête
`Authorization: Bearer <jeton>` pour consulter le profil ou se déconnecter.
Une adresse e-mail valide et un mot de passe de 8 à 256 octets sont requis.
Les mots de passe sont stockés sous forme de hachages PBKDF2-SHA256 salés et
les jetons ne sont stockés qu'avec leur empreinte. Les réponses utilisateur
n'exposent jamais le hachage du mot de passe.

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

Le serveur s'exécute depuis `cmd/server` :

```sh
go run ./cmd/server
```

Les autres exécutables sont disponibles dans `cmd/client`, `cmd/admin`,
`cmd/ssh-client` et `cmd/ssh-admin`.

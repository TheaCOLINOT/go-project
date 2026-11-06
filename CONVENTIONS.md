# Conventions Go

Ces conventions s'appliquent au projet e-commerce et complètent les conventions
idiomatiques de Go. En cas de doute, privilégier la simplicité et les pratiques
standard de Go.

## Structure et responsabilités

- `cmd/<nom>/` contient le point d'entrée de chaque exécutable. `main` assemble
  les dépendances et démarre l'application ; la logique métier reste dans
  `internal/`.
- `internal/domain/<domaine>/` contient les modèles, règles métier, erreurs
  métier et interfaces nécessaires au domaine.
- `internal/application/<domaine>/` orchestre les cas d'usage et dépend des
  abstractions du domaine, pas des détails techniques.
- `internal/infrastructure/` contient les implémentations techniques, par
  exemple l'accès PostgreSQL, l'envoi d'e-mails et la sécurité.
- `internal/delivery/` adapte les entrées/sorties (HTTP, CLI, SSH) et délègue
  les cas d'usage à la couche application.
- `internal/config/` centralise le chargement et la validation de la
  configuration.
- Les migrations SQL restent dans `migrations/` et sont numérotées pour
  conserver leur ordre d'application.
- Garder les dépendances orientées vers le domaine : celui-ci ne dépend ni de
  l'infrastructure ni des mécanismes de delivery.

## Packages, fichiers et noms

- Nommer les packages en minuscules, sans underscore ni nom générique (`util`,
  `common`, `misc`). Choisir des noms courts et précis, généralement au
  singulier : `user`, `product`, `postgres`.
- Éviter les noms de package redondants dans les identifiants : dans le package
  `user`, préférer `User` à `UserEntity`.
- Utiliser `PascalCase` pour les identifiants exportés et `camelCase` pour les
  identifiants internes. Les acronymes connus gardent leur casse Go : `ID`,
  `URL`, `HTTP` (par exemple `userID`, `HTTPServer`).
- Choisir des noms concis mais explicites. Les abréviations courantes de Go
  (`ctx`, `err`, `db`) sont acceptables ; éviter les abréviations inventées.
- Nommer les fichiers selon leur contenu et leur rôle (`entity.go`,
  `repository.go`, `routes.go`). Ajouter `_test.go` pour les tests.
- Ne pas imposer un suffixe à toutes les interfaces ou fonctions : les noms
  doivent décrire le comportement et rester naturels dans leur package.

## Modèle métier et dépendances

- Les modèles du domaine représentent le métier, sans dépendance aux types
  HTTP, SQL ou aux structures de transport.
- Définir les interfaces là où elles sont consommées, avec uniquement les
  méthodes requises. Une interface de dépôt métier peut par exemple vivre dans
  `internal/domain/user`, tandis que son implémentation PostgreSQL vit dans
  `internal/infrastructure/persistence/postgres`.
- Injecter les dépendances explicitement. Utiliser un constructeur `New...`
  lorsqu'il améliore la création ou garantit un état valide ; éviter les
  variables globales et les couches d'abstraction sans besoin concret.
- Les handlers valident et convertissent les entrées/sorties propres au
  transport ; les règles métier appartiennent au domaine ou aux cas d'usage.

## Erreurs et contexte

- Retourner les erreurs attendues explicitement ; ne pas utiliser `panic` pour
  un échec métier ou opérationnel récupérable.
- Ajouter du contexte aux erreurs propagées avec `%w` :
  `fmt.Errorf("user repository: find by ID: %w", err)`.
- Définir des erreurs métier reconnaissables (sentinelles ou types) quand
  l'appelant doit adapter son comportement. Les comparer avec `errors.Is` ou
  `errors.As`, pas avec leur texte.
- Ne pas ignorer silencieusement une erreur. La traiter, la retourner ou, si
  elle est volontairement ignorée, rendre ce choix explicite.
- Passer `context.Context` en premier paramètre aux opérations qui font de
  l'I/O ou doivent être annulables. Propager le contexte aux requêtes SQL et
  appels externes ; ne pas le stocker dans une structure.

## Sécurité et données

- Ne jamais journaliser mots de passe, jetons, secrets ou données sensibles.
  Stocker uniquement des empreintes de jetons et des mots de passe hachés avec
  un mécanisme adapté.
- Utiliser des requêtes SQL paramétrées, jamais de concaténation de données
  utilisateur dans une requête.
- Valider les entrées à la frontière du système et faire respecter également
  les invariants métier dans le domaine. Ne pas considérer une validation de
  transport comme une garantie suffisante.
- Garder les secrets dans la configuration d'environnement appropriée, jamais
  dans le code ou les migrations.

## Style et commentaires

- Formater chaque changement Go avec `gofmt`. Garder les imports gérés par les
  outils Go et supprimer les imports inutilisés.
- Préférer un code direct, lisible et découpé par responsabilité. N'ajouter une
  abstraction que lorsqu'elle clarifie le code ou répond à un besoin réel.
- Documenter les identifiants exportés dont le rôle n'est pas évident. Le
  commentaire commence par le nom de l'identifiant, par exemple
  `// User représente un compte utilisateur.`.
- Écrire les noms du code en anglais, conformément aux fichiers existants.
  Garder les explications destinées à l'équipe en français si elles sont
  documentaires.

## Tests et vérifications

- Placer les tests unitaires à côté du code testé dans des fichiers
  `*_test.go`. Utiliser des tests table-driven lorsque plusieurs cas partagent
  la même structure.
- Tester les comportements observables : succès, erreurs, cas limites et
  règles métier. Les tests ne doivent pas dépendre d'un ordre d'exécution.
- Isoler les tests d'intégration et de bout en bout dans `tests/` et indiquer
  clairement leurs dépendances externes.
- Avant de proposer un changement, exécuter :

  ```sh
  gofmt -w chemin/vers/fichier.go
  go test ./...
  go build ./...
  ```

- Après toute modification de migration, vérifier le schéma et son ordre
  d'application ; les migrations déjà appliquées ne sont pas rejouées
  automatiquement sur un volume PostgreSQL existant.

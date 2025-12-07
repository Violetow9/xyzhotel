# XYZ Hotel - Projet DDD en Go

Ce projet est une application de gestion hôtelière réalisée en Go. 
L'objectif principal était de mettre en pratique le DDD.

L'application permet de gérer des clients, des portefeuilles électroniques et le cycle complet de réservation de chambres (disponibilité, acompte, confirmation, annulation).

Ce projet a été réalisé dans le cadre d'un exercice technique sur le DDD.
Lien vers l'énoncé complet : [XYZ Hotel - Exercice DDD](https://github.com/sallareznov/ddd/blob/main/projet.md)

L'ubiquitous language se trouve dans le fichier [DESIGN.md](./DESIGN.md).
Pour visualiser le design stratégique, un [diagramme Mermaid](./DESIGN.md) est également disponible.
Sur intellij ou GoLand, installer le plugin "Mermaid" pour le visualiser directement dans l'ide.

## Technos

* **Langage** : Go (Golang)
* **Architecture** : Hexagonale (Ports & Adapters)
* **Base de données** : MySQL (via Docker)
* **Requêtes SQL** : Générées avec `sqlc` pour la sécurité de typage.
* **API HTTP** : Framework `Gin`.
* **CLI** : Interface en ligne de commande native.

## Pour lancer le projet

Il faut utiliser le `Makefile` qui simplifie toutes les commandes.

### Pré-requis
* Go installé.
* Docker et Docker Compose installés et lancés.

### 1. Démarrer la base de données
Avant toute chose, il faut lancer le conteneur MySQL.

```bash
make docker-db
```

(Attendre quelques secondes que la base s'initialise). 
Note : Si vous voulez réinitialiser la base de données à zéro (supprimer les données et relancer le seed), 
utilisez: 

```bash
make docker-reset
```
### 2. Lancer l'application

Pour lancer uniquement le serveur HTTP :
```bash
make run
```

Pour lancer l'interface en ligne de commande (CLI) :
```bash
make run-cli
```

Pour lancer les deux en même temps (serveur HTTP + CLI) :
Il suffit d'ouvrir 2 terminal et de faire les 2 commandes.

### 3. Lancer les tests

```bash
make test
```
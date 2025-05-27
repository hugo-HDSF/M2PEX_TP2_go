# TP : GoLog Analyzer - Analyse de Logs Distribuée

### Contexte

Votre équipe est chargée de développer un outil en ligne de commande (CLI) en Go, nommé `loganalyzer`. Son but est d'aider les administrateurs système à analyser des fichiers de logs (journaux) provenant de diverses sources (serveurs, applications). L'objectif est de pouvoir centraliser l'analyse de multiples logs en parallèle et d'en extraire des informations clés, tout en gérant les erreurs potentielles de manière robuste.

### Objectifs d'apprentissage

Ce TP vous permettra de renforcer vos compétences sur les concepts suivants :

- **Concurrence :** Utiliser les **goroutines** et les **WaitGroups** pour traiter plusieurs tâches en parallèle.
- **Gestion des Erreurs :** Implémenter des **erreurs personnalisées** et les gérer proprement avec `errors.Is` et `errors.As`.
- **Outil CLI avec Cobra :** Structurer une application en ligne de commande avec des **sous-commandes** et des **drapeaux (flags)**.
- **Import/Export JSON :** Manipuler des données au format JSON pour la configuration d'entrée et le rapport de sortie.
- **Modularité :** Organiser le code en **packages** logiques (`internal/`).

---

### Cahier des charges

Votre outil `loganalyzer` devra implémenter les fonctionnalités suivantes :

#### 1. Commande `analyze`

- **Entrée JSON :** La commande prendra un chemin vers un **fichier de configuration JSON** via un drapeau `--config <path>` (raccourci `-c`). Ce fichier contiendra la liste des logs à analyser.

  **Exemple de fichier `config.json` :**
    ```json
    [
      {
        "id": "web-server-1",
        "path": "/var/log/nginx/access.log",
        "type": "nginx-access"
      },
      {
        "id": "app-backend-2",
        "path": "/var/log/my_app/errors.log",
        "type": "custom-app"
      }
    ]
    ```
  - `id` : Un identifiant unique pour le log.
  - `path` : Le chemin (absolu ou relatif) vers le fichier de log.
  - `type` : Le type de log (peut être ignoré mais doit être présent).

- **Traitement concurrentiel :** Une **goroutine** sera lancée pour chaque log :
  - Vérifier si le fichier existe et est lisible.
  - Simuler l'analyse avec un `time.Sleep` aléatoire (50 à 200 ms).
  - **Erreur aléatoire simulée :** 10% de chance de générer une erreur de parsing.

- **Collecte et Exportation des résultats :**
  - Résultats collectés via un **canal sécurisé**.
  - Export possible via `--output <path>` (raccourci `-o`) dans un fichier JSON.

    **Exemple de fichier `report.json` :**
    ```json
    [
      {
        "log_id": "web-server-1",
        "file_path": "/var/log/nginx/access.log",
        "status": "OK",
        "message": "Analyse terminée avec succès.",
        "error_details": ""
      },
      {
        "log_id": "invalid-path",
        "file_path": "/non/existent/log.log",
        "status": "FAILED",
        "message": "Fichier introuvable.",
        "error_details": "open /non/existent/log.log: no such file or directory"
      }
    ]
    ```

- **Affichage sur console :** Un résumé doit être affiché pour chaque log : ID, chemin, statut, message, erreur (si applicable).

#### 2. Gestion des Erreurs Personnalisées

- Implémenter au moins **deux types d'erreurs personnalisées** :
  - Fichier introuvable/inaccessible.
  - Erreur de parsing.
- Utiliser `errors.Is()` et `errors.As()` pour les gérer proprement.

---

### Architecture suggérée (packages `internal/`)

Organisez le projet comme suit :

- `internal/config` : Lecture des configurations JSON.
- `internal/analyzer` : Analyse, erreurs personnalisées, rapport.
- `internal/reporter` : Export JSON des résultats.
- `cmd/` :
  - `root.go` : Commande racine.
  - `analyze.go` : Commande `analyze`.

---

### Critères d'évaluation

L’évaluation portera sur :

- **Fonctionnalité :** La commande `analyze` fonctionne-t-elle comme spécifié ?
- **Concurrence :** Traitement en parallèle via `goroutines` et `WaitGroup` ? Résultats collectés via `channel` ?
- **Gestion des Erreurs :** Utilisation et gestion correcte des erreurs personnalisées ? Messages d’erreur clairs ?
- **CLI :** Interface Cobra fonctionnelle, avec drapeaux et descriptions ?
- **JSON :** Import/export respectant les structures attendues ?
- **Modularité :** Code organisé proprement en packages ?
- **Documentation :** Je veux voir **un beau readme** qui explique le fonctionnement de votre programme, vos commandes, et j'en passe ET **la documentation de votre code** et **les membres de votre team**.

### Type de rendu

- Un lien github


### 🎁 BONUS

Vous avez l'âme d'un.e développeur.euse courageux.euse ? Je vous laisse ici quelques bonus si vous voulez vous amuser un peu et avoir un programme plus complet.

**1. Gestion des dossiers d'exportation **
* Si le chemin de sortie JSON (`--output`) inclut des répertoires qui n'existent pas (ex: `rapports/2024/mon_rapport.json`), faire en sorte que le programme crée automatiquement ces répertoires avant d'écrire le fichier.
* **Indice** : `os.MkdirAll(filepath.Dir(path), 0755)`
* **Intérêt** : Rend l'outil plus robuste et convivial.

**2. Horodatage des Exports JSON**
* Nommer les fichiers de sortie JSON avec une date :
  * **Modifier la commande `analyze`** pour que, si le drapeau `--output` est fourni, le nom du fichier de sortie JSON inclue la date du jour au format AAMMJJ (AnnéeMoisJour).
  * **Exemple** : au lieu de `report.json`, le fichier serait nommé `240524_report.json` (pour le 24 mai 2024).
  * **Indice** : Utiliser le package `time` de Go (`time.Now()`, `time.Format()`).
  * **Intérêt** : Ajoute une fonctionnalité pratique pour l'organisation des rapports, et force à manipuler les dates en Go.

**2. Commande `add-log`**
* **Ajouter une nouvelle sous-commande add-log** qui permettrait d'ajouter manuellement une configuration de log au fichier config.json existant.
* **Drapeaux nécessaires** : `--id`, `--path`, `--type`, `--file` (chemin du fichier `config.json`).

**3. Filtrage des résultats d'analyse**
* **Ajouter un drapeau `--status <status>`** (ex: `--status FAILED` ou `--status OK`) à la commande analyze pour n'afficher et/ou n'exporter que les logs ayant un certain statut.
* **Intérêt** : Ajoute une fonctionnalité utile et demande de la logique de filtrage avant l'affichage/l'export.


---

### Pour démarrer (Prérequis)

1. Créer un module : `go mod init`
2. Installer Cobra : `go get github.com/spf13/cobra@latest`
3. Avoir bien lu le readme ;)

---

Bon courage !

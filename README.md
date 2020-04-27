# terraform-graph-beautifier
Terraform graph beautifier

## TODO
- Dans samples, configs TF de tests
  - TF only (local file, random) avec:
    - modules imbriqués
    - data
    - variables, locals et outputs
    - maps et lists
  - Utilisation des modules GCP
- Config file
```yaml
tf-types:
  - type (resource, provider, var, local, output) :
    - shape: box(ressource), arrondi(var, local), ?? (output)
    - bg-colors: [...]
    - fg-colors: [...]
    - arrow-to-shape: ""
```
- Param pour générer un fichier de config avec valeurs par défaut. 
- Customs attributes
    - tf-type
- Voir (orphan), (close) (removed)
- Utilisation `name` plutôt que `label`
- Param `--exclude-prefix`
- Name : remplacer "[root] " par "root."
- Exemples utilisation [`sfdp`](https://linux.die.net/man/1/sfdp)
- Github actions to build
- Virer junk (--keep-tf-junk)

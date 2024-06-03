Pour construire le main, simplement utiliser le fichier `build.py`
Pour chercher à démontrer une relation *x y z*, utiliser `main.exe -script demo "x" "y" "z"` (attention les epsaces doivent être notés `%20`)
Pour trouver tous les *z* possibles dans une relation *x y z*, utiliser `main.exe -script search "x" "y" "..."`. Utiliser `positive-weights` pour filter les relations à poids positif et `-lim n` pour afficher uniquement les `n` premiers candidats.

#!/usr/bin/env python3
"""
Script pour ajouter des commentaires avant tous les returns dans les fichiers Go.
"""

import re
import sys

# Mapping fichier -> {ligne: commentaire}
fixes = {
    "/workspace/src/pkg/analyzer/func.go": {
        80: "\t\t\t// Retourne true car le fichier est un fichier de test\n",
        83: "\t// Retourne false car le fichier n'est pas un fichier de test\n",
        124: "\t\t// Retourne immédiatement car pas de paramètres à vérifier\n",
        185: "\t\t// Retourne immédiatement car la fonction n'a pas de body\n",
        214: "\t\t\t// Retourne false pour arrêter l'inspection de cette branche\n",
        217: "\t\t\t// Retourne false pour arrêter l'inspection de cette branche\n",
        220: "\t\t\t// Retourne false pour arrêter l'inspection de cette branche\n",
        223: "\t\t\t// Retourne false pour arrêter l'inspection de cette branche\n",
        226: "\t\t\t// Retourne false pour arrêter l'inspection de cette branche\n",
        228: "\t\t// Retourne true pour continuer l'inspection\n",
        231: "\t// Retourne la profondeur maximale trouvée\n",
        246: "\t\t// Retourne la nouvelle profondeur si elle est supérieure\n",
        248: "\t// Retourne la profondeur actuelle sinon\n",
        266: "\t\t// Retourne immédiatement si le format est correct\n",
        286: "\t\t// Retourne immédiatement si pas de params ni returns\n",
        382: "\t// Retourne le contenu de la section extraite\n",
        398: "\t// Retourne l'exemple formaté de section Params\n",
        414: "\t// Retourne l'exemple formaté de section Returns\n",
        434: "\t// Retourne le nombre total de paramètres\n",
        453: "\t// Retourne la liste des noms de paramètres\n",
        466: "\t\t// Retourne 0 car pas de body\n",
        486: "\t\t// Retourne 1 car pas de body\n",
        492: "\t\t// Retourne true pour continuer l'inspection\n",
        495: "\t// Retourne la complexité totale\n",
        508: "\t\t// Retourne 1 pour un if statement\n",
        510: "\t\t// Retourne 1 pour un for ou range statement\n",
        513: "\t\t\t// Retourne 1 car le case a une liste\n",
        517: "\t\t\t// Retourne 1 car le comm clause a une comm\n",
        521: "\t\t\t// Retourne 1 pour un opérateur logique\n",
        524: "\t// Retourne 0 pour les autres nœuds\n",
    },
    "/workspace/src/pkg/analyzer/interface.go": {
        70: "\t\t// Retourne immédiatement car le package est exempté\n",
        91: "\t// Retourne nil car l'analyseur rapporte via pass.Reportf\n",
        133: "\t// Retourne les informations collectées\n",
        212: "\t\t// Retourne immédiatement car c'est une méthode\n",
        233: "\t\t// Retourne 0 car pas de méthodes\n",
        235: "\t// Retourne le nombre de méthodes\n",
        253: "\t\t\t// Retourne true car le package est exempté\n",
        257: "\t// Retourne false car le package n'est pas exempté\n",
        323: "\t\t\t\t// Retourne true car le package a besoin d'interfaces.go\n",
        330: "\t\t\t// Retourne true car le package a des interfaces publiques\n",
        429: "\t\t\t// Retourne true car le type est autorisé\n",
        433: "\t// Retourne false car le type n'est pas dans la liste autorisée\n",
    },
    "/workspace/src/pkg/analyzer/interface_test.go": {
        511: "\t// Retourne true car la sous-chaîne a été trouvée\n",
        525: "\t\t\t// Retourne true car la sous-chaîne a été trouvée\n",
        528: "\t// Retourne false car la sous-chaîne n'a pas été trouvée\n",
    },
    "/workspace/src/pkg/analyzer/test.go": {
        58: "\t\t// Retourne immédiatement car le package est exempté\n",
        70: "\t// Retourne nil car l'analyseur rapporte via pass.Reportf\n",
        98: "\t// Retourne les informations collectées\n",
        119: "\t\t\t// Retourne true car le fichier contient des tests\n",
        122: "\t// Retourne false car le fichier ne contient pas de tests\n",
        161: "\t\t\t// Retourne le fichier AST trouvé\n",
        164: "\t// Retourne nil car le fichier n'a pas été trouvé\n",
        176: "\t\t// Retourne immédiatement si on est dans un package _test\n",
        219: "\t\t// Retourne immédiatement si on n'est pas dans un package _test\n",
        314: "\t// Retourne true si le fichier existe\n",
    },
    "/workspace/src/pkg/analyzer/test_export_test.go": {
        21: "\t// Retourne le résultat de l'exécution\n",
        33: "\t// Retourne le fichier AST trouvé ou nil\n",
        46: "\t// Retourne une fileInfo pour utilisation interne\n",
    },
    "/workspace/src/pkg/analyzer/test_test.go": {
        288: "\t// Retourne le mock FileSystem créé\n",
        366: "\t// Retourne true si le fichier existe\n",
    },
    "/workspace/src/pkg/analyzer/var.go": {
        46: "\t// Retourne nil car l'analyseur rapporte via pass.Reportf\n",
        80: "\t\t\t// Retourne true pour continuer l'inspection\n",
        90: "\t\t// Retourne true pour continuer l'inspection\n",
        241: "\t\t\t// Retourne true car un commentaire individuel existe\n",
        244: "\t\t// Retourne true car un commentaire en ligne existe\n",
        246: "\t// Retourne false car aucun commentaire individuel trouvé\n",
        271: "\t\t// Retourne immédiatement si pas un channel make\n",
        275: "\t\t// Retourne immédiatement si unbuffered est intentionnel\n",
        292: "\t\t// Retourne false car pas de type\n",
        297: "\t\t// Retourne false car pas un channel\n",
        301: "\t\t// Retourne false car pas de valeur initiale\n",
        306: "\t\t// Retourne false car pas un appel make\n",
        310: "\t// Retourne true si c'est un channel make sans buffer\n",
        332: "\t// Retourne true si unbuffered est mentionné dans le commentaire\n",
        345: "\t\t// Retourne immédiatement si pas de valeur initiale ou pas de type\n",
        350: "\t\t// Retourne immédiatement si le type n'est pas compatible const\n",
        355: "\t\t// Retourne immédiatement si la valeur n'est pas littérale\n",
        360: "\t\t// Retourne immédiatement si la variable est réassignée\n",
    },
}

def add_comment_before_line(file_path, line_num, comment):
    """Ajoute un commentaire avant une ligne spécifique."""
    try:
        with open(file_path, 'r') as f:
            lines = f.readlines()

        if line_num <= len(lines):
            # Insérer le commentaire avant la ligne
            lines.insert(line_num - 1, comment)

            with open(file_path, 'w') as f:
                f.writelines(lines)
            return True
        else:
            print(f"Warning: Line {line_num} out of range for {file_path}")
            return False
    except Exception as e:
        print(f"Error processing {file_path}: {e}")
        return False

def main():
    total_fixes = 0
    for file_path, line_comments in fixes.items():
        print(f"Processing {file_path}...")
        # Trier par ordre décroissant pour ne pas décaler les numéros de ligne
        for line_num in sorted(line_comments.keys(), reverse=True):
            comment = line_comments[line_num]
            if add_comment_before_line(file_path, line_num, comment):
                total_fixes += 1

    print(f"\nTotal: {total_fixes} commentaires ajoutés")
    return 0

if __name__ == "__main__":
    sys.exit(main())

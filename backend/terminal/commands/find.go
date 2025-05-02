package commands

import (
    "fmt"
	"sort"
    "strings"
    global "terminal/global"
)

// matchPattern compara un nombre con un patrón usando * y ?
func matchPattern(name, pattern string) bool {
    // Soporta * y ?
    if pattern == "*" {
        return true
    }
    pi, ni := 0, 0
    for pi < len(pattern) && ni < len(name) {
        if pattern[pi] == '*' {
            if pi+1 == len(pattern) {
                return true
            }
            for k := ni; k <= len(name); k++ {
                if matchPattern(name[k:], pattern[pi+1:]) {
                    return true
                }
            }
            return false
        } else if pattern[pi] == '?' || pattern[pi] == name[ni] {
            pi++
            ni++
        } else {
            return false
        }
    }
    for pi < len(pattern) && pattern[pi] == '*' {
        pi++
    }
    return pi == len(pattern) && ni == len(name)
}

// ParseFind procesa el comando find
func ParseFind(args []string) (string, error) {
    var path, namePattern string
    for i := 0; i < len(args); i++ {
        if strings.HasPrefix(args[i], "-path=") {
            path = strings.Trim(strings.TrimPrefix(args[i], "-path="), "\"")
        } else if strings.HasPrefix(args[i], "-name=") {
            namePattern = strings.Trim(strings.TrimPrefix(args[i], "-name="), "\"")
        }
    }
    if path == "" || namePattern == "" {
        return "", fmt.Errorf("debe especificar los parámetros -path y -name")
    }

    // Unir carpetas y archivos en una sola lista de rutas
    allPaths := append(global.ValidPaths, global.GetValidFilePathsMkfile()...)

    // Filtrar solo los que empiezan con el path base
    var filtered []string
    for _, p := range allPaths {
        if strings.HasPrefix(p, path) {
            filtered = append(filtered, p)
        }
    }

    // Buscar coincidencias por patrón
    var matches []string
    for _, p := range filtered {
        parts := strings.Split(p, "/")
        if len(parts) == 0 {
            continue
        }
        last := parts[len(parts)-1]
        if matchPattern(last, namePattern) {
            matches = append(matches, p)
        }
    }

    // Confirmación de éxito antes del árbol
    confirm := "FIND: Comando ejecutado exitosamente\n"

    // Construir árbol
    tree := buildTree(matches, path)
    return confirm + tree, nil
}

// buildTree construye la representación tipo árbol de las rutas encontradas
func buildTree(paths []string, base string) string {
    if len(paths) == 0 {
        return "# No se encontraron coincidencias"
    }

    // Estructura de árbol
    type Node struct {
        Children map[string]*Node
        IsFile   bool
    }
    root := &Node{Children: make(map[string]*Node)}

    // Construir árbol
    for _, p := range paths {
        rel := strings.TrimPrefix(p, base)
        rel = strings.TrimPrefix(rel, "/")
        if rel == "" {
            continue
        }
        parts := strings.Split(rel, "/")
        curr := root
        for i, part := range parts {
            if part == "" {
                continue
            }
            if curr.Children[part] == nil {
                curr.Children[part] = &Node{Children: make(map[string]*Node)}
            }
            if i == len(parts)-1 && strings.Contains(part, ".") {
                curr.Children[part].IsFile = true
            }
            curr = curr.Children[part]
        }
    }

    // Recursivo para imprimir árbol
    var printTree func(node *Node, prefix string, isLast bool) string
    printTree = func(node *Node, prefix string, isLast bool) string {
        var out string
        keys := make([]string, 0, len(node.Children))
        for k := range node.Children {
            keys = append(keys, k)
        }
        sort.Strings(keys)
        for i, k := range keys {
            child := node.Children[k]
            connector := "├── "
            newPrefix := prefix + "│   "
            if i == len(keys)-1 {
                connector = "└── "
                newPrefix = prefix + "    "
            }
            out += prefix + connector + k + "\n"
            if len(child.Children) > 0 {
                out += printTree(child, newPrefix, i == len(keys)-1)
            }
        }
        return out
    }

    // Imprimir árbol desde la raíz
    tree := base
    if !strings.HasSuffix(tree, "/") {
        tree += "/"
    }
    tree += "\n"
    tree += printTree(root, "", true)
    return tree
}
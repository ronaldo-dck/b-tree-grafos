package main

import (
	"fmt"
	"strings"
)

const t = 2 // Grau (ou ordem) da árvore B
// // No. mínimo de chaves/nó = t - 1
// // No. máximo de chaves/nó = 2*t - 1

type DataType struct {
	nome  string
	index int
}

/******************************
 * Declaração da classe BTree *
 ******************************/

/*******************
 * Estrutura do nó *
 *******************/
type BTreeNode struct {
	leaf     bool         // identifica se um nó é folha (true) ou não (false)
	keys     []DataType   // vetor de chaves em cada nó
	children []*BTreeNode // vetor de ponteiros em cada nó
}

/*********************************************************
 * InitNode(leaf): Criação de um novo nó da árvore B *
 * leaf indica se o novo nó será uma folha ou não        *
 *********************************************************/
func InitNode(leaf bool) *BTreeNode {
	return &BTreeNode{
		leaf:     leaf,
		keys:     []DataType{},
		children: []*BTreeNode{},
	}
}

/**************************************
 * Definição da estrutura da Árvore B *
 **************************************/
type BTree struct {
	root *BTreeNode
}

/*************************************
 * Init(): Inicialização da árvore-B *
 *************************************/
func Init() *BTree {

	return &BTree{
		root: InitNode(true),
	}
}

/*********************************************************
 * Impressão da árvore B em forma de árvore de diretório *
 *********************************************************/
func (node *BTreeNode) Print(indent string, last bool) {
	fmt.Print(indent)
	if last {
		fmt.Print("└─ ")
		indent += "    "
	} else {
		fmt.Print("├─ ")
		indent += "|   "
	}
	keys := make([]string, len(node.keys))
	fmt.Print("[")
	for i, key := range node.keys {
		keys[i] = fmt.Sprintf("%v", key)
	}
	fmt.Println(strings.Join(keys, "|"), "]")

	childCount := len(node.children)
	for i, child := range node.children {
		child.Print(indent, i == childCount-1)
	}
}

/*********************************************************
 * splitChild(i): implementa a Divisão de um filho cheio *
 * i é o ponto onde o nó será dividido                   *
 *********************************************************/
func (node *BTreeNode) splitChild(i int16) {
	child := node.children[i]
	newChild := InitNode(child.leaf)

	// Move as chaves e os filhos para o novo filho
	newChild.keys = append(newChild.keys, child.keys[t:]...)
	child.keys = child.keys[:t]
	if !child.leaf { // divide o nó em dois
		newChild.children = append(newChild.children, child.children[t:]...)
		child.children = child.children[:t]
	}

	// Insere o novo filho no nó
	node.children = append(node.children, nil)
	copy(node.children[i+2:], node.children[i+1:])
	node.children[i+1] = newChild

	noVazio := DataType{"", 0}

	// Move a chave correspondente para cima
	node.keys = append(node.keys, noVazio)
	copy(node.keys[i+1:], node.keys[i:])
	node.keys[i] = child.keys[t-1]
	child.keys = child.keys[:t-1]
}

/***********************************************************
 * Insert(key): Inserção de uma chave em um nó da árvore B *
 * key é a chave que será inserida                        *
 ***********************************************************/
func (node *BTreeNode) Insert(key DataType) {
	if !node.leaf {
		// Encontra o filho apropriado para inserir a chave
		i := len(node.keys) - 1
		for i >= 0 && key.nome < node.keys[i].nome {
			i--
		}

		// Insere a chave no filho apropriado
		if len(node.children[i+1].keys) == 2*t-1 {
			node.splitChild(int16(i) + 1)
			if key.nome > node.keys[i+1].nome {
				i++
			}
		}
		node.children[i+1].Insert(key)
	} else {
		// Insere a chave no nó folha
		noVazio := DataType{"", 0}
		i := len(node.keys) - 1
		node.keys = append(node.keys, noVazio)
		for i >= 0 && key.nome < node.keys[i].nome {
			node.keys[i+1] = node.keys[i]
			i--
		}
		node.keys[i+1] = key
	}
}

/******************************************************
 * Insert(key): Inserção de uma chave na árvore B     *
 * Esta é a função que deve ser chamada para realizar *
 * a inserção. key é a chave a ser inserida.          *
 ******************************************************/
func (tree *BTree) Insert(key DataType) {
	root := tree.root
	if len(root.keys) == 2*t-1 {
		newRoot := InitNode(false)
		newRoot.children = append(newRoot.children, root)
		newRoot.splitChild(0)
		tree.root = newRoot
	}
	tree.root.Insert(key)
}

// Busca de uma chave na árvore B
func (node *BTreeNode) Search(key string) (*DataType, int) {
	i := 0
	for i < len(node.keys) && key > node.keys[i].nome {
		i++
	}

	if i < len(node.keys) && key == node.keys[i].nome {
		return &node.keys[i], i
	} else if node.leaf {
		return nil, -1
	} else {
		return node.children[i].Search(key)
	}
}

// Busca de uma chave na árvore B
func (tree *BTree) Search(key string) (*DataType, int) {
	return tree.root.Search(key)
}

func (node *BTreeNode) Remove(key string) {
	i := 0
	for i < len(node.keys) && key > node.keys[i].nome {
		i++
	}

	if i < len(node.keys) && key == node.keys[i].nome {
		// Caso 1: A chave está presente no nó folha
		if node.leaf {
			copy(node.keys[i:], node.keys[i+1:])
			node.keys = node.keys[:len(node.keys)-1]
		} else {
			// Caso 2: A chave está presente em um nó interno
			// Encontrar o predecessor (maior elemento) da chave na subárvore esquerda
			predecessor := node.children[i]
			for !predecessor.leaf {
				predecessor = predecessor.children[len(predecessor.children)-1]
			}
			node.keys[i] = predecessor.keys[len(predecessor.keys)-1]
			predecessor.Remove(node.keys[i].nome)
		}
	} else {
		// Caso 3: A chave está presente em um nó interno, mas não neste nó em particular
		if node.leaf {
			fmt.Println("A chave não está presente na árvore.")
			return
		}

		// Verificar se o filho em que a chave deve estar tem um número mínimo de chaves
		if len(node.children[i].keys) < t {
			// Caso 3a: O filho não possui chaves suficientes
			child := node.children[i]
			if i > 0 && len(node.children[i-1].keys) >= t {
				// Caso 3a(i): Se o irmão esquerdo do filho tiver pelo menos t chaves, emprestamos uma chave
				leftSibling := node.children[i-1]
				child.keys = append([]DataType{node.keys[i-1]}, child.keys...)
				node.keys[i-1] = leftSibling.keys[len(leftSibling.keys)-1]
				leftSibling.keys = leftSibling.keys[:len(leftSibling.keys)-1]
				if !child.leaf {
					child.children = append([]*BTreeNode{leftSibling.children[len(leftSibling.children)-1]}, child.children...)
					leftSibling.children = leftSibling.children[:len(leftSibling.children)-1]
				}
			} else if i < len(node.children)-1 && len(node.children[i+1].keys) >= t {
				// Caso 3a(ii): Se o irmão direito do filho tiver pelo menos t chaves, emprestamos uma chave
				rightSibling := node.children[i+1]
				child.keys = append(child.keys, node.keys[i])
				node.keys[i] = rightSibling.keys[0]
				rightSibling.keys = rightSibling.keys[1:]
				if !child.leaf {
					child.children = append(child.children, rightSibling.children[0])
					rightSibling.children = rightSibling.children[1:]
				}
			} else {
				// Caso 3b: Fusão com irmão
				if i > 0 {
					// Fusão com o irmão esquerdo
					leftSibling := node.children[i-1]
					leftSibling.keys = append(leftSibling.keys, node.keys[i-1])
					leftSibling.keys = append(leftSibling.keys, child.keys...)
					if !child.leaf {
						leftSibling.children = append(leftSibling.children, child.children...)
					}
					copy(node.keys[i-1:], node.keys[i:])
					copy(node.children[i:], node.children[i+1:])
				} else {
					// Fusão com o irmão direito
					rightSibling := node.children[i+1]
					child.keys = append(child.keys, node.keys[i])
					child.keys = append(child.keys, rightSibling.keys...)
					if !child.leaf {
						child.children = append(child.children, rightSibling.children...)
					}
					copy(node.keys[i:], node.keys[i+1:])
					copy(node.children[i+1:], node.children[i+2:])
				}
				node.keys = node.keys[:len(node.keys)-1]
				node.children = node.children[:len(node.children)-1]
			}
			child.Remove(key)
		} else {
			// Caso 4: A chave está em um filho e o filho tem chaves suficientes
			node.children[i].Remove(key)
		}
	}
}

func (tree *BTree) Remove(key string) {
	tree.root.Remove(key)
}

func main() {
	// Declaração de variáveis
	var (
		op   int
		tree *BTree
	)
	indexAtual := 0

	for {
		// Limpa a tela
		fmt.Println("ÁRVORES-B (B-Trees)")
		fmt.Println("\t Menu Principal")
		fmt.Println("[ 0] Sair")
		fmt.Println("[ 1] Criar árvore")
		fmt.Println("[ 2] Imprimir árvore")
		fmt.Println("[ 3] Inserir elemento")
		fmt.Println("[ 4] Remover elemento")
		fmt.Println("[ 5] Procurar elemento")

		fmt.Print("\nQual a sua opção? >> ")
		fmt.Scan(&op)
		switch op {
		case 0:
			{
				// Clear()
				fmt.Println("Programa Encerrado!\nTecle [ENTER]")
				fmt.Scanln(&op)
				return
			}
		case 1:
			{
				// Clear()
				fmt.Println("Criando árvore...")
				tree = Init()
				fmt.Println("Árvore criada!\n Tecle [ENTER]")

			}
		case 2:
			{
				// Clear()
				fmt.Println("Árvore armazenada:")
				tree.root.Print("", true)
				fmt.Println("\nTecle [ENTER]")

			}
		case 3:
			{

				fmt.Println("Inserindo elemento")
				fmt.Print("Número a inserir >> ")
				x := DataType{"", 0}
				x.index = indexAtual
				indexAtual++
				fmt.Scanln(&x.nome)
				tree.Insert(x)
				fmt.Println("Elemento inserido. Tecle [ENTER]")
			}
		case 4:
			{

				fmt.Println("Removendo elemento")
				fmt.Print("Número a remover >> ")
				nome := ""
				fmt.Scanln(&nome)
				tree.Remove(nome)
				fmt.Println("Elemento Removido. Tecle [ENTER]")

			}
		case 5:
			{

				fmt.Println("Procurar elemento")
				fmt.Print("Número a procurar >> ")
				nome := ""
				fmt.Scanln(&nome)
				a, i := tree.Search(nome)
				if a == nil {
					fmt.Println("Elemento não encontrado. Tecle [ENTER]")

				} else {
					fmt.Println(nome, "está na árvore:", a, " ", i)

				}
			}
		default:
			{
				fmt.Println("Opção Inválida!\nTecle [ENTER]")
				fmt.Scanln(&op)
			}
		}
	}
}

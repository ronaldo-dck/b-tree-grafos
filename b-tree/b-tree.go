package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

const t = 2 // Grau (ou ordem) da árvore B
// // No. mínimo de chaves/nó = t - 1
// // No. máximo de chaves/nó = 2*t - 1

type DataType struct {
	nome  string
	index int
}

type Contato struct {
	Nome     string
	Endereco string
	Telefone string
	Apagado  bool
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

func (nodo *BTreeNode) PercursoEmOrdem(listaIndex *[]int) {
	if nodo == nil {
		return
	}

	for i := 0; i < len(nodo.keys); i++ {
		if nodo.leaf == false {
			nodo.children[i].PercursoEmOrdem(listaIndex)
		}
		fmt.Printf("%s %d\n", nodo.keys[i].nome, nodo.keys[i].index)
		*listaIndex = append(*listaIndex, nodo.keys[i].index)
	}

	// Percorre o último filho (se houver)
	if nodo.leaf == false {
		nodo.children[len(nodo.keys)].PercursoEmOrdem(listaIndex)
	}
}

///////////////////////////////

func setDados() Contato {

	var contato Contato

	scanner := bufio.NewScanner(os.Stdin)

	fmt.Printf("Nome: ")
	scanner.Scan()
	contato.Nome = scanner.Text()

	fmt.Printf("Endereço: ")
	scanner.Scan()
	contato.Endereco = scanner.Text()

	fmt.Printf("Telefone: ")
	scanner.Scan()
	contato.Telefone = scanner.Text()

	contato.Apagado = false

	return contato
}

func appendDados(contato Contato, file *os.File) {
	fmt.Fprintf(file, "%v;%v;%v;%v\n", contato.Nome, contato.Endereco, contato.Telefone, contato.Apagado)
}

func appendIndex(nodo DataType, file *os.File) {
	fmt.Fprintf(file, "%v;%v\n", nodo.index, nodo.nome)
}

func initTreeFromFile(fileName string, tree *BTree) int {

	file, err := os.Open(fileName)
	if err != nil {
		fmt.Println("Erro ao abrir ou criar o arquivo:", err)
		return 0
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	maxIndex := 0

	for scanner.Scan() {
		line := scanner.Text()
		parts := strings.Split(line, ";")

		if len(parts) != 2 {
			// Ignorar linhas mal formatadas (sem ; ou com mais de um ;)
			continue
		}

		index, err := strconv.Atoi(parts[0])

		if err != nil {
			// Ignorar linhas em que o índice não é um número válido
			continue
		}

		if index > maxIndex {
			maxIndex = index
		}

		nodo := DataType{parts[1], index}
		tree.Insert(nodo)
	}

	return maxIndex + 1
}

func loadIndex() []DataType {

	file, err := os.Open("indexFile.txt")
	if err != nil {
		panic("Erro ao abrir o arquivo")
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	listaIndex := make([]DataType, 0)

	for scanner.Scan() {
		line := scanner.Text()
		parts := strings.Split(line, ";")

		if len(parts) != 2 {
			// Ignorar linhas mal formatadas (sem ; ou com mais de um ;)
			continue
		}

		index, err := strconv.Atoi(parts[0])

		if err != nil {
			// Ignorar linhas em que o índice não é um número válido
			continue
		}

		nodo := DataType{parts[1], index}
		listaIndex = append(listaIndex, nodo)
	}

	return listaIndex
}

func listarContatos(tree *BTree, lixeira bool) {
	lista := make([]int, 0)
	tree.root.PercursoEmOrdem(&lista)

	agenda := loadData()

	for _, v := range lista {
		if lixeira == true && agenda[v].Apagado == true {
			fmt.Println(agenda[v])
		} else if lixeira == false && agenda[v].Apagado == false {
			fmt.Println(agenda[v])
		}
	}
}

func loadData() []Contato {

	file, err := os.Open("dataFile.txt")
	if err != nil {
		panic("Erro ao abrir o arquivo")
	}
	defer file.Close()

	agenda := make([]Contato, 0)

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		parts := strings.Split(line, ";")

		if len(parts) != 4 {
			// Ignorar linhas mal formatadas (sem ; ou com mais de um ;)
			continue
		}
		apagado := false
		if parts[3] == "true" {
			apagado = true
		}

		contato := Contato{parts[0], parts[1], parts[2], apagado}
		agenda = append(agenda, contato)
	}

	return agenda
}

func saveData(agenda []Contato) {
	file, err := os.OpenFile("dataFile.txt", os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0666)
	if err != nil {
		panic("Erro ao abrir o arquivo")
	}
	defer file.Close()

	for _, contato := range agenda {
		appendDados(contato, file)
	}
}

func saveIndex(listaIndex []DataType) {
	file, err := os.OpenFile("indexFile.txt", os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0666)
	if err != nil {
		panic("Erro ao abrir o arquivo")
	}
	defer file.Close()

	for _, contato := range listaIndex {
		appendIndex(contato, file)
	}
}

func enviarParaLixeira(nome string, tree *BTree) {
	contato, err := tree.Search(nome)
	if err == -1 {
		fmt.Println("Contato não encontrado")
		return
	}

	agenda := loadData()

	aux := agenda[contato.index]
	aux.Apagado = true
	agenda[contato.index] = aux

	saveData(agenda)
}

func restaurarDaLixeira(tree *BTree) {
	listarContatos(tree, true)
	scanner := bufio.NewScanner(os.Stdin)
	fmt.Printf("Nome: ")
	scanner.Scan()

	contato, err := tree.Search(scanner.Text())
	if err == -1 {
		fmt.Println("Contato não encontrado")
		return
	}

	agenda := loadData()

	aux := agenda[contato.index]
	aux.Apagado = false
	agenda[contato.index] = aux

	saveData(agenda)
}

func esvaziarLixeira(tree *BTree) {

	lista := make([]int, 0)
	tree.root.PercursoEmOrdem(&lista)
	agenda := loadData()
	apagados := make([]int, 0)

	for _, v := range lista {
		if agenda[v].Apagado == true {
			apagados = append(apagados, v)
			tree.Remove(agenda[v].Nome)
		}
	}

	listaIndex := loadIndex()
	newAgenda := make([]Contato, 0)
	newIndex := make([]DataType, 0)
	fmt.Println(listaIndex)

	for i := range listaIndex {
		achou := false
		for _, p := range apagados {
			if listaIndex[i].index == p {
				achou = true
				break
			}
		}
		if achou == false {
			newIndex = append(newIndex, listaIndex[i])
			newAgenda = append(newAgenda, agenda[i])
		}
	}

	tree.root.Print("", true)
	fmt.Println(newIndex)
	fmt.Printf("atualizando index:\n")
	for i, v := range newIndex {
		fmt.Printf("buscando por: %v\n", v.nome)
		newIndex[i].index = i
		nodo, err := tree.Search(v.nome)
		println(nodo.nome, nodo.index, err)
		nodo.index = i
	}

	saveData(newAgenda)
	saveIndex(newIndex)

}

func main() {
	var (
		tree       *BTree
		indexAtual int
		op         int
	)
	tree = Init()
	indexAtual = initTreeFromFile("indexFile.txt", tree)

	// indexFile, err := os.OpenFile("indexFile.txt", os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0666)
	// if err != nil {
	// 	fmt.Println("Erro ao abrir ou criar o arquivo:", err)
	// 	return
	// }
	// defer indexFile.Close()

	// dataFile, err := os.OpenFile("dataFile.txt", os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0666)
	// if err != nil {
	// 	fmt.Println("Erro ao abrir ou criar o arquivo:", err)
	// 	return
	// }
	// defer dataFile.Close()

	println(indexAtual)
	tree.root.Print("", true)

	for {
		fmt.Println("0 - Sair")
		fmt.Println("1 - Inserir Contato")
		fmt.Println("2 - Listar Contatos")
		fmt.Println("3 - Enviar para lixeira")
		fmt.Println("4 - Restaurar da lixeira")
		fmt.Println("5 - Esvaziar lixeira")

		fmt.Scan(&op)
		switch op {
		case 0:
			return
		case 1:
			contato := setDados()
			// appendDados(contato, dataFile)
			nodo := DataType{contato.Nome, indexAtual}
			tree.Insert(nodo)
			// appendIndex(nodo, indexFile)
			indexAtual++
		case 2:
			listarContatos(tree, false)
		case 3:
			scanner := bufio.NewScanner(os.Stdin)
			fmt.Printf("Nome: ")
			scanner.Scan()
			enviarParaLixeira(scanner.Text(), tree)
		case 4:
			restaurarDaLixeira(tree)
		case 5:
			esvaziarLixeira(tree)
		case 9:
			tree.root.Print("", true)
		}

	}

	/* // Inserção
	// salval file;
	//	constriu o dados Datatype

	tree = Init()
	// insere na arvore;

	tree.root.Print("", true) */

}

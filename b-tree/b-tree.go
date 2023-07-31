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
				irmaoEsquerdo := node.children[i-1]
				child.keys = append([]DataType{node.keys[i-1]}, child.keys...)
				node.keys[i-1] = irmaoEsquerdo.keys[len(irmaoEsquerdo.keys)-1]
				irmaoEsquerdo.keys = irmaoEsquerdo.keys[:len(irmaoEsquerdo.keys)-1]
				if !child.leaf {
					child.children = append([]*BTreeNode{irmaoEsquerdo.children[len(irmaoEsquerdo.children)-1]}, child.children...)
					irmaoEsquerdo.children = irmaoEsquerdo.children[:len(irmaoEsquerdo.children)-1]
				}
			} else if i < len(node.children)-1 && len(node.children[i+1].keys) >= t {
				// Caso 3a(ii): Se o irmão direito do filho tiver pelo menos t chaves, emprestamos uma chave
				irmaoDireito := node.children[i+1]
				child.keys = append(child.keys, node.keys[i])
				node.keys[i] = irmaoDireito.keys[0]
				irmaoDireito.keys = irmaoDireito.keys[1:]
				if !child.leaf {
					child.children = append(child.children, irmaoDireito.children[0])
					irmaoDireito.children = irmaoDireito.children[1:]
				}
			} else {
				// Caso 3b: Fusão com irmão
				if i > 0 {
					// Fusão com o irmão esquerdo
					irmaoEsquerdo := node.children[i-1]
					irmaoEsquerdo.keys = append(irmaoEsquerdo.keys, node.keys[i-1])
					irmaoEsquerdo.keys = append(irmaoEsquerdo.keys, child.keys...)
					if !child.leaf {
						irmaoEsquerdo.children = append(irmaoEsquerdo.children, child.children...)
					}
					copy(node.keys[i-1:], node.keys[i:])
					copy(node.children[i:], node.children[i+1:])
				} else {
					// Fusão com o irmão direito
					irmaoDireito := node.children[i+1]
					child.keys = append(child.keys, node.keys[i])
					child.keys = append(child.keys, irmaoDireito.keys...)
					if !child.leaf {
						child.children = append(child.children, irmaoDireito.children...)
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
		// fmt.Printf("%s %d\n", nodo.keys[i].nome, nodo.keys[i].index)
		*listaIndex = append(*listaIndex, nodo.keys[i].index)
	}

	// Percorre o último filho (se houver)
	if nodo.leaf == false {
		nodo.children[len(nodo.keys)].PercursoEmOrdem(listaIndex)
	}
}

///////////////////////////////

func limString(scanner *bufio.Scanner, info string, limite int) string {

	var res string

	fmt.Print(info)
	scanner.Scan()
	res = scanner.Text()
	for res == "" || len(res) > limite {
		fmt.Printf("Dados inválidos, insira novamente\n")
		fmt.Print(info)
		scanner.Scan()
		res = scanner.Text()
	}

	return res

}

func setDados(tree *BTree) Contato {

	fmt.Printf("Dados para o contato:\n")
	var contato Contato

	scanner := bufio.NewScanner(os.Stdin)

	err := 0
	for err != -1 {
		contato.Nome = limString(scanner, "Nome: ", 30)
		_, err = tree.Search(contato.Nome)
		if err != -1 {
			fmt.Println("Nome já existente, insira outro")
		}
	}
	contato.Endereco = limString(scanner, "Endereço: ", 50)
	contato.Telefone = limString(scanner, "Telefone: ", 15)

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

		os.Create("indexFile.txt")
		os.Create("dataFile.txt")

		fmt.Println("Arquivos de data e index Criados")
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

	if maxIndex == 0 {
		return 0
	} else {
		return maxIndex + 1
	}
}

func getIndex() int {
	file, err := os.Open("indexFile.txt")
	if err != nil {

		os.Create("indexFile.txt")
		os.Create("dataFile.txt")

		fmt.Println("Arquivos de data e index Criados")
		return 0
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	maxIndex := 0

	if scanner.Scan(){
		maxIndex++
	} else {
		return 0
	}

	for scanner.Scan() {
		maxIndex++
	}

	return maxIndex
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

	if len(lista) == 0 {
		fmt.Println("Não há contatos\nPressione enter")
		fmt.Scanln()
		return
	}
	fmt.Println("\nNome | Endereço | Telefone:")
	for _, v := range lista {
		if lixeira == true && agenda[v].Apagado == true {
			fmt.Printf("%s | %s | %s\n", agenda[v].Nome, agenda[v].Endereco, agenda[v].Telefone)
		} else if lixeira == false && agenda[v].Apagado == false {
			fmt.Printf("%s | %s | %s\n", agenda[v].Nome, agenda[v].Endereco, agenda[v].Telefone)
		}
	}

	fmt.Printf("Pressione enter para Continuar: ")
	fmt.Scanln()
	Clear()

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
	fmt.Println("Contato enviado para lixeira")
}

func restaurarDaLixeira(tree *BTree) {
	fmt.Printf("Lista de contatos na lixeira:\n")
	listarContatos(tree, true)
	scanner := bufio.NewScanner(os.Stdin)
	fmt.Printf("Nome para restaurar da lixeira: ")
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
	fmt.Println("Contato restaurado da lixeira.")
}

func esvaziarLixeira(tree *BTree) *BTree {

	lista := make([]int, 0)
	tree.root.PercursoEmOrdem(&lista)
	agenda := loadData()
	apagados := make([]int, 0)

	apagaTudo := true
	for _, v := range lista {
		if agenda[v].Apagado == false {
			apagaTudo = false
		}
	}

	if apagaTudo {
		os.Create("indexFile.txt")
		os.Create("dataFile.txt")
		var newTree *BTree
		newTree = Init()
		return newTree
	}

	for _, v := range lista {
		if agenda[v].Apagado == true {
			apagados = append(apagados, v)
			tree.Remove(agenda[v].Nome)
		}
	}

	listaIndex := loadIndex()
	newAgenda := make([]Contato, 0)
	newIndex := make([]DataType, 0)

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

	for i, v := range newIndex {
		newIndex[i].index = i
		nodo, _ := tree.Search(v.nome)
		nodo.index = i
	}

	saveData(newAgenda)
	saveIndex(newIndex)

	return tree
}

func editarContato(tree *BTree) *BTree {
	scanner := bufio.NewScanner(os.Stdin)
	fmt.Printf("Nome do contato a ser editado: ")
	scanner.Scan()
	nodo, err := tree.Search(scanner.Text())
	
	listaDados := loadData()
	if err == -1 || listaDados[nodo.index].Apagado == true {
		fmt.Println("Contato não encontrado\nPressione enter")
		fmt.Scanln()
		return tree
	}

	novosDados := setDados(tree)
	listaDados[nodo.index] = novosDados
	saveData(listaDados)

	if novosDados.Nome == nodo.nome {
		return tree
	}

	novoNodo := DataType{nome: novosDados.Nome, index: nodo.index}
	listaIndex := loadIndex()
	listaIndex[novoNodo.index] = novoNodo
	saveIndex(listaIndex)

	newTree := Init()
	initTreeFromFile("indexFile.txt", newTree)

	fmt.Println("Contato editado.")
	return newTree
}

func Clear() {
	fmt.Print("\033[H\033[2J") // escape codes para limpar a tela (Unix)
}

func main() {
	var (
		tree       *BTree
		indexAtual int
		op         int
	)
	tree = Init()
	indexAtual = initTreeFromFile("indexFile.txt", tree)

	for {
		fmt.Println("Menu:")
		fmt.Println(" 0 - Sair")
		fmt.Println(" 1 - Inserir Contato")
		fmt.Println(" 2 - Listar Contatos")
		fmt.Println(" 3 - Enviar para lixeira")
		fmt.Println(" 4 - Restaurar da lixeira")
		fmt.Println(" 5 - Esvaziar lixeira")
		fmt.Println(" 6 - Editar contato")

		fmt.Print("Operação desejada: ")
		fmt.Scanf("%d\n", &op)
		switch op {
		case 0:
			return
		case 1:
			listaIndex := loadIndex()
			listaData := loadData()
			indexAtual = getIndex()

			contato := setDados(tree)
			listaData = append(listaData, contato)

			nodo := DataType{contato.Nome, indexAtual}
			listaIndex = append(listaIndex, nodo)
			tree.Insert(nodo)

			indexAtual++
			saveIndex(listaIndex)
			saveData(listaData)
			fmt.Printf("Contato inserido.\nPressione enter")
			fmt.Scanln()
			Clear()
		case 2:
			Clear()
			listarContatos(tree, false)
		case 3:
			scanner := bufio.NewScanner(os.Stdin)
			fmt.Printf("Nome a mandar para lixeira: ")
			scanner.Scan()
			enviarParaLixeira(scanner.Text(), tree)

			fmt.Printf("Pressione enter")
			fmt.Scanln()
			Clear()
		case 4:
			restaurarDaLixeira(tree)
			fmt.Printf("Pressione enter")
			fmt.Scanln()
			Clear()
		case 5:
			fmt.Printf("Todos os contatos na lixeira serão excluídos, deseja continuar? [s/N] ")
			var confirma rune
			fmt.Scanf("%c\n", &confirma)
			if confirma == 's' {
				tree = esvaziarLixeira(tree)
				fmt.Printf("Lixeira esvaziada.\nPressione enter")
			} else {
				fmt.Printf("Operação cancelada.\nPressione enter")
			}

			fmt.Scanln()
			Clear()
		case 6:
			tree = editarContato(tree)
		}

	}

}

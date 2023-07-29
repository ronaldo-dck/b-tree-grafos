#include <iostream>
#include <vector>
#include <iomanip>
#include <fstream>
#include <algorithm>
#include <cmath>
#include <queue>
#include <sstream>
#include <limits>
using namespace std;

const int INF = numeric_limits<int>::max();

class Vertice
{
public:
    float x;
    float y;
    string label;
    Vertice()
    {
        x = 0;
        y = 0;
        label = "";
    }
    Vertice(float x, float y, string label)
    {
        this->x = x;
        this->y = y;
        this->label = label;
    }
    Vertice(const Vertice &other) : x(other.x), y(other.y), label(other.label) {}
};

class Aresta
{
public:
    float custo;
    string label;
    string v1, v2;
    Aresta() {}
    Aresta(float custo, string label, string v1, string v2)
    {
        this->custo = custo;
        this->label = label;
        this->v1 = v1;
        this->v2 = v2;
    }

    Aresta(const Aresta &other)
        : custo(other.custo), label(other.label), v1(other.v1), v2(other.v2) {}
};

class Fleury
{
private:
    vector<vector<float>> grafo;

    bool isPonte(int u, int v)
    {
        int grau = 0;
        for (int i = 0; i < grafo.size(); i++)
        {
            if (grafo[v][i] != 0.0)
            {
                grau++;
                if (grau > 1)
                {
                    return false;
                }
            }
        }
        return true;
    }

    int countArestas()
    {
        int count = 0;
        for (int i = 0; i < grafo.size(); i++)
        {
            for (int j = i; j < grafo.size(); j++)
            {
                if (grafo[i][j] != 0.0 || grafo[j][i] != 0.0)
                {
                    count++;
                }
            }
        }
        return count;
    }

public:
    Fleury(const vector<vector<float>> &matrizCustos)
    {
        grafo = matrizCustos;
    }

    void CicloEuleriano(int verticeInicial, vector<int> &ciclo)
    {
        int numeroDeArestas = countArestas();
        for (int i = 0; i < grafo.size(); i++)
        {
            if (grafo[verticeInicial][i] != 0.0)
            {
                if (numeroDeArestas <= 1 || !isPonte(verticeInicial, i))
                {
                    ciclo.push_back(i);
                    grafo[verticeInicial][i] = 0.0;
                    grafo[i][verticeInicial] = 0.0;
                    numeroDeArestas--;
                    CicloEuleriano(i, ciclo);
                }
            }
        }
    }
};

class Grafo
{
private:
    vector<Vertice> vertices;
    vector<Aresta> arestas;
    vector<vector<float>> custos;

    void geraMatrizCustos()
    {
        vector<vector<float>> aux(size(vertices), vector<float>(size(vertices)));

        for (int i = 0; i < size(arestas); i++)
        {
            int linha, coluna;

            for (int j = 0; j < size(vertices); j++)
            {
                if (arestas[i].v1 == vertices[j].label)
                {
                    linha = j;
                }
                else if (arestas[i].v2 == vertices[j].label)
                {
                    coluna = j;
                }
            }

            aux[linha][coluna] = arestas[i].custo;
            aux[coluna][linha] = arestas[i].custo;
        }

        custos = aux;
    }

    void dfs(int node, vector<bool> &visitados, bool logico)
    {
        visitados[node] = true;
        if (!logico)
            cout << vertices[node].label;

        for (int vizinho = 0; vizinho < size(custos); ++vizinho)
        {
            if (custos[node][vizinho] != 0 && !visitados[vizinho])
            {
                if (!logico)
                    cout << "->";

                dfs(vizinho, visitados, logico);
            }
        }
    }

    void bfs(int startNode, vector<bool> &visitados)
    {
        int numNodes = size(vertices);
        queue<int> q;

        visitados[startNode] = true;
        q.push(startNode);

        while (!q.empty())
        {
            int nodoAtual = q.front();
            q.pop();
            cout << vertices[nodoAtual].label << "->";

            for (int vizinho = 0; vizinho < numNodes; ++vizinho)
            {
                if (custos[nodoAtual][vizinho] != 0 && !visitados[vizinho])
                {
                    visitados[vizinho] = true;
                    q.push(vizinho);
                }
            }
        }
    }

    vector<int> __dijkstra(const vector<vector<float>> &grafo, int origem, int destino)
    {
        int numVertices = grafo.size();
        vector<float> distancias(numVertices, INF);
        vector<int> caminhoAnterior(numVertices, -1);
        distancias[origem] = 0;

        queue<int> fila;
        fila.push(origem);

        while (!fila.empty())
        {
            int u = fila.front();
            fila.pop();

            for (int v = 0; v < numVertices; v++)
            {
                if (grafo[u][v] != 0.0 && distancias[u] + grafo[u][v] < distancias[v])
                {
                    distancias[v] = distancias[u] + grafo[u][v];
                    caminhoAnterior[v] = u;
                    fila.push(v);
                }
            }
        }

        vector<int> caminho;
        int atual = destino;
        while (atual != origem && atual != -1)
        {
            caminho.push_back(atual);
            atual = caminhoAnterior[atual];
        }
        caminho.push_back(origem);
        reverse(caminho.begin(), caminho.end());

        return caminho;
    }

public:
    Grafo(){};
    Grafo(string nameFile)
    {
        ifstream file(nameFile, ios::in);
        string linhaN;
        getline(file, linhaN);
        int n = stoi(linhaN);

        for (int i = 0; i < n; i++)
        {
            string linha, label, linhaX, linhaY;
            float x, y;

            getline(file, linha);
            stringstream parser(linha);

            getline(parser, label, ';');
            getline(parser, linhaX, ';');
            getline(parser, linhaY, ';');

            InserirVertice(stoi(linhaX), stoi(linhaY), label);
        }

        getline(file, linhaN);
        n = stoi(linhaN);
        for (int i = 0; i < n; i++)
        {
            string linha, label, linhaV1, linhaV2;

            getline(file, linha);
            stringstream parser(linha);

            getline(parser, label, ';');
            getline(parser, linhaV1, ';');
            getline(parser, linhaV2, ';');

            InserirAresta(label, linhaV1, linhaV2);
        }

        file.close();
    };

    // Construtor de cópia
    Grafo(const Grafo &other)
    {
        // Copiar os vértices
        for (const Vertice &v : other.vertices)
        {
            vertices.push_back(v); // Supondo que a classe Vertice possui um construtor de cópia apropriado
        }

        // Copiar as arestas
        for (const Aresta &a : other.arestas)
        {
            arestas.push_back(a); // Supondo que a classe Aresta possui um construtor de cópia apropriado
        }

        // Copiar a matriz de custos
        int n = other.custos.size();
        custos.resize(n, vector<float>(n, 0.0f)); // Redimensiona a matriz de custos
        for (int i = 0; i < n; ++i)
        {
            for (int j = 0; j < n; ++j)
            {
                custos[i][j] = other.custos[i][j];
            }
        }
    }

    void Dijkstra()
    {
        string v1, v2;
        cout << "Vértice de Origen: ";
        getline(cin, v1);
        cout << "Vértice de Destino: ";
        getline(cin, v2);

        int i = 0, origin, destino;
        while (v1 != vertices[i++].label)
        {
        }
        origin = --i;
        i = 0;
        while (v2 != vertices[i++].label)
        {
        }
        destino = --i;

        vector<int> caminho = __dijkstra(custos, origin, destino);

        cout << vertices[caminho[0]].label;
        for (int i = 1; i < size(caminho); i++)
        {
            cout << "->" << vertices[caminho[i]].label;
        }
        cout << endl;
    }

    void SalvarGrafo(string nameFile)
    {
        ofstream file(nameFile, ios::trunc | ios::out);

        file << size(vertices) << endl;
        for (int i = 0; i < size(vertices); i++)
        {
            file << vertices[i].label << " " << vertices[i].x << " " << vertices[i].y << endl;
        }
        file << size(arestas) << endl;
        for (int i = 0; i < size(arestas); i++)
        {
            file << arestas[i].label << " " << arestas[i].v1 << " " << arestas[i].v2 << endl;
        }

        file.close();
    }

    void InserirVertice(float x, float y, string label)
    {
        Vertice novoVertice(x, y, label);
        vertices.push_back(novoVertice);
        geraMatrizCustos();
    };

    void InserirAresta(string label, string v1, string v2)
    {

        Vertice A, B;

        for (int i = 0; i < size(vertices); i++)
        {
            if (v1 == vertices[i].label)
            {
                A = vertices[i];
            }
            else if (v2 == vertices[i].label)
            {
                B = vertices[i];
            }
        }
        if (A.label == "" || B.label == "")
        {
            cout << "Vértices não encontrados.\n";
            return;
        }

        float custo = sqrt(pow(A.x - B.x, 2) + pow(A.y - B.y, 2));

        Aresta novaAresta(custo, label, v1, v2);
        arestas.push_back(novaAresta);
        geraMatrizCustos();
        // cout << "Aresta inserida.\n";
    }

    void RemoverAresta(string label)
    {
        int index = -1;
        for (int i = 0; i < size(arestas); i++)
        {
            if (arestas[i].label == label)
            {
                index = i;
                break;
            }
        }

        if (index == -1)
        {
            cout << "Aresta não encontrada.\n";
            return;
        }

        arestas.erase(arestas.begin() + index);
        geraMatrizCustos();
        cout << "Aresta '" << label << "' removida\n";
    }

    void RemoverVertice(string label)
    {
        int index = -1;
        for (int i = 0; i < size(vertices); i++)
        {
            if (vertices[i].label == label)
            {
                index = i;
                break;
            }
        }

        if (index == -1)
        {
            cout << "Vértice não encontrado.\n";
            return;
        }

        vertices.erase(vertices.begin() + index);

        vector<Aresta> tempAresta;
        for (int i = 0; i < size(arestas); i++)
        {
            if (arestas[i].v1 != label && arestas[i].v2 != label)
            {
                tempAresta.push_back(arestas[i]);
            }
        }
        arestas = tempAresta;
        geraMatrizCustos();
        cout << "Vértice '" << label << "' removido\n";
    }

    void depthFirstSearch(string label)
    {
        vector<bool> visitados(size(vertices), false);

        for (int i = 0; i < size(vertices); i++)
        {
            if (label == vertices[i].label)
            {
                dfs(i, visitados, false);
                cout << endl;
                return;
            }
        }

        cout << "Vértice não encontrado.\n";
    };

    void breadthFirstSearch(string label)
    {
        vector<bool> visitados(size(vertices), false);
        for (int i = 0; i < size(vertices); i++)
        {
            if (label == vertices[i].label)
            {
                bfs(i, visitados);
                cout << endl;
                return;
            }
        }
        cout << "Vértice não encontrado.\n";
    }

    int GoodMan()
    {
        vector<bool> visitados(size(vertices), false);
        int componentes = 0;

        for (int i = 0; i < size(vertices); i++)
        {
            if (!visitados[i])
            {
                dfs(i, visitados, true);
                componentes++;
            }
        }

        return componentes;
    }

    void PrintInfos()
    {
        cout << "lista de arestas" << endl;
        for (int i = 0; i < size(arestas); i++)
        {
            cout << arestas[i].label << ": " << arestas[i].v1 << "->" << arestas[i].v2 << " custo: " << arestas[i].custo << endl;
        }

        cout << "lista de vertices" << endl;
        for (int i = 0; i < size(vertices); i++)
        {
            cout << vertices[i].label << ": (" << vertices[i].x << ", " << vertices[i].y << ")\n";
        }

        cout << "Matriz de custos" << endl
             << "     ";
        for (int i = 0; i < size(custos); i++)
        {
            cout << " '" << vertices[i].label << "'";
        }
        cout << endl;
        for (int i = 0; i < size(custos); i++)
        {
            cout << "'" << vertices[i].label << "' :";
            for (int j = 0; j < size(custos); j++)
            {
                cout << fixed << setprecision(2) << custos[i][j] << " ";
            }
            cout << endl;
        }
    }

    bool IsEuleriano()
    {
        for (int i = 0; i < size(custos); i++)
        {
            int contador = 0;
            for (int j = 0; j < size(custos); j++)
            {
                if (custos[i][j] != 0)
                {
                    contador++;
                }
            }
            if (contador % 2)
                return false;
        }

        return true;
    }

    void FleuryCiclo()
    {
        if (!this->IsEuleriano())
        {
            cout << "O grafo não é euleriano.\n";
            return;
        }
        Fleury ciclo(custos);

        string label;
        cout << "Vértice de Inicio: ";
        cin >> label;

        int origin = 0;
        while (label != vertices[origin].label)
        {
            origin++;
        }

        vector<int> cicloEuler;
        ciclo.CicloEuleriano(origin, cicloEuler);

        cout << vertices[origin].label;
        for (int i = 0; i < size(cicloEuler); i++)
        {
            cout << "->" << vertices[cicloEuler[i]].label;
        }
        cout << endl;
    }
};

int main()
{
    Grafo G("grafoIn.txt");
    int op;
    string label, A, B;
    bool isEuler;

    while (1)
    {
        cout << " 0 - Sair\n"
             << " 1 - Inserir Vértice\n"
             << " 2 - Inserir Aresta\n"
             << " 3 - Remover Vértice\n"
             << " 4 - Remover Aresta\n"
             << " 5 - Salvar em arquivo\n"
             << " 6 - Verifica Euleriano\n"
             << " 7 - Componentes conexos\n"
             << " 8 - Busca em profundidade\n"
             << " 9 - Busca em largura\n"
             << " 10 - Dijkstra\n"
             << " 11 - Fleury\n";
        cin >> op;
        cin.ignore();
        switch (op)
        {
        case 0:
            return 0;
            break;
        case 1:
            float x, y;
            cout << "Label: ";
            getline(cin, label);
            cout << "X: ";
            cin >> x;
            cout << "Y: ";
            cin >> y;

            G.InserirVertice(x, y, label);
            cout << "Vértice inserido.\n";
            break;
        case 2:
            cout << "Label: ";
            getline(cin, label);
            cout << "Vertice A: ";
            getline(cin, A);
            cout << "Vertice B: ";
            getline(cin, B);

            G.InserirAresta(label, A, B);
            break;
        case 3:
            cout << "Vertice a ser removido: ";
            getline(cin, label);
            G.RemoverVertice(label);
            break;
        case 4:
            cout << "Aresta a ser removida: ";
            getline(cin, label);
            G.RemoverAresta(label);
            break;
        case 5:
            cout << "Nome do arquivo de saída: ";
            getline(cin, label);
            G.SalvarGrafo(label);
            break;
        case 6:
            isEuler = G.IsEuleriano();
            if (isEuler)
                cout << "O grafo é Euleriano!" << endl;
            else
                cout << "O grafo NÃO é Euleriano!" << endl;
            break;
        case 7:
            cout << "Número de componentes conexos : " << G.GoodMan() << endl;
            break;

        case 8:

            cout << "Informe o vértice de origem: ";
            getline(cin, label);
            G.depthFirstSearch(label);
            break;
        case 9:

            cout << "Informe o vértice de origem: ";
            getline(cin, label);
            G.breadthFirstSearch(label);
            break;
        case 10:
            G.Dijkstra();
            break;
        case 11:
            G.FleuryCiclo();
            break;
        case 12:
            G.PrintInfos();
            break;
        default:
            break;
        }
    }

    return 0;
}
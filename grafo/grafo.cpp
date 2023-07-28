#include <iostream>
#include <list>
#include <vector>
#include <iomanip>
#include <fstream>
#include <algorithm>
#include <cmath>
#include <queue>
using namespace std;

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
};

class Aresta
{
public:
    float custo;
    string label;
    string v1, v2;
    Aresta(float custo, string label, string v1, string v2)
    {
        this->custo = custo;
        this->label = label;
        this->v1 = v1;
        this->v2 = v2;
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

    void dfs(int node, vector<bool> &visitados)
    {
        visitados[node] = true;
        cout << vertices[node].label << "->";

        for (int vizinho = 0; vizinho < size(custos); ++vizinho)
        {
            if (custos[node][vizinho] != 0 && !visitados[vizinho])
            {
                dfs(vizinho, visitados);
            }
        }
    }

    void bfs(int startNode, vector<bool> &visitados)
    {
        int numNodes = size(vertices);
        std::queue<int> q;

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

public:
    Grafo(/* args */){};
    Grafo(string nameFile)
    {
        ifstream file(nameFile, ios::in);
        int n;
        file >> n;
        string label;
        float x, y;
        for (int i = 0; i < n; i++)
        {
            file >> label >> x >> y;
            InserirVertice(x, y, label);
        }

        file >> n;
        string v1, v2;
        for (int i = 0; i < n; i++)
        {
            file >> label >> v1 >> v2;
            InserirAresta(label, v1, v2);
        }

        file.close();
    };

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
            cout << "Vértices não encontrados";
        }

        float custo = sqrt(pow(A.x - B.x, 2) + pow(A.y - B.y, 2));

        Aresta novaAresta(custo, label, v1, v2);
        arestas.push_back(novaAresta);
        geraMatrizCustos();
    }

    void RemoverAresta(string label)
    {
        int index;
        for (int i = 0; i < size(arestas); i++)
        {
            if (arestas[i].label == label)
            {
                index = i;
                break;
            }
        }

        arestas.erase(arestas.begin() + index);
        geraMatrizCustos();
    }

    void RemoverVertice(string label)
    {
        int index;
        for (int i = 0; i < size(vertices); i++)
        {
            if (vertices[i].label == label)
            {
                index = i;
                break;
            }
        }

        vertices.erase(vertices.begin() + index);
        for (int i = 0; i < size(arestas); i++)
        {
            if (arestas[i].v1 == label || arestas[i].v2 == label)
            {
                arestas.erase(arestas.begin() + i);
            }
        }
        geraMatrizCustos();
    }

    void depthFirstSearch(string label)
    {
        int index;
        vector<bool> visitados(size(vertices), false);

        for (int i = 0; i < size(vertices); i++)
        {
            if (label == vertices[i].label)
            {
                dfs(i, visitados);
            }
        }
    };

    void breadthFirstSearch()
    {
        string label;
        cout << "Informe o vértice de origem: ";
        cin >> label;
        vector<bool> visitados(size(vertices), false);
        for (int i = 0; i < size(vertices); i++)
        {
            if (label == vertices[i].label)
            {
                bfs(i, visitados);
                return;
            }
        }
    }

    int GoodMan()
    {
        return 0;
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
            cout << " " << vertices[i].label;
        }
        cout << endl;
        for (int i = 0; i < size(custos); i++)
        {
            cout << vertices[i].label << " :";
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

    int compsConexos()
    {
        vector<int> listaComps(size(vertices));

        for (int i = 1; i < size(vertices); i++)
        {
            for (int j = 0; j < i; j++)
            {
                if (custos[i][j] != 0)
                {
                    listaComps[i] = 1;
                    listaComps[j] = 1;
                }
            }
        }

        int componentes = 1;
        for (int i = 0; i < size(listaComps); i++)
        {
            cout << listaComps[i] << " ";
            if (listaComps[i] == 0)
                componentes++;
        }

        return componentes;
    }
};

int main()
{
    Grafo formiga("grafoIn.txt");
    int op;
    string label, A, B;
    bool isEuler;

    while (1)
    {
        cout << " 0 - Sair\n 1 - Inserir Vértice\n 2 - Inserir Aresta\n 3 - Remover Vértice\n 4 - Remover Aresta\n 5 - Salvar em arquivo\n 6 - Verifica Euleriano\n 7 - Componentes desconexos\n";
        cin >> op;
        switch (op)
        {
        case 0:
            return 0;
            break;
        case 1:
            float x, y;
            cout << "Label: ";
            cin >> label;
            cout << "X: ";
            cin >> x;
            cout << "Y: ";
            cin >> y;

            formiga.InserirVertice(x, y, label);
            break;
        case 2:
            cout << "Label: ";
            cin >> label;
            cout << "Vertice A: ";
            cin >> A;
            cout << "Vertice B: ";
            cin >> B;

            formiga.InserirAresta(label, A, B);
            break;
        case 3:
            cout << "Label: ";
            cin >> label;
            formiga.RemoverVertice(label);
            break;
        case 4:
            cout << "Label: ";
            cin >> label;
            formiga.RemoverAresta(label);
            break;
        case 5:
            // cout << "Nome do arquivo: ";
            // cin >> label;
            formiga.SalvarGrafo("grafoOut.txt");
            break;
        case 6:
            isEuler = formiga.IsEuleriano();
            if (isEuler)
                cout << "O grafo é Euleriano!" << endl;
            else
                cout << "O grafo NÃO é Euleriano!" << endl;
            break;
        case 7:
            cout << "Número de componentes conexos : " << formiga.compsConexos() << endl;
            break;

        case 8:

            cout << "Informe o vértice de origem: ";
            cin >> label;
            formiga.depthFirstSearch(label);
            break;
        case 9:
            formiga.breadthFirstSearch();
            break;
        case 10:
            formiga.PrintInfos();
            break;
        default:
            break;
        }
    }

    return 0;
}